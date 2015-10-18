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
