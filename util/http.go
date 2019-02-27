package util

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

func GetRequest(url string) (bytes []byte, err error) {

	client := &http.Client{}
	client.Timeout = time.Second * 5

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logrus.Error(err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer resp.Body.Close()
	bytes, err = ioutil.ReadAll(resp.Body)
	return

}
