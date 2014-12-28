wayback [![GoDoc](https://godoc.org/github.com/rjeczalik/wayback?status.svg)](https://godoc.org/github.com/rjeczalik/wayback) [![Build Status](https://img.shields.io/travis/rjeczalik/wayback/master.svg)](https://travis-ci.org/rjeczalik/wayback "linux_amd64") [![Build status](https://img.shields.io/appveyor/ci/rjeczalik/wayback.svg)](https://ci.appveyor.com/project/rjeczalik/wayback "windows_amd64")
=========

Package wayback implements a client for Wayback Availability JSON API.
See its website for details: https://archive.org/help/wayback_api.php

*Example usage*

```go
package main

import (
	"log"

	"github.com/rjeczalik/wayback"
)

func main() {
	url, time, err := wayback.Available("github.com")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(url, "captured at", time)
}
```

## cmd/wayback [![GoDoc](https://godoc.org/github.com/rjeczalik/wayback/cmd/wayback?status.png)](https://godoc.org/github.com/rjeczalik/wayback/cmd/wayback)

*Installation*

```
~ $ go get -u github.com/rjeczalik/wayback/cmd/wayback
```

*Documentation*

[godoc.org/github.com/rjeczalik/wayback/cmd/wayback](http://godoc.org/github.com/rjeczalik/wayback/cmd/wayback)

*Example usage*

```bash
~ $ wayback github.com
http://web.archive.org/web/20141226100456/https://github.com/	Fri Dec 26 10:04:56 2014
```
```bash
~ $ wayback -t 2010 github.com
http://web.archive.org/web/20100102002654/http://github.com/	Sat Jan  2 00:26:54 2010
```
