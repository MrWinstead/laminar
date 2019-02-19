# Laminar [![GoDoc](https://godoc.org/github.com/mrwinstead/laminar?status.svg)](https://godoc.org/github.com/mrwinstead/laminar) [![Go Report Card](https://goreportcard.com/badge/github.com/mrwinstead/laminar)](https://goreportcard.com/report/github.com/mrwinstead/laminar)

**Laminar** is a middlewares library for [Goka](https://github.com/lovoo/goka)
heavily inspired by the certain functionality present in
[Go-Kit](https://github.com/go-kit/kit). Using Laminar, engineers familiar with
go-kit should feel at home wriing Goka processors.

## Motivation

Go has an excellent ecosystem and tradition of middlewares for HTTP processing.
While a streaming application built with Goka is certainly not a web
application, the patterns common there can be adapted and re-used in order to
get over the initial hurdle of building stream processing applications.
