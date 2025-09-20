#!/usr/bin/env bats

load ../test_helper

setup() {
  # zipped EICAR file
  hash="275a021bbfb6489e54d471899f7db9d1663fc695ec2fe2a2c4538aabf651fd0f"
}

teardown() {
  [ -f "${hash}.zip" ] && rm -f "${hash}.zip"
}

@test "file" {
  run ./dist/urlscan pro file $hash
  assert_success
  assert [ -f "${hash}.zip" ]
}


