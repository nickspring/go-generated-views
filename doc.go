// Copyright (c) 2024 Nikolay Yarovoy

// go-generated-views is a utility for generating views in go (using one base model).
//
// The generator looks for view names between `(` and `)` in tags.
// It creates alternative views with own tags (based on model tags) with methods to convert from/to base model.
// This is a useful when you have one model but need to have different tags validation rules on input and output sides.
//
// Installation
//
//	go get github.com/nickspring/go-generated-views
//
// Usage:
// Sample File
//
// package example
//
//	type Book struct {
//		 ID          uint64 `binding(get):"required,gt=0" json:"id" json(add):"-,omitempty"`
//		 Name        string `json:"name" binding(add):"required"`
//		 Description string `json:"description"`
//	}
//
// Command to generate views
//
//	go generate ./
package main
