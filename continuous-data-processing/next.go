package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/fnproject/fdk-go"
)

func reduceAndCallNext(fctx *fdk.Ctx, data *ContinuousData, reverseMode int) error {
	var b bytes.Buffer
	if reverseMode == 1 {
		data.DataChunks = data.DataChunks[0 : len(data.DataChunks)-2]
	} else {
		data.DataChunks = data.DataChunks[1:]
	}

	json.NewEncoder(&b).Encode(data)

	req, err := http.NewRequest(http.MethodPost, fctx.RequestURL, &b)
	if err != nil {
		return err
	}
	return DoRequest(req, http.DefaultClient)
}
