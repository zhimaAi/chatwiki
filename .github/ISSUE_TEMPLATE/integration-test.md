name: Integration Test
description: Report integration testing issues
title: "[Integration Test] "
labels: integration-test
body:
  - type: textarea
    id: description
    attributes:
      label: Description
      description: Describe the integration issue
    validations:
      required: true
  - type: textarea
    id: steps
    attributes:
      label: Steps to Reproduce
      description: How to reproduce this issue
    validations:
      required: true
  - type: textarea
    id: expected
    attributes:
      label: Expected Behavior
    validations:
      required: true
  - type: textarea
    id: actual
    attributes:
      label: Actual Behavior
    validations:
      required: true
