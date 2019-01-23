package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	//Actions
	apply := flag.Bool("apply", false, "Apply a file")

	//General
	saveEndpoint := flag.Bool("save-endpoint", false, "Whether or not to save the url to file.")
	endpoint := flag.String("endpoint", "http://127.0.0.1:6227", "Endpoint of the PC API server.")

	//Specific
	file := flag.String("file", "", "Json file to apply.")

	flag.Parse()
	fmt.Printf("apply: %v, saveEndpoint: %v, file: %s, endpoint: %s\n", *apply, *saveEndpoint, *file, *endpoint)

	//TODO: lots of checking....

	//TODO: persist to file
	if *saveEndpoint {
		//Where to save?
		//runtime.GOOS
	}

	//Route the action
	if *apply {
		handleApply(endpoint, file)
	} else {
		panic("Nothing to do!")
	}

}

func handleApply(endpoint, file *string) {
	b, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("Failed to read file %s: '%s'\n", *file, err)
	}
	client := &http.Client{}
	client.Timeout = time.Second * 5

	uri := fmt.Sprint(*endpoint, "/route")
	body := bytes.NewBuffer(b)
	req, err := http.NewRequest(http.MethodPut, uri, body)
	if err != nil {
		log.Fatalf("http.NewRequest() failed with '%s'\n", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("client.Do() failed with '%s'\n", err)
	}

	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll() failed with '%s'\n", err)
	}

	fmt.Printf("Response status code: %d, text:\n%s\n", resp.StatusCode, string(b))
}
