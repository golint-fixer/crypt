# Crypt

Crypt is a library for the [Go Programming Language][go]. It provides some
cryptographic operations.

## Status

[![Build Status](https://img.shields.io/travis/raiqub/crypt/master.svg?style=flat&label=linux%20build)](https://travis-ci.org/raiqub/crypt)
[![AppVeyor Build](https://img.shields.io/appveyor/ci/skarllot/crypt/master.svg?style=flat&label=windows%20build)](https://ci.appveyor.com/project/skarllot/crypt)
[![Coverage Status](https://coveralls.io/repos/raiqub/crypt/badge.svg?branch=master&service=github)](https://coveralls.io/github/raiqub/crypt?branch=master)
[![GoDoc](https://godoc.org/github.com/raiqub/crypt?status.svg)](http://godoc.org/github.com/raiqub/crypt)

## Features

 * **RandomAggr** type which provides an aggregated random data sources.
 * **Salter** type to create password salts and unique session IDs.
 * **SSTDEG** type which provides a System Sleep Time Delta Entropy Gathering.

## Installation

To install raiqub/crypt library run the following command:

~~~ bash
go get gopkg.in/raiqub/crypt.v0
~~~

To import this package, add the following line to your code:

~~~ bash
import "gopkg.in/raiqub/crypt.v0"
~~~

## Examples

Examples can be found on [library documentation][doc].

## Running tests

The tests can be run via the usual `go test` procedure:

~~~ bash
go test -v --race ./...
~~~

## License

raiqub/crypt is made available under the [Apache Version 2.0 License][license].


[go]: http://golang.org/
[doc]: http://godoc.org/github.com/raiqub/crypt
[license]: http://www.apache.org/licenses/LICENSE-2.0
