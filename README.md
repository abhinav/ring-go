# container/ring

[![Go Reference](https://pkg.go.dev/badge/go.abhg.dev/container/ring.svg)](https://pkg.go.dev/go.abhg.dev/container/ring)
[![CI](https://github.com/abhinav/ring-go/actions/workflows/ci.yml/badge.svg)](https://github.com/abhinav/ring-go/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/abhinav/ring-go/branch/main/graph/badge.svg?token=zXfHANxPoF)](https://codecov.io/gh/abhinav/ring-go)

container/ring provides an efficient queue data structure
intended to be used in performance-sensitive contexts
with usage patterns that have:

- a similar rate of pushes and pops
- bursts of pushes followed by bursts of pops

See [API Reference](https://abhinav.github.io/ring-go) for more details.

## Motivation

I often find myself needing a queue in projects.
In a hurry, I reach for `container/list` or a slice.
However, they are both allocation-heavy (depending on usage pattern).

This repository largely exists so that there's a known working version
of a more efficient queue implementation (for specific usage patterns)
that I can use or copy-paste directly where I need.
Feel free to use it in the same way if it meets your needs.

## License

This software is made available under the MIT license.
