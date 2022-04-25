# zaprec

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
![Go](https://github.com/lefinal/zaprec/workflows/Go/badge.svg?branch=main)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lefinal/zaprec)
[![GoReportCard example](https://goreportcard.com/badge/github.com/lefinal/zaprec)](https://goreportcard.com/report/github.com/lefinal/zaprec)
[![codecov](https://codecov.io/gh/lefinal/zaprec/branch/main/graph/badge.svg?token=ema8Z2HEk5)](https://codecov.io/gh/lefinal/zaprec)
[![GitHub issues](https://img.shields.io/github/issues/lefinal/zaprec)](https://github.com/lefinal/zaprec/issues)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/lefinal/zaprec)

Recording of log-calls to `zap.Logger` (see [uber-go/zap](https://github.com/uber-go/zap)). This package is intended to be used for easier testing of logging with zap.

# Installation

In order to use this package, run:

```shell
go get github.com/lefinal/zaprec
```

# Usage

Create a new recorder via:

```go
logger, rec := NewRecorder(zapcore.InfoLevel)
```

You can now use `logger` as regular `*zap.Logger` and check the returned record store for records with methods like `Records()` or `RecordsByLevel(level zapcore.Level)`.