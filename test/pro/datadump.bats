#!/usr/bin/env bats

load ../test_helper

TODAY=$(date +%Y%m%d)

@test "list" {
  run bash -c "./dist/urlscan pro datadump list hours/api/$TODAY | jq ."
  assert_success
}

@test "download" {
  # NOTE: this can be flaky depending on the time of the day the test is run...
  run bash -c "./dist/urlscan pro datadump download hours/search/$TODAY/$TODAY-00.gz --output /tmp/$TODAY.gz "
  assert_success
}

