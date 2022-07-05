#!/usr/bin/env bash

go build -o pcfg_tool main.go

GOARCH=amd64 GOOS=darwin go build -o pcfg_tool-amd64-darwin main.go
GOARCH=amd64 GOOS=linux go build -o pcfg_tool-amd64-linux main.go