#!/usr/bin/env bats

load test_helper

@test "quotas" {
  run bash -c "./dist/urlscan quotas | jq -r '.source'"
  assert_output --regexp '(^team$|^user$)'
}
