package main

import "github.com/dwirx/ghex/cmd/ghex/commands"

// Version is set during build via ldflags
var Version = "1.0.0"

func main() {
	commands.Version = Version
	commands.Execute()
}
