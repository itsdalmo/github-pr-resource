#!/bin/sh

set -eu

apk add --update --no-cache git-lfs
git lfs install
