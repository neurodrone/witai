# wit.ai [![GoDoc](https://godoc.org/github.com/neurodrone/witai?status.svg)](https://godoc.org/github.com/neurodrone/witai) [![](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/neurodrone/witai/blob/master/LICENSE)

Go library for `wit.ai` Natural Language Processing API. This library
integrates with `libwit` locally and can perform both blocking and
non-blocking voice queries using `wit.ai` service.

## Dependencies

In order to start using this library you will need an account on https://wit.ai/. Make sure you complete the [Quick Start](https://wit.ai/docs/console/quickstart) guide to get yourself all set up.

Other dependencies:

 * Latest `libwit` (can be downloaded from [wit.ai Releases](https://github.com/wit-ai/libwit/releases)).
 * `libsox` (`brew install sox` on OSX).
 * `libcurl` (`brew install curl` on OSX).

## Usage

```go
package main

import (
	"flag"
	"log"

	"github.com/ianschenck/envflag"
	"github.com/neurodrone/witai"
)

func main() {
	var (
		// This points to the client access token for your account on wit.ai, that you
		// will possibly have stored within an environment variable.
		accessToken = envflag.String("ACCESS_TOKEN", "", "WIT client access token")

		// The recording device you will use for voice queries.
		// Usually, you can leave it be default.
		device = flag.String("device", witai.DefaultDevice, "device name for recording input")
	)
	envflag.Parse()
	flag.Parse()

	// Create a new wit-ai context that will be used for queries.
	ctx, err := witai.NewContext(*device, *accessToken, witai.Error)
	if err != nil {
		log.Fatalln("cannot create new wit-ai context:", err)
	}

	// Always make sure to close the context once you are done.
	defer ctx.Close()

	log.Println("Say something nice now: ...")

	done := make(chan struct{})

	// Query the wit.ai voice service asyncly.
	if err := ctx.VoiceQueryAutoAsync(func(s string) {
		r, err := witai.NewResult(s)
		if err != nil || !r.IsValid() {
			return
		}
		log.Printf("Result: %q\n", r.Outcomes[0].Text)

		// We can exit now that we have the result.
		close(done)
	}); err != nil {
		log.Fatalln("cannot query wit-ai:", err)
	}

	// Wait exiting the process until the async result returns.
	<-done
}
```

## Installation

As simple as:

```
go get github.com/neurodrone/witai
```
