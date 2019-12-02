package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kokukuma/siwa-cli/siwa"
)

var (
	addr     = flag.String("addr", ":443", "address")
	certfile = flag.String("certfile", "/etc/letsencrypt/live/siwa.kokukuma.com/fullchain.pem", "certfile")
	keyfile  = flag.String("keyfile", "/etc/letsencrypt/live/siwa.kokukuma.com/privkey.pem", "keyfile")
)

func main() {
	http.HandleFunc("/localhost", siwa.Redirector("http://localhost:8080"))
	err := http.ListenAndServeTLS(*addr, *certfile, *keyfile, nil)
	if err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}
