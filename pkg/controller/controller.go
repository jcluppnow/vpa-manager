package controller

import (
	"context"
	"log/slog"
	"vpa-manager/pkg/config"
	"vpa-manager/pkg/events"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Controller struct {
	kubeclientset kubernetes.Interface
	factory       informers.SharedInformerFactory
}

func NewController(
	env config.ControllerEnv,
	restConfig *rest.Config,
	clientset kubernetes.Interface,
) *Controller {
	slog.Info("Validating environment config")

	config.ValidateControllerEnv(env)

	slog.Info("Setting up event handlers")

	controller := &Controller{
		kubeclientset: clientset,
		factory:       events.CreateInformers(env, restConfig, clientset),
	}

	return controller
}

func (c *Controller) Run(ctx context.Context) error {
	slog.Info("Starting controller")
	c.factory.Start(ctx.Done())
	c.factory.WaitForCacheSync(ctx.Done())

	return nil
}
