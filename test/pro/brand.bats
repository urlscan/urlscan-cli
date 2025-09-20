#!/usr/bin/env bats

load ../test_helper

@test "available" {
  run ./dist/urlscan pro brand available
  assert_success
}

@test "list" {
  run ./dist/urlscan pro brand list
  assert_success
}

