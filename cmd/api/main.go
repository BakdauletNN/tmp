package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type LogginAction struct {
	logger io.Writer
	next   http.RoundTripper
}

func (l LogginAction) logging(r *http.Request) (*http.Response, error) {
	fmt.Fprintf(l.logger, "[%s] %s %s\n", time.Now().Format(time.ANSIC), r.Method, r.URL)
	return l.next.RoundTrip(r)
}

func (l LogginAction) RoundTrip(r *http.Request) (*http.Response, error) {
	return l.logging(r)
}

func main() {
	client := &http.Client{

		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Print(req.Response.Status)
			fmt.Println("redirect")
			return nil
		},
		Transport: &LogginAction{
			logger: os.Stdout,
			next:   http.DefaultTransport,
		},
	}

	rsp, err := client.Get("https://platonus.iitu.edu.kz")
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Body.Close()
	fmt.Println("Code:->", rsp.StatusCode)

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Fatal(err)

	}
	fmt.Println(string(body))

}
