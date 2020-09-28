package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync/atomic"
	"time"

	"github.com/p9c/pkg/app/data"
	"github.com/p9c/pkg/app/disrupt"
	"github.com/p9c/pkg/app/slog"

	"github.com/urfave/cli"
)

type HookFunc *func(interface{}) error

func ToHookFunc(fn func(interface{}) error) HookFunc {
	return &fn
}

type confItem struct {
	// value is atomic so there is no conflict between goroutines accessing the value
	value *atomic.Value
	// hooks are attached by calling the AddHook method and are called when the value is set
	hooks []HookFunc
	typ   string
}

type ConfCategory map[string]confItem
type ConfByCategories map[string]ConfCategory

type defaults struct {
	flags []cli.Flag
	confs []ConfByCategories
}

type App struct {
	*cli.App
	name string
	// Note that the namespace for the name field is flat, the urfave/cli library will throw an error upon generating
	// the app data structure if multiple flags share a name (this means the first word before the comma if shortcuts
	// are present)
	conf      ConfByCategories
	defaults  defaults
	jsonCache []byte
}

type marshalledConfItem struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type marshalledConfCategory map[string]marshalledConfItem
type marshalledConfByCategories map[string]marshalledConfCategory

// MarshalJSON creates the wire form of the configuration item
func (c confItem) MarshalJSON() (out []byte, err error) {
	mc := marshalledConfItem{}
	mc.Type = c.typ
	mc.Value = c.value.Load()
	out, err = json.Marshal(mc)
	slog.Check(err)
	return
}

// UnmarshalJSON unmarshals a configuration item
func (c confItem) UnmarshalJSON(data []byte) (err error) {
	mc := marshalledConfItem{}
	if err = json.Unmarshal(data, &mc); !slog.Check(err) {
		c.value.Store(mc.Value)
		c.typ = mc.Type
	}
	return
}

// NewApp returns a new App note that this by default includes a datadir for saving the flags defined for the app and a
// version subcommand (in addition to the flags)
func NewApp(name, version, description, copyright string,
	action, before, after func(c *cli.Context) (err error),
) (a *App) {
	// incorporate shutdown interrupt library so 'after' can centralise the clean shutdown
	before = func(c *cli.Context) (err error) {
		disrupt.AddHandler(func() {
			after(c)
		})
		return before(c)
	}
	a = &App{
		name: name,
		App: &cli.App{
			Name:                 name,
			Version:              version,
			Description:          description,
			Copyright:            copyright,
			Action:               action,
			Before:               before,
			After:                after,
			EnableBashCompletion: true,
		},
		conf: make(ConfByCategories),
	}
	datadir := data.Dir(a.Name, false)
	a.defaults.flags = []cli.Flag{
		a.Bool("app", "defaults",
			"sets to base default", ""),
	}
	a.defaults.confs = []ConfByCategories{
		// create first empty one to run on init
		make(ConfByCategories),
	}
	a.SetCommands(
		a.NewCommand("version", "print version and exit",
			func(c *cli.Context) error {
				fmt.Println(c.App.Name, c.App.Version)
				return nil
			}, a.SubCommands(), nil, "v"),
		a.NewCommand(
			"reset",
			"reset settings to default",
			func(c *cli.Context) error {
				slog.Debug("resetting settings to default")
				// todo: check for flags and select index by this
				defOpt := 0
				for i := range a.defaults.confs[defOpt] {
					for j := range a.conf[i] {
						a.conf[i][j] =
							a.defaults.confs[defOpt][i][j]
					}
				}
				return nil
			},
			a.SubCommands(), a.defaults.flags),
	)
	a.SetFlags(
		a.String("app", "datadir, D",
			"sets the directory where app data will be stored",
			datadir,
			strings.ToUpper(a.Name)+"_DATADIR",
		),
	)
	return a
}

func (a *App) AddToAfter(before bool, fn func(c *cli.Context) (err error)) {
	if before {
		a.After = func(c *cli.Context) (err error) {
			if err = a.After(c); slog.Check(err) {
				return fn(c)
			}
			return
		}
	} else {
		a.After = func(c *cli.Context) (err error) {
			if err = fn(c); slog.Check(err) {
				return a.After(c)
			}
			return
		}
	}
}

// LoadConf reads the configuration located in the datadir folder into the app config database
func (a *App) LoadConf() (err error) {
	var datadir string
	if datadir, err = a.GetString("app", "datadir"); !slog.Check(err) {
		path := Join(datadir, a.name+".json")
		if FileExists(path) {
			slog.Debug("found config at path", datadir)
			var d []byte
			if d, err = ioutil.ReadFile(path); slog.Check(err) {
			} else {
				mc := make(marshalledConfByCategories)
				if err = json.Unmarshal(d, &mc); slog.Check(err) {
				} else {
					slog.Debug("loaded configuration")
					for i := range mc {
						for j := range mc[i] {
							ac := a.conf[i][j]
							ac.typ = mc[i][j].Type
							switch ac.typ {
							case "string":
								typed := mc[i][j].Value.(string)
								ac.value.Store(&typed)
							case "stringslice":
								typed := mc[i][j].Value.([]interface{})
								s := cli.StringSlice{}
								for i := range typed {
									s = append(s, typed[i].(string))
								}
								ac.value.Store(&s)
							case "bool":
								typed := mc[i][j].Value.(bool)
								ac.value.Store(&typed)
							case "int":
								typed := mc[i][j].Value.(float64)
								s := int(typed)
								ac.value.Store(&s)
							case "uint":
								typed := mc[i][j].Value.(float64)
								s := uint(typed)
								ac.value.Store(&s)
							case "duration":
								typed := mc[i][j].Value.(float64)
								s := time.Duration(typed)
								ac.value.Store(&s)
							case "float64":
								typed := mc[i][j].Value.(float64)
								ac.value.Store(&typed)
							default:
								panic("element " + i + " " + j + "lacks a type spec")
							}
						}
					}
					// slog.Debug(spew.Sdump(a.conf))
				}

			}
		} else {
			err = errors.New("no configuration file found")
		}
	}
	return
}

// SaveConf writes the configuration data to disk in json format
func (a *App) SaveConf() {
	var err error
	if a.jsonCache, err = json.MarshalIndent(a.conf, "", "  "); slog.Check(err) {
	}
	// load configuration from disk if it exists
	if datadir, err := a.GetString("app", "datadir"); !slog.Check(err) {
		path := Join(datadir, a.name+".json")
		slog.Debug("saving configuration to", path)
		EnsureDir(path)
		if err := ioutil.WriteFile(path, a.jsonCache, 0700); slog.Check(err) {
		}
	} else {
		// This should not happen since we know we made it in the NewApp function
		panic(err)
	}
}

// Initialize values that can only be generated after the generators run
func (a *App) Initialize() {
	// This saves the defaults as in the spec into a slot so it can be selected by the reset function (index 0) The
	// consuming app can add more using the functions below
	for i := range a.conf {
		a.defaults.confs[0][i] = make(ConfCategory)
		for j := range a.conf[i] {
			a.defaults.confs[0][i][j] = a.conf[i][j]
		}
	}
	// load configuration from disk
	if err := a.LoadConf(); slog.Check(err) {
		// if there is no configuration, save the default
		a.SaveConf()
	}
}

// SetDefaults loads a set of flag and configured structures into the app this is used to completely reset the defaults
// This should be run after Initialize
func (a *App) SetDefaults(d []cli.Flag, c []ConfByCategories) {
	a.defaults.flags = d
	a.defaults.confs = c
}

// AddDefaults appends a set of flags and preset configurations to the existing this is usually what the consuming app
// will use This should be run after Initialize
func (a *App) AddDefaults(d []cli.Flag, c []ConfByCategories) {
	a.defaults.flags = append(a.defaults.flags, d...)
	a.defaults.confs = append(a.defaults.confs, c...)
}

// AddHook adds a function to react when a configuration Value is changed
func (a *App) AddHook(group, name string, fn HookFunc) {
	if a.FlagExists(group, name) {
		v := a.conf[group][name]
		v.hooks = append(v.hooks, fn)
	}
}

// RemoveHook removes a change hook from a flag
func (a *App) RemoveHook(group, name string, fn HookFunc) {
	if a.FlagExists(group, name) {
		f := a.conf[group][name]
		for v, i := range f.hooks {
			if i == fn {
				if v < len(f.hooks) {
					f.hooks = append(f.hooks[:v], f.hooks[v+1:]...)
					return
				} else {
					f.hooks = f.hooks[:v]
				}
			}
		}
	}
}

// SetCommands places a slice of cli.Command into an App
func (a *App) SetCommands(c ...cli.Command) {
	a.App.Commands = c
}

// SetFlags copies in the flags
func (a *App) SetFlags(f ...cli.Flag) {
	a.App.Flags = f
}

// AddCommands appends a slice of cli.Command into an App's commands
func (a *App) AddCommands(c ...cli.Command) {
	a.App.Commands = append(a.App.Commands, c...)
}

// AddFlags appends flags to existing flags
func (a *App) AddFlags(f ...cli.Flag) {
	a.App.Flags = append(a.App.Flags, f...)
}

// NewCommand returns a cli.Command
func (a *App) NewCommand(name string, usage string, action interface{},
	subcommands cli.Commands, flags []cli.Flag, aliases ...string) cli.Command {
	return cli.Command{
		Name:        name,
		Aliases:     aliases,
		Usage:       usage,
		Action:      action,
		Subcommands: subcommands,
		Flags:       flags,
	}
}

// SubCommands returns a slice of cli.Command
func (a *App) SubCommands(sc ...cli.Command) []cli.Command {
	var c []cli.Command
	return append(c, sc...)
}

func getFirstName(name string) string {
	if strings.Contains(name, ",") {
		tmp := strings.Split(name, ",")
		return tmp[0]
	}
	return name
}

func (a *App) addItem(group, name string, value interface{}, typ string, hooks ...HookFunc) {
	if typ == "" {
		panic("type not set on item")
	}
	first := getFirstName(name)
	if !a.GroupExists(group) {
		a.conf[group] = make(ConfCategory)
	}
	val := atomic.Value{}
	val.Store(value)
	f := a.conf[group]
	f[first] = confItem{
		value: &val,
		hooks: hooks,
		typ:   typ,
	}
}

// GroupExists returns true if the configuration group exists
func (a *App) GroupExists(group string) (ok bool) {
	_, ok = a.conf[group]
	return ok
}

// FlagExists returns true if the configuration item exists
func (a *App) FlagExists(group, name string) (ok bool) {
	var v map[string]confItem
	if v, ok = a.conf[group]; ok {
		_, ok = v[name]
	}
	return ok
}

// SetConf changes a flag Value and runs any hooks attached to it
func (a *App) SetConf(group, name string, value interface{}) (err error) {
	if !a.GroupExists(group) {
		err = errors.New("group named " + group + " does not exist")
	} else if a.FlagExists(group, name) {
		f := a.conf[group][name]
		if f.value != nil {
			f.value.Store(value)
			for _, i := range f.hooks {
				if err = (*i)(value); slog.Check(err) {
					break
				}
			}
		} else {
			err = errors.New("atomic storage for " + group + "/" + name +
				" has not been created")
		}
		// flush new config to disk
		a.SaveConf()
	} else {
		err = errors.New("unable to set flag " + group + "/" +
			name + " item does not exist")
	}
	return
}

// GetConf returns the interface containing a configuration item at a given path
func (a *App) GetConf(group, name string) (o interface{}, err error) {
	if a.FlagExists(group, name) {
		v := a.conf[group][name]
		if v.value != nil {
			o = v.value.Load()
		} else {
			err = errors.New("atomic storage for " + group + "/" + name + " has not been created")
		}
	} else {
		err = errors.New("index not found: " + group + "/" + name)
	}
	return
}

// String returns an cli.StringFlag
func (a *App) String(group, name, usage, value, envVar string,
	sanitizer ...HookFunc) *cli.StringFlag {
	a.addItem(group, name, &value, "string", sanitizer...)
	for _, i := range sanitizer {
		if err := (*i)(&value); slog.Check(err) {
			break
		}
	}
	return &cli.StringFlag{
		Name:        name,
		Usage:       usage,
		Value:       value,
		EnvVar:      envVar,
		Destination: &value,
	}
}

// GetString returns the pointer to the string of the config item Value
func (a *App) GetString(group, name string) (o string, err error) {
	var val interface{}
	val, err = a.GetConf(group, name)
	if !slog.Check(err) {
		if v, ok := val.(*string); !ok {
			err = errors.New("Value at " + group + "/" + name + "is not a string")
		} else {
			o = *v
		}
	}
	return
}

// SetString sets the Value stored in the config item
func (a *App) SetString(group, name, value string) (err error) {
	return a.SetConf(group, name, &value)
}

// BoolTrue returns a CliBoolFlag that defaults to true
func (a *App) BoolTrue(group, name, usage, envVar string,
	sanitizer ...HookFunc) *cli.BoolTFlag {
	value := true
	a.addItem(group, name, &value, "bool", sanitizer...)
	for _, i := range sanitizer {
		if err := (*i)(&value); slog.Check(err) {
			break
		}
	}
	return &cli.BoolTFlag{
		Name:        name,
		Usage:       usage,
		EnvVar:      envVar,
		Destination: &value,
	}
}

// Bool returns an cli.BoolFlag
func (a *App) Bool(group, name, usage, envVar string,
	sanitizer ...HookFunc) *cli.BoolFlag {
	value := false
	a.addItem(group, name, &value, "bool", sanitizer...)
	for _, i := range sanitizer {
		if err := (*i)(&value); slog.Check(err) {
			break
		}
	}
	return &cli.BoolFlag{
		Name:        name,
		Usage:       usage,
		EnvVar:      envVar,
		Destination: &value,
	}
}

// GetBool returns the pointer to the bool of the config item Value
func (a *App) GetBool(group, name string) (o bool, err error) {
	var val interface{}
	val, err = a.GetConf(group, name)
	if !slog.Check(err) {
		if v, ok := val.(*bool); !ok {
			err = errors.New("Value at " + group + "/" + name + "is not a bool")
		} else {
			o = *v
		}
	}
	return
}

// SetBool sets the Value stored in the config item
func (a *App) SetBool(group, name string, value bool) (err error) {
	return a.SetConf(group, name, &value)
}

// StringSlice returns and cli.StringSliceFlag
func (a *App) StringSlice(group, name, usage string, val []string, envVar string,
	sanitizer ...HookFunc) *cli.StringSliceFlag {
	value := cli.StringSlice(val)
	a.addItem(group, name, &value, "stringslice", sanitizer...)
	for _, i := range sanitizer {
		if err := (*i)(&value); slog.Check(err) {
			break
		}
	}
	return &cli.StringSliceFlag{
		Name:   name,
		Usage:  usage,
		EnvVar: envVar,
		Value:  &value,
	}
}

// GetStringSlice returns the pointer to the string slice of the config item Value
func (a *App) GetStringSlice(group, name string) (o cli.StringSlice,
	err error) {
	var val interface{}
	val, err = a.GetConf(group, name)
	if !slog.Check(err) {
		if v, ok := val.(*cli.StringSlice); !ok {
			err = errors.New("Value at " + group + "/" + name + "is not a string slice")
		} else {
			o = *v
		}
	}
	return
}

// SetStringSlice sets the Value stored in the config item
func (a *App) SetStringSlice(group, name string, value []string) (err error) {
	return a.SetConf(group, name, &value)
}

// Int returns an cli.IntFlag
func (a *App) Int(group, name, usage string, value int, envVar string,
	sanitizer ...HookFunc) *cli.IntFlag {
	a.addItem(group, name, &value, "int", sanitizer...)
	for _, i := range sanitizer {
		if err := (*i)(&value); slog.Check(err) {
			break
		}
	}
	return &cli.IntFlag{
		Name:        name,
		Value:       value,
		Usage:       usage,
		EnvVar:      envVar,
		Destination: &value,
	}
}

// GetInt returns the pointer to the int of the config item Value
func (a *App) GetInt(group, name string) (o int, err error) {
	var val interface{}
	val, err = a.GetConf(group, name)
	if !slog.Check(err) {
		if v, ok := val.(*int); !ok {
			err = errors.New("Value at " + group + "/" + name + "is not a int")
		} else {
			o = *v
		}
	}
	return
}

// SetInt sets the Value stored in the config item
func (a *App) SetInt(group, name string, value int) (err error) {
	return a.SetConf(group, name, &value)
}

// Uint returns an cli.UintFlag
func (a *App) Uint(group, name, usage string, value uint, envVar string,
	sanitizer ...HookFunc) *cli.UintFlag {
	a.addItem(group, name, &value, "uint", sanitizer...)
	for _, i := range sanitizer {
		if err := (*i)(&value); slog.Check(err) {
			break
		}
	}
	return &cli.UintFlag{
		Name:        name,
		Value:       value,
		Usage:       usage,
		EnvVar:      envVar,
		Destination: &value,
	}
}

// GetUint returns the pointer to the uint of the config item Value
func (a *App) GetUint(group, name string) (o uint, err error) {
	var val interface{}
	val, err = a.GetConf(group, name)
	if !slog.Check(err) {
		if v, ok := val.(*uint); !ok {
			err = errors.New("Value at " + group + "/" + name + "is not a uint")
		} else {
			o = *v
		}
	}
	return
}

// SetUint sets the Value stored in the config item
func (a *App) SetUint(group, name string, value uint) (err error) {
	return a.SetConf(group, name, &value)
}

// Duration returns an cli.DurationFlag
func (a *App) Duration(group, name, usage string, value time.Duration, envVar string,
	sanitizer ...HookFunc) *cli.DurationFlag {
	a.addItem(group, name, &value, "duration", sanitizer...)
	for _, i := range sanitizer {
		if err := (*i)(&value); slog.Check(err) {
			break
		}
	}
	return &cli.DurationFlag{
		Name:        name,
		Value:       value,
		Usage:       usage,
		EnvVar:      envVar,
		Destination: &value,
	}
}

// GetDuration returns the pointer to the time.Duration of the config item Value
func (a *App) GetDuration(group, name string) (o time.Duration, err error) {
	var val interface{}
	val, err = a.GetConf(group, name)
	if !slog.Check(err) {
		if v, ok := val.(*time.Duration); !ok {
			err = errors.New("Value at " + group + "/" + name + "is not a duration")
		} else {
			o = *v
		}
	}
	return
}

// SetDuration sets the Value stored in the config item
func (a *App) SetDuration(group, name string, value time.Duration) (err error) {
	return a.SetConf(group, name, &value)
}

// Float64 returns an cli.Float64Flag
func (a *App) Float64(group, name, usage string, value float64, envVar string,
	sanitizer ...HookFunc) *cli.Float64Flag {
	a.addItem(group, name, &value, "float64", sanitizer...)
	for _, i := range sanitizer {
		if err := (*i)(&value); slog.Check(err) {
			break
		}
	}
	return &cli.Float64Flag{
		Name:        name,
		Value:       value,
		Usage:       usage,
		EnvVar:      envVar,
		Destination: &value,
	}
}

// GetFloat64 returns the pointer to the float64 of the config item Value
func (a *App) GetFloat64(group, name string) (o float64, err error) {
	var val interface{}
	val, err = a.GetConf(group, name)
	if !slog.Check(err) {
		if v, ok := val.(*float64); !ok {
			err = errors.New("Value at " + group + "/" + name + "is not a float64")
		} else {
			o = *v
		}
	}
	return

}

// SetFloat64 sets the Value stored in the config item
func (a *App) SetFloat64(group, name string, value float64) (err error) {
	return a.SetConf(group, name, &value)
}
