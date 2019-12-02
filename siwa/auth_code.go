package siwa

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

// AuthResp represents auth endpoint result
type AuthResp struct {
	Code    string
	IDToken string
	State   string
	User    string
}

// GetAuthCode start localserver to get auth code.
// To get the auth code, response has to be redirect to local:8080
func GetAuthCode(ctx context.Context, state string) (*AuthResp, error) {
	ch := make(chan AuthResp, 5)
	fmt.Println("Start local server")
	srv := startHTTPServer(ch)

	var ar AuthResp
	for {
		select {
		case ar = <-ch:
			if ar.State != state {
				fmt.Printf("Unexpected response %v\n", ar)
				continue
			}
			fmt.Println("Got result")
			fmt.Println("Shutdown local server")
			if err := srv.Shutdown(context.TODO()); err != nil {
				return nil, err
			}
			return &ar, nil
		}
	}
}

func startHTTPServer(ch chan<- AuthResp) *http.Server {
	srv := &http.Server{Addr: "localhost:8080"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		fmt.Fprintf(w, "Got response parameter. Close the browser.")

		spew.Dump(r.Form)

		ch <- AuthResp{
			Code:    r.Form.Get("code"),
			IDToken: r.Form.Get("id_token"),
			State:   r.Form.Get("state"),
			User:    r.Form.Get("user"),
		}
	})

	go func() {
		// returns ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	return srv
}
