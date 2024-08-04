package controller

import (
	"os"
	"strings"
	"vpa-manager/controller/utils"
)

type ControllerEnv struct {
	Cronjobs         bool
	Deployments      bool
	Jobs             bool
	Pods             bool
	TargetNamespaces []string
}

func LoadEnv() ControllerEnv {
	targetNamespaces := os.Getenv("TARGET_NAMESPACES")
	formattedNamespaces := []string{}

	if targetNamespaces != "" {
		formattedNamespaces = strings.Split(targetNamespaces, ",")
	}

	env := ControllerEnv{
		Cronjobs:         utils.ParseBoolFromEnv("ENABLE_CRONJOBS"),
		Deployments:      utils.ParseBoolFromEnv("ENABLE_DEPLOYMENTS"),
		Jobs:             utils.ParseBoolFromEnv("ENABLE_JOBS"),
		Pods:             utils.ParseBoolFromEnv("ENABLE_PODS"),
		TargetNamespaces: formattedNamespaces,
	}

	return env
}
