package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/dewan-ahmed/drone-openai-plugin/pkg/plugin"
)

func main() {
	// Print startup banner for debugging
	fmt.Printf("===========================================\n")
	fmt.Printf("Drone OpenAI Plugin v0.1.2\n")
	fmt.Printf("Built for: linux/amd64\n")
	fmt.Printf("Running on: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("===========================================\n\n")

	// Verify we're running in the expected environment
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		fmt.Printf("⚠️  WARNING: This binary was built for linux/amd64 but is running on %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Printf("This may cause compatibility issues.\n\n")
	}

	if err := plugin.Run(); err != nil {
		log.Fatalf("❌ Plugin execution failed: %v", err)
		os.Exit(1)
	}
}

