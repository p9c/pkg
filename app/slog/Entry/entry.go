// Package Entry is a message type for logi log entries
package Entry

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/p9c/pod/pkg/coding/simplebuffer/String"

	"github.com/p9c/pkg/coding/simplebuffer"
	"github.com/p9c/pkg/coding/simplebuffer/Time"
)

var EntryMagic = []byte{'e', 'n', 't', 'r'}

type Container struct {
	simplebuffer.Container
}

// Entry is a log entry to be printed as json to the log file
type Entry struct {
	Time         time.Time
	Level        string
	Package      string
	CodeLocation string
	Text         string
}

func Get(ent *Entry) Container {
	return Container{*simplebuffer.Serializers{
		Time.New().Put(ent.Time),
		String.New().Put(ent.Level),
		String.New().Put(ent.Package),
		String.New().Put(ent.CodeLocation),
		String.New().Put(ent.Text),
	}.CreateContainer(EntryMagic)}
}

// LoadContainer takes a message byte slice payload and loads it into a container ready to be decoded
func LoadContainer(b []byte) (out *Container) {
	out = &Container{simplebuffer.Container{Data: b}}
	return
}

func (c *Container) GetTime() time.Time {
	return Time.New().DecodeOne(c.Get(0)).Get()
}

func (c *Container) GetLevel() string {
	return String.New().DecodeOne(c.Get(1)).Get()
}

func (c *Container) GetPackage() string {
	return String.New().DecodeOne(c.Get(2)).Get()
}

func (c *Container) GetCodeLocation() string {
	return String.New().DecodeOne(c.Get(3)).Get()
}

func (c *Container) GetText() string {
	return String.New().DecodeOne(c.Get(4)).Get()
}

func (c *Container) String() (s string) {
	spew.Sdump(*c.Struct())
	return
}

// Struct deserializes the data all in one go by calling the field deserializing functions into a structure containing
// the fields. The height is given in this report as it is part of the job message and makes it faster for clients to
// look up the algorithm name according to the block height, which can change between hard fork versions
func (c *Container) Struct() (out *Entry) {
	out = &Entry{
		Time:         c.GetTime(),
		Package:      c.GetPackage(),
		Level:        c.GetLevel(),
		CodeLocation: c.GetCodeLocation(),
		Text:         c.GetText(),
	}
	return
}
