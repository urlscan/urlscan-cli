#!/usr/bin/env bats

load ../test_helper

@test "list" {
  run bash -c "./dist/urlscan pro datadump list --time-window days --file-type api --date 20260101 | jq ."
  assert_success
}

@test "download" {
  run bash -c "./dist/urlscan pro datadump download hours/search/20260101/20260101-00.gz --output /tmp/20260101.gz "
  assert_success
}

