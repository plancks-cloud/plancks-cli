package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//Commands
var apply bool
var delete bool
var get bool

//Flags
var filename string
var endpoint string
var object string

func main() {
	err := readFirst()
	if err != nil {
		logrus.Error(err)
		return
	}
	readFlags()

	if !apply && !delete {
		logrus.Error(errors.New("No command. Supported commands are apply, delete."))
		return
	}

	if !apply && !delete {
		logrus.Error(errors.New("No command. Supported commands are apply, delete."))
		return
	}

	//TODO: load filename if "" from .plancks in ~/

	if (filename == "" && apply) || (filename == "" && delete) {
		logrus.Error(errors.New("No filename for command"))
		return
	}

	if (endpoint == "" && apply) || (endpoint == "" && delete) {
		endpoint = "127.0.0.1:6227"
		logrus.Println("Assuming endpoint 127.0.0.1:6227")
	}

	/// End of checking

	if apply {
		handleApply(&endpoint, &filename)
		return
	}

	if delete {
		handleDelete(&endpoint, &filename)
	}

	if get {
		//TODO
	}

}

func readFirst() (err error) {
	if len(os.Args) < 2 {
		err = errors.New("Not enough arguments. Provide either apply or delete.")
		return
	}
	if os.Args[1] == "apply" || os.Args[0] == "a" {
		apply = true
	} else if os.Args[1] == "delete" || os.Args[0] == "d" {
		delete = true
	} else if os.Args[1] == "get" || os.Args[0] == "g" {
		get = true
	}
	return
}

func readFlags() {
	for i, s := range os.Args {
		if i < 2 {
			continue
		}
		f, v, e := split(s)
		if e != nil {
			continue
		}
		if f == "-f" || f == "-filename" {
			filename = v
			continue
		}
		if f == "-e" || f == "-endpoint" {
			endpoint = v
			continue
		}
		if f == "-o" || f == "-object" {
			endpoint = v
			continue
		}
	}
}

func split(in string) (f, v string, err error) {
	if !strings.Contains(in, "=") {
		err = errors.New("No value")
		return
	}
	s := strings.Split(in, "=")
	f, v = s[0], s[1]
	return
}

func handleApply(endpoint, file *string) {
	b, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("Failed to read file %s: '%s'\n", *file, err)
	}
	client := &http.Client{}
	client.Timeout = time.Second * 5

	uri := fmt.Sprint(*endpoint, "/apply")
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

func handleDelete(endpoint, file *string) {
	b, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("Failed to read file %s: '%s'\n", *file, err)
	}
	client := &http.Client{}
	client.Timeout = time.Second * 5

	uri := fmt.Sprint(*endpoint, "/delete")
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
