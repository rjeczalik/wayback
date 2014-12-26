package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rjeczalik/wayback"
)

const usage = `usage: wayback [-t TIME|TIMESTAMP] URL`

var (
	url string
	t   Timestamp
)

var client wayback.Client

var layouts = []string{
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC850,
	time.RFC1123,
	time.RFC3339,
}

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
func (t *Timestamp) Set(s string) error {
	if d, err := strconv.ParseUint(s, 10, 64); err == nil {
		t.Timestamp = wayback.Timestamp(d) // TODO(rjeczalik): validate timestamp
		return nil
	}
	for _, layout := range layouts {
		if d, err := time.Parse(layout, s); err == nil {
			t.Timestamp = wayback.NewTimestamp(d)
			return nil
		} else {
			fmt.Printf("Timestamp.Set(%s)=%v\n", s, err)
		}
	}
	return errors.New("invalid time/timestamp value: " + s)
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
	if t.Timestamp != 0 {
		cached, when, err = client.AvailableAt(url, uint64(t.Timestamp))
	} else {
		cached, when, err = client.Available(url)
	}
	if err != nil {
		die(err)
	}
	fmt.Printf("%s\t%s\n", cached, when.Format(time.ANSIC))
}
