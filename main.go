package main

import (
	"os"
	"strings"
)

func main() {
	targetNamespaces := os.Getenv("TARGET_NAMESPACES")
	formattedNamespaces := []string{}

	if targetNamespaces != "" {
		formattedNamespaces = strings.Split(targetNamespaces, ",")
	}

	// Setup event listeners
	createListeners(formattedNamespaces)

	// Wait forever
	select {}
}
