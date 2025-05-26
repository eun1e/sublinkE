//go:build !prod
// +build !prod

package main

import "io/fs"

var StaticFiles fs.FS = nil
