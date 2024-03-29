package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/dustin/go-humanize"
)

var opt options

type options struct {
	Chunks      int64  `short:"n" default:"3" help:"split into n parts"`
	SizePerPart string `short:"b" help:"split into sized parts"`
	Verbose     bool   `short:"v" help:"Verbose printing."`
	About       bool   `help:"About."`

	Files []string `arg:"" optional:""`
}

func (o *options) splitSize(total int64) int64 {
	if o.partSize() > 0 {
		return o.partSize()
	} else {
		size := total / o.Chunks
		if size*o.Chunks < total {
			size += 1
		}
		return size
	}
}

func (o *options) partSize() int64 {
	size, _ := humanize.ParseBytes(o.SizePerPart)
	return int64(size)
}

func main() {
	ctx := kong.Parse(&opt,
		kong.Name("splits"),
		kong.Description("Command line tool for splitting file into parts."),
		kong.UsageOnError(),
	)

	switch {
	case opt.About:
		fmt.Println("Visit https://github.com/gonejack/splits")
		return
	case len(opt.Files) == 0, opt.partSize() == 0 && opt.Chunks == 0:
		_ = ctx.PrintUsage(false)
	default:
		for _, f := range opt.Files {
			if e := split(f); e != nil {
				log.Fatalf("split %s failed: %s", f, e)
			}
		}
	}
}

func split(name string) (err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return
	}

	if opt.Verbose {
		log.Printf("split %s(%s)", name, humanize.Bytes(uint64(stat.Size())))
	}

	wrt, idx := int64(0), 0
	for wrt < stat.Size() {
		p := fmt.Sprintf("%s.%d", filepath.Base(name), idx)

		n, e := part(f, p, stat.Mode(), stat.Size())
		if opt.Verbose {
			log.Printf("part#%d => %s(%s)", idx, p, humanize.Bytes(uint64(n)))
		}

		switch e {
		case nil:
			wrt, idx = wrt+n, idx+1
		case io.EOF:
			return nil
		default:
			return e
		}
	}
	return
}

func part(r io.Reader, name string, mod os.FileMode, total int64) (n int64, err error) {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mod)
	if err != nil {
		return
	}
	defer f.Close()
	return io.CopyN(f, r, opt.splitSize(total))
}
