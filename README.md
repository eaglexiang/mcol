# mcol

[![GoDoc](https://godoc.org/github.com/eaglexiang/mcol?status.svg)](https://godoc.org/github.com/eaglexiang/mcol)
![build](https://api.travis-ci.com/eaglexiang/mcol.svg?branch=master)
[![codebeat badge](https://codebeat.co/badges/b4703e04-0239-4c04-86fd-d4a4623e5470)](https://codebeat.co/projects/github-com-eaglexiang-mcol-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/eaglexiang/mcol)](https://goreportcard.com/report/github.com/eaglexiang/mcol)
![license](https://img.shields.io/badge/license-MIT-black.svg)

easy to to search mongo col

## install

```shell
go get github.com/eaglexiang/mcol
```

## config

create or edit `$HOME/.mcol.confg` with JSON.

sample:

```json
{
    "addr": "127.0.0.1:27017",
    "db": "admin",
    "username": "admin",
    "password": "123456"
}
```

## init cache

run command:

```shell
mcol --cache
```

## search

run command:

```shell
mcol $key0 $key1
```

sample:

```shell
mcol db2 col1
```

output:

![search](doc/mcol_search.png)

> keys found will be highlighted.
