#!/bin/bash

GOOS=linux GOARCH=amd64 go build .
gcloud compute scp redirector vm-1:~/
