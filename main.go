package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"vpa-manager/controller"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Handle OS signals for graceful shutdown
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)

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

	go controller.Run(ctx)

	// Block until a signal is received
	<-stopCh
	cancel() // Trigger context cancellation

	controller.ShutDown(ctx)
}
