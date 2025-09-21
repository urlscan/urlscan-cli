#!/usr/bin/env bats

load ../test_helper

# NOTE: don't test create command to avoid unnecessary channel creation

setup() {
  channel_id=${CHANNEL_ID}
  if [ -z "$channel_id" ]; then
    echo "Make sure to set CHANNEL_ID for testing the channel commands"
    exit 1
  fi
}

@test "get" {
  run ./dist/urlscan pro channel get $channel_id
  assert_success
}

@test "update" {
  name=$(./dist/urlscan pro channel get $channel_id | jq -r '.channel.name')

  run ./dist/urlscan pro channel update $channel_id -n "bats-test"
  assert_success

  updated=$(./dist/urlscan pro channel get $channel_id | jq -r '.channel.name')
  assert_equal "bats-test" "$updated"

  # revert the change
  run ./dist/urlscan pro channel update $channel_id -n "$name"
  assert_success
}

@test "list" {
  run ./dist/urlscan pro channel list
  assert_success
}
