package controller

import (
	"context"
	"log/slog"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Controller struct {
	kubeclientset kubernetes.Interface
	factory       informers.SharedInformerFactory
}

func NewController(
	env ControllerEnv,
	config *rest.Config,
	clientset *kubernetes.Clientset,
) *Controller {
	slog.Info("Setting up event handlers")

	controller := &Controller{
		kubeclientset: clientset,
		factory:       CreateInformers(env, config, clientset),
	}

	return controller
}

func (c *Controller) Run(ctx context.Context) error {
	slog.Info("Starting controller")
	c.factory.Start(ctx.Done())
	c.factory.WaitForCacheSync(ctx.Done())

	return nil
}
