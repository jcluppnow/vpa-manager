package controller

import (
	"log/slog"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const controllerAgentName = "vpa-manager"

// Controller is the controller implementation for Foo resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
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

func (c *Controller) Run() error {
	// Start the informer factories to begin populating the informer caches
	slog.Info("Starting ", controllerAgentName, "controller")
	c.factory.Start(wait.NeverStop)
	c.factory.WaitForCacheSync(wait.NeverStop)

	return nil
}
