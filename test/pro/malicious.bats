#!/usr/bin/env bats

load ../test_helper

@test "lookup ip" {
  run bash -c "./dist/urlscan pro malicious lookup ip 1.1.1.1 | jq -r '.type'"
  assert_output "ip"
}

@test "lookup hostname" {
  run bash -c "./dist/urlscan pro malicious lookup hostname www.example.com | jq -r '.type'"
  assert_output "hostname"
}

@test "lookup domain" {
  run bash -c "./dist/urlscan pro malicious lookup domain example.com | jq -r '.type'"
  assert_output "domain"
}

@test "lookup url" {
  run bash -c "./dist/urlscan pro malicious lookup url 'https://example.com' | jq -r '.type'"
  assert_output "url"
}

@test "lookup with invalid type" {
  run ./dist/urlscan pro malicious lookup invalid 1.1.1.1
  assert_failure
}

