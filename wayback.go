// Package wayback implements a client for Wayback Availability JSON API.
// See its website for details: https://archive.org/help/wayback_api.php
package wayback

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const req = `http://archive.org/wayback/available?url=%s`
const reqTimestamp = `http://archive.org/wayback/available?url=%s&timestamp=%s`

var errNotAvailable = errors.New("no snapshot is available for the given url")

type response struct {
	ArchivedSnapshots struct {
		Closest struct {
			Available bool      `json:"available"`
			URL       string    `json:"url"`
			Timestamp Timestamp `json:"timestamp"`
			Status    string
		} `json:"closest"`
	} `json:"archived_snapshots"`
}

func ndigit(i uint64) (n int) {
	for ; i != 0; i /= 10 {
		n++
	}
	return
}

// TimeLayout is Timestamp's time format for use with time.Parse and
// time.(Time).Format functions.
const TimeLayout = "20060102150405"

// Timestamp represents a timestamp value which format is YYYYMMDDhhmmss.
type Timestamp struct {
	s string
}

// NewTimestamp gives new timestamp value converted from the given time.
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp{s: t.Format(TimeLayout)}
}

var layouts = []string{
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC850,
	time.RFC1123,
	time.RFC3339,
}

// ParseTimestamp tries to parse the given string into a Timestamp value. If it
// is a digit, it parses it into uint64 and calls NewTimestamp on it. Otherwise
// it tries to parse a time.Time and upon success gives new Timestamp value via
// NewTime function.
//
// TODO(rjeczalik): update doc
func ParseTimestamp(layout, s string) (Timestamp, error) {
	if layout != "" {
		t, err := time.Parse(layout, s)
		if err != nil {
			return Timestamp{}, err
		}
		return NewTimestamp(t), nil
	}
	if t, err := strconv.ParseUint(s, 10, 64); err == nil {
		s = strconv.FormatUint(t, 10)
		switch n, m := len(s), len(TimeLayout); {
		case n < m:
			s = s + strings.Repeat("0", m-n)
		case n > m:
			s = s[:m]
		}
		return Timestamp{s: s}, nil
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return NewTimestamp(t), nil
		}
	}
	return Timestamp{}, errors.New("invalid time/timestamp value: " + s)
}

// Time converts the timestamp to a time value.
func (t Timestamp) Time() time.Time {
	if t, err := time.Parse(TimeLayout, t.s); err == nil {
		return t
	}
	return time.Time{}
}

// String converts the timestamp to string. It cuts off any following zeros.
func (t Timestamp) String() string {
	i := len(t.s) - 1
	for ; i > 3; i-- {
		if t.s[i] != '0' {
			break
		}
	}
	return t.s[:i+1]
}

// UnamarshalJSON decodes quoted string into a timestamp value.
//
// TODO(rjeczalik): validate
func (t *Timestamp) UnmarshalJSON(p []byte) error {
	s, err := strconv.Unquote(string(p))
	if err != nil {
		return err
	}
	*t = Timestamp{s: s}
	return nil
}

// Client is a client for Wayback Availability JSON API.
type Client struct {
	// ClientHTTP is a HTTP client used to communicate with the Wayback Machine
	// endpoint. If it's nil, the http.DefaultClient is used.
	ClientHTTP *http.Client
}

func (c *Client) clientHTTP() *http.Client {
	if c.ClientHTTP != nil {
		return c.ClientHTTP
	}
	return http.DefaultClient
}

func (c *Client) available(url string) (s string, t time.Time, err error) {
	res, err := c.clientHTTP().Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = errors.New(http.StatusText(res.StatusCode))
		return
	}
	resp := &response{}
	buf := &bytes.Buffer{}
	if err = json.NewDecoder(io.TeeReader(res.Body, buf)).Decode(resp); err != nil {
		return
	}
	if !resp.ArchivedSnapshots.Closest.Available {
		err = errNotAvailable
		return
	}
	s = resp.ArchivedSnapshots.Closest.URL
	t = resp.ArchivedSnapshots.Closest.Timestamp.Time()
	return
}

// Available queries the archive for cached snapshot of website given by the url.
// If it's available, it returns its URL and creation time of the most recent
// snapshot. Otherwise it returns non-nil error.
func (c *Client) Available(url string) (string, time.Time, error) {
	return c.available(fmt.Sprintf(req, url))
}

// AvailableAt queries the archive for cached snapshot of website given by the url
// and creation time. If it's available, the function returns its URL and creation
// time which is the closest to the requested one.
func (c *Client) AvailableAt(url string, timestamp Timestamp) (string, time.Time, error) {
	return c.available(fmt.Sprintf(reqTimestamp, url, timestamp))
}
