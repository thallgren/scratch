package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/datawire/dlib/dlog"
	"github.com/datawire/dtest"
)

func main() {
	lgr := logrus.StandardLogger()
	lgr.SetLevel(logrus.ErrorLevel)
	ctx := dlog.WithLogger(context.Background(), dlog.WrapLogrus(lgr))
	fmt.Println(dtest.Kubeconfig(ctx))
}
