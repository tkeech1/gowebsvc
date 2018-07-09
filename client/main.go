package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main__() {
	ctx := context.Background()
	// cancel the request after one second
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	// this defer is needed for the timeout to work
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatal(res.Status)
	}

	io.Copy(os.Stdout, res.Body)

}
