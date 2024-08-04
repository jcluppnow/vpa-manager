package controller

import (
	"log/slog"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func isTargetNamespace(targetNamespaces []string, namespace string) bool {
	if len(targetNamespaces) == 0 {
		return true
	}

	for index := range targetNamespaces {
		if targetNamespaces[index] == namespace {
			return true
		}
	}

	return false
}

func CreateInformers(env ControllerEnv) {
	if !env.EnableCronjobs && !env.EnableDeployments && !env.EnableJobs && !env.EnablePods {
		slog.Warn("All resources types are disabled, as a result no Vertical Pod Autoscalers will be created. If this is not expected, review your configuration")
	}

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

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		slog.Error("Error creating dynamic client")
		panic(err)
	}

	factory := informers.NewSharedInformerFactory(clientset, time.Minute)

	jobInformer := factory.Batch().V1().Jobs().Informer()
	podInformer := factory.Core().V1().Pods().Informer()

	if env.EnableCronjobs {
		cronJobInformer := factory.Batch().V1().CronJobs().Informer()
		cronJobInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				cronJob := obj.(*batchv1.CronJob)
				if isTargetNamespace(env.TargetNamespaces, cronJob.Namespace) {
					CreateVPA(*client, "CronJob", cronJob.Name, cronJob.Namespace)
				}
			},
			DeleteFunc: func(obj interface{}) {
				cronJob := obj.(*batchv1.CronJob)
				if isTargetNamespace(env.TargetNamespaces, cronJob.Namespace) {
					DeleteVPA(*client, cronJob.Name, cronJob.Namespace)
				}
			},
		})
	}

	if env.EnableDeployments {
		deploymentInformer := factory.Apps().V1().Deployments().Informer()
		deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				deployment := obj.(*appsv1.Deployment)
				if isTargetNamespace(env.TargetNamespaces, deployment.Namespace) {
					CreateVPA(*client, "Deployment", deployment.Name, deployment.Namespace)
				}
			},
			DeleteFunc: func(obj interface{}) {
				deployment := obj.(*batchv1.CronJob)
				if isTargetNamespace(env.TargetNamespaces, deployment.Namespace) {
					DeleteVPA(*client, deployment.Name, deployment.Namespace)
				}
			},
		})
	}

	if env.EnableJobs {
		jobInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				job := obj.(*batchv1.Job)
				if len(job.OwnerReferences) == 0 && isTargetNamespace(env.TargetNamespaces, job.Namespace) {
					CreateVPA(*client, "Job", job.Name, job.Namespace)
				}
			},
			DeleteFunc: func(obj interface{}) {
				job := obj.(*batchv1.Job)
				if len(job.OwnerReferences) == 0 && isTargetNamespace(env.TargetNamespaces, job.Namespace) {
					DeleteVPA(*client, job.Name, job.Namespace)
				}
			},
		})
	}

	if env.EnablePods {
		podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				if len(pod.OwnerReferences) == 0 && isTargetNamespace(env.TargetNamespaces, pod.Namespace) {
					CreateVPA(*client, "Pod", pod.Name, pod.Namespace)
				}
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				if len(pod.OwnerReferences) == 0 && isTargetNamespace(env.TargetNamespaces, pod.Namespace) {
					DeleteVPA(*client, pod.Name, pod.Namespace)
				}
			},
		})
	}

	factory.Start(wait.NeverStop)
	factory.WaitForCacheSync(wait.NeverStop)
}
