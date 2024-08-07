package main

import (
	"embed"
	"github.com/workflow-intelligence/gather/cmd"
	"github.com/workflow-intelligence/gather/server"
)

// It will add the specified files.
//go:embed index.html openapi.yaml
// It will add all non-hidden file in images, css, and js.
//go:embed audio fonts images img css js swagger

var Static embed.FS

func main() {
	server.Static = Static
	cmd.Execute()
}
