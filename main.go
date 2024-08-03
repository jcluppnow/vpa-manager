package main

import (
	"os"
	"strings"
	"vpa-manager/controller"
)

func main() {
	targetNamespaces := os.Getenv("TARGET_NAMESPACES")
	formattedNamespaces := []string{}

	if targetNamespaces != "" {
		formattedNamespaces = strings.Split(targetNamespaces, ",")
	}

	// Setup event listeners
	controller.CreateInformers(formattedNamespaces)

	// Wait forever
	select {}
}
