# ![Infralight Logo](project-logo.png)

**Infralight Command Line Interface, Terminal User Interface and Client Library**

<!-- vim-markdown-toc GFM -->

* [Overview](#overview)
* [Features](#features)
* [Installation](#installation)
* [Quick Start](#quick-start)
    * [Example 1: Start the TUI](#example-1-start-the-tui)
    * [Example 2: List Available Environments](#example-2-list-available-environments)
* [Development](#development)
    * [Requirements](#requirements)
    * [Unit Tests and Static Code Analysis](#unit-tests-and-static-code-analysis)

<!-- vim-markdown-toc -->

## Overview

This repository contains a CLI, TUI and client library for the
[Infralight SaaS](https://infralight.co). Customers can use it to automate Infralight
operations in CI systems; manually execute such operations via the command line;
or integrate the client library into the customer's applications.

The client is written in Go and distributed as a single, statically-linked
executable.

![](screenshot.jpg)

## Features

- Beautiful view-only TUI to fetch information from the Infralight API. Only a
  subset of the features provided by the Infralight Dashboard are also included
  in the TUI.
- Comprehensive suite of CI-friendly commands to access the Infralight API.
- Go client library for integration with customer applications.

## Installation

TODO: once published, modify instructions to download the executable.

```sh
go clone git@github.com:infralight/cli.git
cd cli
CGO_ENABLED=0 \
    go build -a -tags netgo \
    -ldflags '-w -extldflags "-static"' \
    -o infralight \
    main.go
sudo install -Dm755 infralight /usr/local/bin/infralight
```

## Quick Start

```sh
infralight --help
```

You must provide an access key and secret key for the CLI to authenticate with
the Infralight API. Create a keypair through the Infralight dashboard.

The keypair can be provided via the `--access-key` and `--secret-key` command
line options. If not provided, the user will be prompted to enter them before
the CLI continues.

By default, the CLI will authenticate with Infralight's production API server.
To use a different server, provide the `--url` option.

If no command is provided, the program will start the Terminal User Interface.

### Example 1: Start the TUI

```sh
infralight
```

### Example 2: List Available Environments

```sh
infralight envs list
```

By default, output is one-lined JSON. To pretty print, add the `--pretty` flag.

## Development

During development, execute the code directly with `go run main.go`.

### Requirements

* [Go](https://golang.org/) v1.16+
* [golangci-lint](https://golangci-lint.run/) v1.35+

### Unit Tests and Static Code Analysis

To execute unit tests and static code analysis, run:

```sh
$ go build
$ go test ./...
$ golangci-lint run ./...
```
