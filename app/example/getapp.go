package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"

	"github.com/p9c/pkg/app/disrupt"
	"github.com/p9c/pkg/app/util"

	"github.com/p9c/pkg/app/slog"
)

func main() {
	a := getApp()
	if err := a.Run(os.Args); slog.Check(err) {
	}
}

func getApp() (a *util.App) {
	a = util.NewApp("example", "v0.0.3",
		"example app",
		"The Unlicense - https://unlicense.org",
		func(c *cli.Context) (err error) {
			slog.Debug("default command")
			return
		},
		func(c *cli.Context) (err error) {
			slog.Debug("running before command")
			return
		},
		func(c *cli.Context) error {
			if disrupt.Restart {
				slog.Debug("restarting")
			} else {
				slog.Debug("subcommand completed")
			}
			return nil
		},
	)
	a.AddCommands(
		a.NewCommand("example", "here as an example",
			func(c *cli.Context) error {
				fmt.Println(c.App.Name, c.App.Version)
				return nil
			}, a.SubCommands(), nil, "e"),
	)
	a.AddFlags(
		a.String("example0", "string, s",
			"example flag",
			"example",
			strings.ToUpper(a.Name)+"_EX",
			util.ToHookFunc(func(i interface{}) (err error) {
				// custom sanitizer
				ii := i.(*string)
				*ii = *ii + "[sanitizeded]"
				return nil
			}),
		),
		a.Bool("example0", "bool, b",
			"enable this", "",
		),
		a.BoolTrue("example0", "boolt, t",
			"disable this", "",
		),
		a.Int("example1", "int, i", "integer (as per OS)",
			-69, "",
			util.IntBounds(-50, 200),
		),
		a.Uint("example1", "uint, u",
			"unsigned integer (as per OS)",
			69, "",
			util.UintBounds(50, 600),
		),
		a.Float64("example2", "float64, f",
			"64 bit floating point value", 0.696969, "",
			util.Float64Bounds(-100.1, 0.23),
		),
		a.Duration("example2", "duration",
			"duration in standard Go notation",
			time.Second*69, "",
			util.DurationBounds(time.Second, time.Hour),
		),
		a.StringSlice("example2", "stringslice, S",
			"slice of strings", cli.StringSlice{
				"one", "two", "sixty nine",
			}, ""),
		a.String("example2", "url",
			"example flag",
			"example.com",
			strings.ToUpper(a.Name)+"_EX",
			util.CheckURL(),
		),
		a.String("example2", "badurl",
			"example flag",
			"wrong!$^!#$!^#$&!#$%!#$%S_S_-','",
			strings.ToUpper(a.Name)+"_EX",
			util.CheckURL(),
		),
		a.String("example2", "path",
			"example flag",
			os.Args[0],
			strings.ToUpper(a.Name)+"_EX",
			util.CheckPath(),
		),
		a.String("example2", "badpath",
			"example flag",
			"aoeu&#&@*($)#",
			strings.ToUpper(a.Name)+"_EX",
			util.CheckPath(),
		),
	)
	a.Initialize()
	if err := a.SetBool("example0", "bool", true); slog.Check(err) {
	}
	return
}
