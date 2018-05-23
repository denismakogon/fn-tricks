package main

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/fnproject/fdk-go"
	"github.com/sirupsen/logrus"
)

var DefaultRecursionDepth = WithDefault("DEFAULT_RECURSION_DEPTH", "100")

func getRecursionValue(fctx *fdk.Ctx) (int, error) {
	log := logrus.New()
	log.Infof("request URL: %s", fctx.RequestURL)
	u, err := url.Parse(fctx.RequestURL)
	if err != nil {
		return 0, err
	}

	q := u.Query()
	recursionDepth := q.Get(RecursionDepthKey)
	if recursionDepth == "" {
		recursionDepth = DefaultRecursionDepth
	}
	intRecursionDepth, err := strconv.Atoi(recursionDepth)
	if err != nil {
		return 0, err
	}

	if intRecursionDepth <= 0 {
		return 0, errors.New("invalid recursion depth")
	}

	return intRecursionDepth, nil
}

func setRecursionValue(fctx *fdk.Ctx, intRecursionDepth int) (*http.Request, error) {
	u, err := url.Parse(fctx.RequestURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set(RecursionDepthKey, strconv.Itoa(intRecursionDepth))
	u.RawQuery = q.Encode()

	return http.NewRequest(http.MethodPost, u.String(), nil)
}
