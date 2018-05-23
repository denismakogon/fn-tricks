package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

func DoRequest(req *http.Request, httpClient *http.Client) error {
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode > http.StatusAccepted {
		return fmt.Errorf("unable to submit webhoot successfully, "+
			"status code: %d, response body: '%s'", resp.StatusCode, string(b))
	}
	logrus.Infof("request accepted, response: %s", string(b))

	return nil
}
