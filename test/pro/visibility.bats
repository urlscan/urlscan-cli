#!/usr/bin/env bats

load ../test_helper

@test "scan, change and reset visibility" {
  scan_id="$(./dist/urlscan scan submit "https://http-test.com/" -v private --wait | jq -r ".task.uuid")"

  message="$(./dist/urlscan pro visibility update "$scan_id" -v unlisted | jq -r ".message")"
  assert_equal "Visibility updated" "$message"


  message="$(./dist/urlscan pro visibility reset "$scan_id" | jq -r ".message")"
  assert_equal "Visibility reset" "$message"
}

