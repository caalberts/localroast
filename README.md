# Localroast

[![CircleCI](https://circleci.com/gh/caalberts/localroast/tree/master.svg?style=svg)](https://circleci.com/gh/caalberts/localroast/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/caalberts/localroast)](https://goreportcard.com/report/github.com/caalberts/localroast)

![localroast](coffee.png)

## Overview

Localroast is a lightweight Go command line tool to run a local stub http service. The endpoints and the stubbed responses are defined in a JSON file.

## Installation

```sh
go get -u github.com/caalberts/localroast
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
- [ ] yml input
- [ ] autoload file changes
