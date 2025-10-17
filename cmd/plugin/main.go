package main

import (
	"log"

	"github.com/dewan-ahmed/drone-openai-plugin/pkg/plugin"
)

func main() {
	if err := plugin.Run(); err != nil {
		log.Fatalf("Plugin execution failed: %v", err)
	}
}

