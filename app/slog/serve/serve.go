package serve

import (
	"github.com/p9c/pod/pkg/comm/pipe"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/Entry"
	"github.com/p9c/pod/pkg/util/logi/Pkg"
	"github.com/p9c/pod/pkg/util/logi/Pkg/Pk"
	"go.uber.org/atomic"

	"github.com/p9c/pkg/app/slog"
)

func Log(quit chan struct{}, saveFunc func(p Pk.Package) (success bool)) {
	slog.Debug("starting log server")
	lc := logi.L.AddLogChan()
	pkgChan := make(chan Pk.Package)
	var logOn atomic.Bool
	logOn.Store(false)
	p := pipe.Serve(quit, func(b []byte) (err error) {
		// listen for commands to enable/disable logging
		if len(b) >= 4 {
			magic := string(b[:4])
			switch magic {
			case "run ":
				slog.Debug("setting to run")
				logOn.Store(true)
			case "stop":
				slog.Debug("stopping")
				logOn.Store(false)
			case "slvl":
				slog.Debug("setting level", logi.Levels[b[4]])
				logi.L.SetLevel(logi.Levels[b[4]], false, "pod")
			case "pkgs":
				pkgs := Pkg.LoadContainer(b).GetPackages()
				for i := range pkgs {
					(*logi.L.Packages)[i] = pkgs[i]
				}
				// save settings
				if !saveFunc(pkgs) {
					slog.Error("failed to save log filter configuration")
				}
			}
		}
		return
	})
	go func() {
	out:
		for {
			select {
			case <-quit:
				break out
			case e := <-lc:
				if logOn.Load() {
					if n, err := p.Write(Entry.Get(&e).Data); !slog.Check(err) {
						if n < 1 {
							slog.Error("short write")
						}
					}
				}
			case pk := <-pkgChan:
				if n, err := p.Write(Pkg.Get(pk).Data); !slog.Check(err) {
					if n < 1 {
						slog.Error("short write")
					}
				}
			}
		}
	}()
}
