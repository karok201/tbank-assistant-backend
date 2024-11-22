#!/bin/bash

set -x
set -e

exec go run /app/cmd/main.go
