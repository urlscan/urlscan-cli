#!/usr/bin/env bats

load ../test_helper

@test "search with limit" {
  run bash -c "./dist/urlscan pro hostname example.com --limit 1 | jq -r '.results | length'"
  assert_output 1
}

@test "search with limit & size" {
  run bash -c "./dist/urlscan pro hostname example.com --limit 10 --size 10 | jq -r '.results | length'"
  assert_output 10
}

@test "search with --params" {
  run bash -c "./dist/urlscan pro hostname example.com --params '{\"limit\":\"10\"}' | jq -r '.results | length'"
  assert_output 10
}

