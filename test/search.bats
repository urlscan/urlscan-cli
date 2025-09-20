#!/usr/bin/env bats

load test_helper

@test "search with limit" {
  run bash -c "./dist/urlscan search 'page.domain:example.com' --limit 1 | jq -r '.results | length'"
  assert_output 1
}

@test "search with limit & size" {
  run bash -c "./dist/urlscan search 'page.domain:example.com' --limit 10 --size 5 | jq -r '.results | length'"
  assert_output 10
}

@test "search with all" {
  run ./dist/urlscan search 'page.domain:example.com AND date:>now-1d' --all
  assert_success
}
