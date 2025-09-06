#!/usr/bin/env bats

load test_helper

@test "user" {
  run bash -c "./dist/urlscan quotas | jq -r '.username'"
  assert_output --regexp '(^.+$)'
}
