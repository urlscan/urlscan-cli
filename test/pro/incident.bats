#!/usr/bin/env bats

load ../test_helper

setup() {
  expire_after=3600
}

@test "create, get, states, update, close and restart" {
  incident_id="$(./dist/urlscan pro incident create -o "example.com" --expire-after "$expire_after" | jq -r ".incident._id")"

  run ./dist/urlscan pro incident get "$incident_id"
  assert_success

  run ./dist/urlscan pro incident states "$incident_id"
  assert_success

  run ./dist/urlscan pro incident update "$incident_id" -o "example.net"
  assert_success

  run ./dist/urlscan pro incident close "$incident_id"
  assert_success

  run ./dist/urlscan pro incident restart "$incident_id"
  assert_success

  # close it again to keep things clean
  run ./dist/urlscan pro incident close "$incident_id"
  assert_success
}


@test "copy" {
  incident_id="$(./dist/urlscan pro incident create -o "example.com" --expire-after "$expire_after" | jq -r ".incident._id")"
  copied_id="$(./dist/urlscan pro incident copy "$incident_id" | jq -r ".incidents._id")"

  run ./dist/urlscan pro incident get "$copied_id"
  assert_success

  # clean up
  run ./dist/urlscan pro incident close "$incident_id"
  assert_success

  run ./dist/urlscan pro incident close "$copied_id"
  assert_success
}

@test "fork" {
  incident_id="$(./dist/urlscan pro incident create -o "example.com" --expire-after "$expire_after" | jq -r ".incident._id")"
  forked_id="$(./dist/urlscan pro incident fork "$incident_id" | jq -r ".incidents._id")"

  run ./dist/urlscan pro incident get "$forked_id"
  assert_success

  # clean up
  run ./dist/urlscan pro incident close "$incident_id"
  assert_success

  run ./dist/urlscan pro incident close "$forked_id"
  assert_success
}
