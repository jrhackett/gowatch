# GoWatch

[![Build Status](https://travis-ci.org/jrhackett/gowatch.svg?branch=master&service=github)](https://travis-ci.org/jrhackett/gowatch)
[![Coverage Status](https://coveralls.io/repos/github/jrhackett/gowatch/badge.svg?branch=master&service=github)](https://coveralls.io/github/jrhackett/gowatch?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/jrhackett/gowatch)](https://goreportcard.com/report/github.com/jrhackett/gowatch)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/jrhackett/gowatch.svg)](https://github.com/jrhackett/gowatch/releases/latest)

This is a development tool that automatically runs tests when files are changed in a Go project.

## Installation

1. `go get github.com/jrhackett/gowatch/cmd/gowatch`
2. `go install github.com/jrhackett/gowatch/cmd/gowatch` 

## Usage

`gowatch -command="go test ./..." -path="./"`

You can supply an arbitrary path and command to watch any directory and run the supplied command when files are changed within that path.
