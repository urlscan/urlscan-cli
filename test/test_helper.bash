#!/usr/bin/env bash
BATS_LOCATION=`which bats`;
BATS_LOCATION="${BATS_LOCATION%/*/*/*/*}"

load "$BATS_LOCATION/bats-assert/load.bash"
load "$BATS_LOCATION/bats-support/load.bash"
