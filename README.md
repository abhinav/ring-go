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

## License

This software is made available under the MIT license.
