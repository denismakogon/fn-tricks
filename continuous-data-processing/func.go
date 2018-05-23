package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/fnproject/fdk-go"
	"github.com/sirupsen/logrus"
)

var (
	ReverseMode = os.Getenv("REVERSE_MODE")
)

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

func myHandler(ctx context.Context, in io.Reader, _ io.Writer) error {
	fctx := fdk.Context(ctx)

	var data ContinuousData
	err := json.NewDecoder(in).Decode(&data)
	if err != nil {
		return err
	}
	total := len(data.DataChunks)
	if total == 0 {
		return nil
	}

	intReverseMode, err := strconv.Atoi(ReverseMode)
	if err != nil {
		return err
	}

	var currentChunk interface{}
	if intReverseMode == 1 {
		// cut off the last item
		currentChunk = data.DataChunks[total-1]
	} else {
		// cut off the first item
		currentChunk = data.DataChunks[0]
	}

	if total-1 > 0 {
		reduceAndCallNext(fctx, &data, intReverseMode)
	}

	return processData(currentChunk)
}

func processData(d interface{}) error {
	logrus.Infof("data: %v", d)
	return nil
}
