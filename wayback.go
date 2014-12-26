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
	} `json:"archived_snapshots""`
}

// Timestamp represents a timestamp value which format is YYYYMMDDhhmmss.
type Timestamp uint64

// NewTimestamp gives new timestamp value converted from the given time.
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp(t.Second() + t.Minute()*100 + t.Hour()*10000 + t.Day()*1000000 +
		int(t.Month())*100000000 + t.Year()*10000000000)
}

// Time converts the timestamp to a time value.
func (t Timestamp) Time() time.Time {
	return time.Now() // TODO
}

// String converts the timestamp to string cutting of any following zeros.
func (t Timestamp) String() string {
	s := strconv.FormatUint(uint64(t), 10)
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] != '0' {
			return s[:i+1]
		}
	}
	return ""
}

// UnamarshalJSON decodes quoted string into a timestamp value.
//
// TODO(rjeczalik): validate
func (t *Timestamp) UnmarshalJSON(p []byte) error {
	s, err := strconv.Unquote(string(p))
	if err != nil {
		return err
	}
	d, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	*t = Timestamp(d) // TODO
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
func (c *Client) AvailableAt(url string, timestamp uint64) (string, time.Time, error) {
	return c.available(fmt.Sprintf(reqTimestamp, url, Timestamp(timestamp)))
}
