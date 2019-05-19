package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"bitbucket.org/mjl/httpasset"
	"github.com/mjl-/nox"
	"github.com/mjl-/sherpa"
	"github.com/mjl-/sherpadoc"
)

var (
	httpFS  http.FileSystem
	version = "dev"
	address = flag.String("address", "localhost:1047", "nox-extended address to serve sherpanox on")
)

func init() {
	httpFS = httpasset.Fs()
	if err := httpasset.Error(); err != nil {
		log.Println("falling back to local assets:", err)
		httpFS = http.Dir("assets")
	}
}

func check(err error, action string) {
	if err != nil {
		log.Fatalf("%s: %s\n", action, err)
	}
}

func main() {
	log.SetFlags(0)
	flag.Usage = func() {
		log.Printf("usage: sherpanox [flags]\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 0 {
		flag.Usage()
		os.Exit(2)
	}

	var doc sherpadoc.Section
	ff, err := httpFS.Open("/example.json")
	check(err, "opening sherpa docs")
	err = json.NewDecoder(ff).Decode(&doc)
	check(err, "parsing sherpa docs")
	err = ff.Close()
	check(err, "closing sherpa docs after parsing")

	exampleHandler, err := sherpa.NewHandler("/example/", version, Example{}, &doc, nil)
	check(err, "making sherpa handler")

	http.Handle("/example/", exampleHandler)
	http.HandleFunc("/", serveAsset)

	config := &nox.Config{}
	listener, err := nox.Listen("tcp", *address, config)
	check(err, "listen")

	log.Printf("sherpanox, version %s, listening on %s, local static public key %s", version, config.Address, config.LocalStaticPublic())
	log.Fatal(http.Serve(listener, nil))
}

func serveAsset(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/") {
		r.URL.Path += "index.html"
	}
	f, err := httpFS.Open("/web" + r.URL.Path)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		log.Printf("serving asset %s: %s\n", r.URL.Path, err)
		http.Error(w, "500 - Server error", 500)
		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		log.Printf("serving asset %s: %s\n", r.URL.Path, err)
		http.Error(w, "500 - Server error", 500)
		return
	}

	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	_, haveCacheBuster := r.URL.Query()["v"]
	cache := "no-cache, max-age=0"
	if haveCacheBuster {
		cache = fmt.Sprintf("public, max-age=%d", 31*24*3600)
	}
	w.Header().Set("Cache-Control", cache)

	http.ServeContent(w, r, r.URL.Path, info.ModTime(), f)
}

// Example is an API for a sherpa-over-nox demo.
type Example struct{}

// ExampleFunction just returns "hello!".
func (Example) ExampleFunction(ctx context.Context) string {
	return "hello!"
}
