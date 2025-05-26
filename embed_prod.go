//go:build prod
// +build prod

package main

import (
	"embed"
	"io/fs"
)

//go:embed static/js/*
//go:embed static/css/*
//go:embed static/img/*
//go:embed static/*
var embeddedFiles embed.FS

var StaticFiles fs.FS = embeddedFiles
