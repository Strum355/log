# Stroom Log™

[![Build Status](https://cloud.drone.io/api/badges/Strum355/log/status.svg)](https://cloud.drone.io/Strum355/log)
[![GoDoc](https://godoc.org/github.com/Strum355/log?status.svg)](https://godoc.org/github.com/Strum355/log)
[![Go report](https://goreportcard.com/badge/Strum355/log)](https://goreportcard.com/report/Strum355/log)

Simple logger inspired by [bwmarrin/lit](https://github.com/bwmarrin/lit) with fields, opentracing and json output support.

## Design Philosophy

The design of the API is inspired by Grafana's Loki log aggregation system and structured logging practices. As a result, it heavily favours using fields to log variable data and having log messages be the same regardless of the contextual data.

Example:

```go
log.WithFields(log.Fields{
    "userId": user.id,
    "requestId": requestId,
}).Info("user logged in successfully")
```

instead of

```go
log.Info("request %s user %d logged in successfully", user.id, requestId)
```

## Usage

Fetch the package:

```shell
go get github.com/Strum355/log
```

and import it:

```go
import (
    "github.com/Strum355/log"
)
```

initialize the logger for development:

```go
log.InitSimpleLogger(&log.Config{...})
```

or for production using a JSON log parser like FluentD

```go
log.InitJSONlogger(&log.Config{...})
```
