package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"vpa-manager/pkg/config"
	"vpa-manager/pkg/controller"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)

	env := config.LoadEnv()

	config, err := rest.InClusterConfig()
	if err != nil {
		slog.Error("Error creating in-cluster config", "error", err)
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		slog.Error("Error creating Kubernetes client", "error", err)
		panic(err)
	}

	controller := controller.NewController(env, config, clientset)

	if err := controller.Run(ctx); err != nil {
		slog.Error("Error running controller", "error", err)
		panic(err)
	}

	<-stopCh
	cancel()

	slog.Info("Shutting down controller")
}
