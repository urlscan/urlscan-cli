#!/usr/bin/env bats

load test_helper

@test "version" {
  run ./dist/urlscan version
  assert_line "urlscan-cli $VERSION"
}
