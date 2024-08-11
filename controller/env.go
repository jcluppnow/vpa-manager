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
	WatchedNamespaces []string
}

func LoadEnv() ControllerEnv {
	watchedNamespaces := os.Getenv("WATCHED_NAMESPACES")
	formattedNamespaces := []string{}

	if watchedNamespaces != "" {
		formattedNamespaces = strings.Split(watchedNamespaces, ",")
	}

	env := ControllerEnv{
		EnableCronjobs:    utils.ParseBoolFromEnv("ENABLE_CRONJOBS"),
		EnableDeployments: utils.ParseBoolFromEnv("ENABLE_DEPLOYMENTS"),
		EnableJobs:        utils.ParseBoolFromEnv("ENABLE_JOBS"),
		EnablePods:        utils.ParseBoolFromEnv("ENABLE_PODS"),
		WatchedNamespaces: formattedNamespaces,
	}

	return env
}
