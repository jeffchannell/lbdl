#!/usr/bin/env bash
set -e
docker build -t lbdl .
docker run --rm -v $(pwd):/go/src/app lbdl make $@
