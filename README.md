# Localroast

[![CircleCI](https://circleci.com/gh/caalberts/localroast/tree/master.svg?style=svg)](https://circleci.com/gh/caalberts/localroast/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/caalberts/localroast)](https://goreportcard.com/report/github.com/caalberts/localroast)

## Overview

Localroast is a lightweight Go command line tool to run a local stub http service. Stub endpoints can be defined either using a JSON file input or command line arguments.

## Installation

```sh
go get -u github.com/caalberts/localroast
```

## Usage

### JSON file

To start a local stub server using a JSON file, use the `-json` flag, followed by the path to the JSON file.

```sh
localroast -json examples/stubs.json
```

The JSON file must be a JSON array containing endpoint definitions. Each endpoint is represented as a JSON object with keys `method`, `path`, `status` and `response`. A valid JSON file would look like this:
```json
[
  {
    "method": "GET",
    "path": "/",
    "status": 200,
    "response": {
        "success": true
    }
  }
]
```

See [examples/stubs.json](examples/stubs.json) for more.


### Command line

Endpoints are defined in the format `'<METHOD> <PATH> <STATUS_CODE>'`, for example `'GET / 200'`. Multiple endpoint definitions are created using successive string arguments.

To start a local stub server using CLI arguments:
```sh
localroast \
  'GET / 200' \
  'POST /users 201'
```
