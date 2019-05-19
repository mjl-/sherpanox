package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/mjl-/nox/noxhttp"
	sherpaclient "github.com/mjl-/sherpa/client"
)

func check(err error, action string) {
	if err != nil {
		log.Fatalf("%s: %s\n", action, err)
	}
}

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		log.Println("usage: sherpanoxclient baseURL")
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(2)
	}

	noxhttp.RegisterDefaultTransport()
	client := &sherpaclient.Client{
		BaseURL:    args[0],
		HTTPClient: http.DefaultClient,
	}

	var result interface{}
	err := client.Call(context.Background(), &result, "exampleFunction")
	check(err, "call")
	log.Printf("result %#v\n", result)
}
