#!/usr/bin/env bats

load test_helper

@test "scan response" {
  # hash of an empty content/string
  run bash -c "./dist/urlscan scan response -f e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  assert_success
}

teardown() {
  rm -f "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
}
