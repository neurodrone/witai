package main

import (
	"flag"
	"log"

	"github.com/ianschenck/envflag"
	"github.com/neurodrone/witai"
)

func main() {
	var (
		accessToken = envflag.String("ACCESS_TOKEN", "", "WIT client access token")

		device = flag.String("device", witai.DefaultDevice, "device name for recording input")
	)
	envflag.Parse()
	flag.Parse()

	ctx, err := witai.NewContext(*device, *accessToken, witai.Debug)
	if err != nil {
		log.Fatalln("cannot create new wit-ai context:", err)
	}
	defer ctx.Close()

	log.Println("Say something nice: ...")

	_, o, err := ctx.VoiceQueryAuto()
	if err != nil {
		log.Fatalln("cannot query wit-ai:", err)
	}

	log.Printf("Interpreted text: %q", o.Text)
}
