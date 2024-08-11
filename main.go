package main

import (
	"log/slog"
	"vpa-manager/controller"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	env := controller.LoadEnv()

	config, err := rest.InClusterConfig()
	if err != nil {
		slog.Error("Error creating in-cluster config")
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		slog.Error("Error creating Kubernetes client")
		panic(err)
	}

	controller := controller.NewController(env, config, clientset)

	controller.Run()

	select {}
}
