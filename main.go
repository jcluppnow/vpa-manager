package main

import (
	"os"
	"strings"
	"vpa-manager/controller"
	"vpa-manager/controller/utils"
)

func main() {
	targetNamespaces := os.Getenv("TARGET_NAMESPACES")
	formattedNamespaces := []string{}

	if targetNamespaces != "" {
		formattedNamespaces = strings.Split(targetNamespaces, ",")
	}

	resourcesToManage := controller.ResourcesToManage{
		Cronjobs:    utils.ParseBoolFromEnv("ENABLE_CRONJOBS"),
		Deployments: utils.ParseBoolFromEnv("ENABLE_DEPLOYMENTS"),
		Jobs:        utils.ParseBoolFromEnv("ENABLE_JOBS"),
		Pods:        utils.ParseBoolFromEnv("ENABLE_PODS"),
	}

	controller.CreateInformers(formattedNamespaces, resourcesToManage)
	select {}
}
