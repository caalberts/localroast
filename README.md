# Localroast

[![CircleCI](https://circleci.com/gh/caalberts/localroast/tree/master.svg?style=svg)](https://circleci.com/gh/caalberts/localroast/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/caalberts/localroast)](https://goreportcard.com/report/github.com/caalberts/localroast)

## Overview

Localroast is a Go command line tool to run a local stub http service. Localroast takes command line string arguments to define routes and stub responses with http status codes.

## Installation

```sh
go get -u github.com/caalberts/localroast
```

## Usage

Routes are defined in the format `'<METHOD> <PATH> <STATUS_CODE>'`, for example `'GET / 200'`. Multiple route definitions are created using successive string arguments. Only `GET` and `POST` are currently supported.

To start Localroast:
```sh
localroast \
  'GET / 200' \
  'POST /users 201'
```
