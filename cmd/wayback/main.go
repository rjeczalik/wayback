package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rjeczalik/wayback"
)

const usage = `usage: wayback [-t TIME|TIMESTAMP] URL`

var (
	url string
	t   Timestamp
)

// Timestamp is a wrapper for wayback.Timestamp which implements flag.Value
// interface.
type Timestamp struct {
	wayback.Timestamp
}

// Set implements flag.Value interface.
func (t Timestamp) String() string {
	return t.Timestamp.String()
}

// Set implements flag.Value interface.
func (t *Timestamp) Set(s string) (err error) {
	t.Timestamp, err = wayback.ParseTimestamp("", s)
	return
}

func die(v interface{}) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}

func init() {
	flag.Var(&t, "t", "snapshot creation time")
	flag.Parse()
	if flag.NArg() != 1 {
		die(usage)
	}
	url = flag.Arg(0)
}

func main() {
	if len(os.Args) == 1 {
		die(usage)
	}
	var (
		cached string
		when   time.Time
		err    error
	)
	if t.String() != "" {
		cached, when, err = wayback.AvailableAt(url, t.Timestamp)
	} else {
		cached, when, err = wayback.Available(url)
	}
	if err != nil {
		die(err)
	}
	fmt.Printf("%s\t%s\n", cached, when.Format(time.ANSIC))
}
