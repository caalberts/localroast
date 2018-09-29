# Localroast

[![CircleCI](https://circleci.com/gh/caalberts/localroast/tree/master.svg?style=svg)](https://circleci.com/gh/caalberts/localroast/tree/master)
[![Coverage Status](https://coveralls.io/repos/github/caalberts/localroast/badge.svg?branch=master)](https://coveralls.io/github/caalberts/localroast?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/caalberts/localroast)](https://goreportcard.com/report/github.com/caalberts/localroast)

![localroast](coffee.png)

## Overview

Localroast quickly stubs a HTTP server. The stubbed endpoints and responses are defined in a JSON file.

## Installation

From brew:

```sh
brew install caalberts/tap/localroast
```

From source:

```sh
go get -u github.com/caalberts/localroast/cmd/localroast
```

## Usage

```sh
localroast examples/stubs.json
```

The command takes a single argument, a path to a JSON file. The JSON file must be a JSON array containing endpoint definitions. Each endpoint is represented as a JSON object with keys `method`, `path`, `status` and `response`. `response` can be any valid JSON object.
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

## Features

- [x] json input
- [x] path variable
- [x] autoload file changes
- [ ] yml input
