#!/usr/bin/env bats

load ../test_helper

# TODO: test store command

setup() {
  url="https://example.com"
  scanner_id="us01"
}

@test "scan, get result, dom, screenshot and response(purge afterwards)" {
  uuid="$(./dist/urlscan pro livescan scan "$url" -s "$scanner_id" | jq -r ".uuid")"

  run ./dist/urlscan pro livescan result "$uuid" -s "$scanner_id"
  assert_success
  result_output="$output"

  run ./dist/urlscan pro livescan dom "$uuid" -s "$scanner_id"
  assert_success
  assert [ -f "./${uuid}.html" ]
  rm -f "./${uuid}.html"

  run ./dist/urlscan pro livescan screenshot "$uuid" -s "$scanner_id"
  assert_success
  assert [ -f "./${uuid}.png" ]
  rm -f "./${uuid}.png"

  hash="$(echo "$result_output" | jq -r ".lists.hashes[0]")"
  run ./dist/urlscan pro livescan response "$hash" -s "$scanner_id"
  assert_success
  assert [ -f "./${hash}" ]

  run ./dist/urlscan pro livescan purge "$uuid" -s "$scanner_id"
  assert_success
}

@test "scanners" {
  run ./dist/urlscan pro livescan scanners
  assert_success
}
