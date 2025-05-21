# go-gin-static

[![tag](https://img.shields.io/github/tag/cyg-pd/go-gin-static.svg)](https://github.com/cyg-pd/go-gin-static/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-%23007d9c)
[![GoDoc](https://godoc.org/github.com/cyg-pd/go-gin-static?status.svg)](https://pkg.go.dev/github.com/cyg-pd/go-gin-static)
![Build Status](https://github.com/cyg-pd/go-gin-static/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/cyg-pd/go-gin-static)](https://goreportcard.com/report/github.com/cyg-pd/go-gin-static)
[![Coverage](https://img.shields.io/codecov/c/github/cyg-pd/go-gin-static)](https://codecov.io/gh/cyg-pd/go-gin-static)
[![Contributors](https://img.shields.io/github/contributors/cyg-pd/go-gin-static)](https://github.com/cyg-pd/go-gin-static/graphs/contributors)
[![License](https://img.shields.io/github/license/cyg-pd/go-gin-static)](./LICENSE)

## 🚀 Install

```sh
go get github.com/cyg-pd/go-gin-static@v1
```

This library is v1 and follows SemVer strictly.

No breaking changes will be made to exported APIs before v2.0.0.

This library has no dependencies outside the Go standard library.

## 💡 Usage

You can import `go-gin-static` using:

```go
import (
	"embed"
	"log/slog"
	"net/http"

	"github.com/cyg-pd/go-gin-static"
	"github.com/gin-gonic/gin"
)

//go:embed all:dist/*
var app embed.FS

func main() {
	r := gin.New()
	h := http.Header{"Cache-Control": {"public", "max-age=86400"}}
	s := static.New(
		app,
		static.WithTemplateValuesFunc(func(req *http.Request) any {
			return map[string]any{
				"baseURL": baseURLParser(req),
				"gtm":     gtmParser(req),
			}
		}),
		static.WithRenderHeader(h),
	)
	r.NoRoute(s.Handler())
	r.Run()
}
```
