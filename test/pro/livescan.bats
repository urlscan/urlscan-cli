#!/usr/bin/env bats

load ../test_helper

# TODO: test store command

setup() {
  url="https://example.com"
  scanner_id="us01"
}

@test "scan, get result and dom and purge" {
  uuid="$(./dist/urlscan pro livescan scan "$url" -s "$scanner_id" | jq -r ".uuid")"

  run ./dist/urlscan pro livescan result "$uuid"
  assert_success

  run ./dist/urlscan pro livescan dom "$uuid"
  assert_success

  run ./dist/urlscan pro livescan purge "$uuid"
  assert_success
}

@test "scanners" {
  run ./dist/urlscan pro livescan scanners
  assert_success
}
