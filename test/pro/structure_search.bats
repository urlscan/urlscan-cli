#!/usr/bin/env bats

load ../test_helper

setup() {
  # ref. https://docs.urlscan.io/apis/urlscan-openapi/search/similarsearch
  uuid="68e26c59-2eae-437b-aeb1-cf750fafe7d7"
}

@test "structure-search with limit" {
  run bash -c "./dist/urlscan pro structure-search $uuid --limit 10 | jq -r '.results | length'"
  assert_output 10
}

@test "structure-search with limit & size" {
  run bash -c "./dist/urlscan pro structure-search $uuid --limit 20 --size 10 | jq -r '.results | length'"
  assert_output 20
}


