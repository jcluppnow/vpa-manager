package controller

import (
	"os"
	"strings"
	"vpa-manager/controller/utils"
)

type ControllerEnv struct {
	EnableCronjobs    bool
	EnableDeployments bool
	EnableJobs        bool
	EnablePods        bool
	TargetNamespaces  []string
}

func LoadEnv() ControllerEnv {
	targetNamespaces := os.Getenv("TARGET_NAMESPACES")
	formattedNamespaces := []string{}

	if targetNamespaces != "" {
		formattedNamespaces = strings.Split(targetNamespaces, ",")
	}

	env := ControllerEnv{
		EnableCronjobs:    utils.ParseBoolFromEnv("ENABLE_CRONJOBS"),
		EnableDeployments: utils.ParseBoolFromEnv("ENABLE_DEPLOYMENTS"),
		EnableJobs:        utils.ParseBoolFromEnv("ENABLE_JOBS"),
		EnablePods:        utils.ParseBoolFromEnv("ENABLE_PODS"),
		TargetNamespaces:  formattedNamespaces,
	}

	return env
}
