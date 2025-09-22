#!/usr/bin/env bats

load ../test_helper

@test "create, get, update and delete" {
  search_id="$(./dist/urlscan pro saved-search create -q "page.domain:example.com" -n "bats test" | jq -r ".search._id")"

  run ./dist/urlscan pro saved-search get "$search_id"
  assert_success

  run ./dist/urlscan pro saved-search update "$search_id" -q "page.domain:example.net"
  assert_success

  run ./dist/urlscan pro saved-search delete "$search_id"
  assert_success
}

@test "list" {
  run ./dist/urlscan pro saved-search list
  assert_success
}
