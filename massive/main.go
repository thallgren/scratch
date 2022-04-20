package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/datawire/dlib/dexec"
	"github.com/datawire/dlib/dlog"
)

func main() {
	if err := doIt(context.Background()); err != nil {
		log.Fatal(err)
	}
}

const purposeLabel = "massive-cluster-testing"

const svcTemplate = `# The echo-double service exposes two ports, 80 and 81 and directs them to two separate containers
---
apiVersion: v1
kind: Service
metadata:
  name: echo-%02[1]d
spec:
  type: ClusterIP
  selector:
    service: echo-%02[1]d
  ports:
    - name: http
      port: 80
      targetPort: http
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: echo-%02[1]d
  labels:
    service: echo-%02[1]d
spec:
  replicas: 1
  selector:
    matchLabels:
      service: echo-%02[1]d
  template:
    metadata:
      labels:
        service: echo-%02[1]d
      annotations:
        telepresence.getambassador.io/inject-traffic-agent: enabled
    spec:
      containers:
        - name: echo
          image: multi:5000/tel2/echo
          ports:
            - containerPort: 8080
              name: http
          env:
            - name: LOG_LEVEL
              value: debug
            - name: FOO_SOME_LONG_NAME
              value: Just a placeholder with text
            - name: BAR_SOME_LONG_NAME
              value: Just another placeholder with text
            - name: FEE_SOME_LONG_NAME
              value: Just a third placeholder with text
            - name: FUM_SOME_LONG_NAME
              value: And then again, another placeholder with text
          resources:
            limits:
              cpu: 10m
              memory: 64Mi
`

func doIt(ctx context.Context) error {
	for i := 0; i < 10; i++ {
		ns := fmt.Sprintf("the-%02d-namespace", i)
		if err := kubectl(ctx, "create", "namespace", ns); err != nil {
			return err
		}
		if err := kubectl(ctx, "label", "namespace", ns, "purpose="+purposeLabel, fmt.Sprintf("app.kubernetes.io/name=%s", ns)); err != nil {
			return err
		}

		retries := 0
		for w := 0; w < 8; w++ {
			cmd := dexec.CommandContext(ctx, "kubectl", "--namespace", ns, "apply", "-f", "-")
			cmd.Stdin = bytes.NewReader([]byte(fmt.Sprintf(svcTemplate, w)))
			cmd.DisableLogging = true
			if err := cmd.Run(); err != nil {
				return err
			}
			n := fmt.Sprintf("echo-%02d", w)
			if err := kubectl(ctx, "--namespace", ns, "rollout", "status", "deployment", "--timeout", "8s", n); err != nil {
				w--
				retries++
				if retries > 3 {
					return err
				}
				dlog.Infof(ctx, "retrying deploy %s", n)
				if err = kubectl(ctx, "--namespace", ns, "delete", "svc,deploy", n); err != nil {
					return err
				}
			} else {
				retries = 0
			}
		}
	}
	return nil
}

func kubectl(ctx context.Context, args ...string) error {
	return dexec.CommandContext(ctx, "kubectl", args...).Run()
}
