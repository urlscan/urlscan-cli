#!/usr/bin/env bash
type brew &>/dev/null && export BATS_LIB_PATH="${BATS_LIB_PATH}:$(brew --prefix)/lib"

bats_load_library bats-support
bats_load_library bats-assert
