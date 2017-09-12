package main

import (
	"net/url"
	"os"

	"github.com/msales/pkg/log"
	"github.com/msales/pkg/stats"
	"github.com/nrwiersma/ren"
	"gopkg.in/inconshreveable/log15.v2"
	"fmt"
	"bytes"
	"strings"
)

// Application =============================

func newApplication(c *Context) (*ren.Application, error) {
	app := ren.NewApplication(c.String(FlagTemplates))

	return app, nil
}

// Logger ==================================

func newLogger(c *Context) (log15.Logger, error) {
	lvl := c.String(FlagLogLevel)
	v, err := log15.LvlFromString(lvl)
	if err != nil {
		return nil, err
	}

	h := log15.LvlFilterHandler(v, log15.StreamHandler(os.Stderr, log15.FormatFunc(logFormat)))
	if lvl == "debug" {
		h = log15.CallerFileHandler(h)
	}

	l := log15.New()
	l.SetHandler(log15.LazyHandler(h))

	return l, nil
}

func logFormat(r *log15.Record) []byte {
	b := &bytes.Buffer{}
	lvl := strings.ToUpper(r.Lvl.String())
	fmt.Fprintf(b, "%s %s %s ", r.Time.UTC().Format("2006-01-02T15:04:05"), lvl, r.Msg)
	b.WriteByte('\n')
	return b.Bytes()
}

// Stats ===================================

func newStats(c *Context) (stats.Stats, error) {
	dsn := c.String(FlagStats)
	if dsn == "" {
		return stats.Null, nil
	}

	uri, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	switch uri.Scheme {
	case "statsd":
		return newStatsdStats(uri.Host)

	case "l2met":
		return newL2metStats(c.logger), nil

	default:
		return stats.Null, nil
	}
}

func newStatsdStats(addr string) (stats.Stats, error) {
	return stats.NewStatsd(addr, "ren")
}

func newL2metStats(log log.Logger) (stats.Stats) {
	return stats.NewL2met(log, "ren");
}
