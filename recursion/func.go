package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/fnproject/fdk-go"
	"github.com/sirupsen/logrus"
)

func WithDefault(key, value string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	return value
}

var (
	BackOff           = WithDefault("BACKOFF_TIMEOUT", "5")
	RecursionDepthKey = "recursion_depth"
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
	logrus.Printf("request accepted, response: %s", string(b))

	return nil
}

func main() {
	fdk.Handle(fdk.HandlerFunc(withError))
}

func withError(ctx context.Context, in io.Reader, out io.Writer) {
	err := myHandler(ctx, in, out)
	if err != nil {
		fdk.WriteStatus(out, http.StatusInternalServerError)
		out.Write([]byte(err.Error()))
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	fdk.WriteStatus(out, http.StatusAccepted)
}

type Payload struct {
	RequirementMet  bool                   `json:"requirementMet"`
	IncomingPayload map[string]interface{} `json:"incomingPayload"`
}

func myHandler(ctx context.Context, in io.Reader, _ io.Writer) error {
	log := logrus.New()
	fctx := fdk.Context(ctx)

	t, err := strconv.Atoi(BackOff)
	if err != nil {
		return err
	}

	intRecursionDepth, err := getRecursionValue(fctx)
	if err != nil {
		return err
	}

	var p Payload
	err = json.NewDecoder(in).Decode(&p)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	if !p.RequirementMet {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			var b bytes.Buffer
			p.RequirementMet = NewRandBool()
			json.NewEncoder(&b).Encode(p)

			log.Infof("starting next recursion at level: %d", intRecursionDepth)
			req, err := setRecursionValue(fctx, intRecursionDepth-1)
			req.Body = ioutil.NopCloser(&b)

			err = DoRequest(req, http.DefaultClient)
			if err != nil {
				log.Fatal(err.Error())
			}

			time.Sleep(time.Duration(t) * time.Second)
		}(&wg)
	} else {
		log.Info("requirement matched, aborting recursive function...")
		return nil
	}

	wg.Wait()

	return nil
}
