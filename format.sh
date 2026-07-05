#!/bin/bash

gofmt -w $(find . -name '*.go' -not -path './tmp/*' -not -path './vendor/*')