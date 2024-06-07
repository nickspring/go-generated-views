# go-generated-views

[![Actions Status](https://github.com/nickspring/go-generated-views/actions/workflows/build_and_test.yml/badge.svg)](https://github.com/nickspring/go-generated-views/actions/workflows/build_and_test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nickspring/go-generated-views)](https://goreportcard.com/report/github.com/nickspring/go-generated-views)
[![GoDoc](https://godoc.org/github.com/nickspring/go-generated-views?status.svg)](https://godoc.org/github.com/nickspring/go-generated-views)

Go Generated views is a view generator for Go - generating views with different tags using one base model struct.

## How it works

go-generated-views will analyze tags declarations like this:

```go
//go:generate go-generated-views

package example

type MyStruct struct {
	Field1 string `json:"usual_json_field" json(Special,VerySpecial):"special_json_field"`
} 
```

The most important part here is `json(Special,VerySpecial)`. It declares two additional views for this struct:
* Special
* VerySpecial
These views will have custom `json` tag with value `special_json_field`.

You could declare what you want (`json`, `form`, `db`, `required`, `validation` tags etc.) in any combination.

After processing go-generated-views generate a file (usually with suffix `_view.go`) with new views definitions and methods.
In this case:

```go
package example

// MyStructSpecialView represents view for model MyStruct
type MyStructSpecialView struct {
    Field1 string `json:"special_json_field"`
}

// MyStructVerySpecialView represents view for model MyStruct
type MyStructVerySpecialView struct {
    Field1 string `json:"special_json_field"`
}

// ToModel change pointer from MyStructSpecialView to MyStruct
func (view *MyStructSpecialView) ToModel() *MyStruct {
    return (*MyStruct)(view)
}

// NewMyStructSpecialView change pointer from MyStruct to MyStructSpecialView
func NewMyStructSpecialView(model *MyStruct) *MyStructSpecialView {
    return (*MyStructSpecialView)(model)
}

// ToModel change pointer from MyStructVerySpecialView to MyStruct
func (view *MyStructVerySpecialView) ToModel() *MyStruct {
    return (*MyStruct)(view)
}

// NewMyStructVerySpecialView change pointer from MyStruct to MyStructVerySpecialView
func NewMyStructVerySpecialView(model *MyStruct) *MyStructVerySpecialView {
    return (*MyStructVerySpecialView)(model)
}
```

As you see we have two auto-generated views with additional methods / receivers to quick conversion from / to base model.

See also examples in `example` folder.

## Goal

The goal of go-generated-view is to create an easy-to-use view auto-generator / manager to reduce structs & tags duplicity in Go code.
Typical and simple example - you need two a little bit of different structs to add object and edit object (first one doesn't require ID of new object, but second one does). 

## Docker image

You can now use a docker image directly for running the command if you do not wish to install anything!

```shell
   docker run -w /app -v $(pwd):/app nickspring/go-generated-views:$(GO_GENERATED_VIEWS_VERSION)
```

## Installation

You can now download a release directly from Github and use that for generating your views! (Thanks to [GoReleaser](https://github.com/goreleaser/goreleaser-action))

```shell
    curl -fsSL "https://github.com/nickspring/go-generated-views/releases/download/v$(GO_GENERATED_VIEWS_VERSION)/go-generated-views_$(uname -s)_$(uname -m)" -o go-generated-views && chmod +x go-generated-views
```

## Adding it to your project

### Using go generate

1. Add a go:generate line to your file like so... `//go:generate go-generated-views`
2. Run go generate like so `go generate ./...`
3. Enjoy your newly created view!

## Command options

``` shell
go-generated-views --help

NAME:
   go-generated-views - A struct views generator for go

USAGE:
   go-generated-views [global options] 

GLOBAL OPTIONS:
   --file value, -f value [ --file value, -f value ]          The file(s) to generate views. Use more than one flag for more files. [$GOFILE]
   --buildtag value, -b value [ --buildtag value, -b value ]  Adds build tags to a generated view file.
   --output-suffix .go                                        Changes the default filename suffix of _view to something else.  .go will be appended to the end of the string no matter what, so that `_test.go` cases can be accommodated
   --help, -h                                                 show help
   --version, -v                                              print the version
```

## Copyright

This project was inspired by [go-enum](https://github.com/nickspring/go-generated-views) and [retag](https://github.com/sevlyar/retag) libraries.
(c) Nikolay Yarovoy. All rights reserved.
