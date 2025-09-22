#!/usr/bin/env bats

load ../test_helper

@test "create, get, update and delete" {
  search_id="$(./dist/urlscan pro saved-search create -q "page.domain:example.com" -n "bats test" | jq -r ".search._id")"
  subscription_id="$(./dist/urlscan pro subscription create -s "$search_id" -f daily  -n "bats test" -e test@example.com  | jq -r ".subscription._id")"

  run ./dist/urlscan pro subscription get "$subscription_id"
  assert_success

  run ./dist/urlscan pro subscription update "$subscription_id" -f weekly -n "bats test" -e test@example.com
  assert_success

  run ./dist/urlscan pro subscription delete "$subscription_id"
  assert_success
}

