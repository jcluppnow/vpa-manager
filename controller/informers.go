package controller

import (
	"log/slog"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func CreateInformers(env ControllerEnv, config *rest.Config, clientset *kubernetes.Clientset) informers.SharedInformerFactory {
	if !env.EnableCronjobs && !env.EnableDeployments && !env.EnableJobs && !env.EnablePods {
		slog.Warn("All resources types are disabled, as a result no Vertical Pod Autoscalers will be created. If this is not expected, review your configuration")
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		slog.Error("Error creating dynamic client", "error", err)
		panic(err)
	}

	factory := informers.NewSharedInformerFactory(clientset, time.Minute)

	if env.EnableCronjobs {
		cronJobInformer := factory.Batch().V1().CronJobs().Informer()
		cronJobInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				cronJob := obj.(*batchv1.CronJob)
				CreateVPA(*client, env.TargetNamespaces, "CronJob", cronJob.Name, cronJob.Namespace)
			},
			DeleteFunc: func(obj interface{}) {
				cronJob := obj.(*batchv1.CronJob)
				DeleteVPA(*client, env.TargetNamespaces, cronJob.Name, cronJob.Namespace)
			},
		})
	}

	if env.EnableDeployments {
		deploymentInformer := factory.Apps().V1().Deployments().Informer()
		deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				deployment := obj.(*appsv1.Deployment)
				CreateVPA(*client, env.TargetNamespaces, "Deployment", deployment.Name, deployment.Namespace)
			},
			DeleteFunc: func(obj interface{}) {
				deployment := obj.(*appsv1.Deployment)
				DeleteVPA(*client, env.TargetNamespaces, deployment.Name, deployment.Namespace)
			},
		})
	}

	if env.EnableJobs {
		jobInformer := factory.Batch().V1().Jobs().Informer()
		jobInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				job := obj.(*batchv1.Job)
				if len(job.OwnerReferences) == 0 {
					CreateVPA(*client, env.TargetNamespaces, "Job", job.Name, job.Namespace)
				}
			},
			DeleteFunc: func(obj interface{}) {
				job := obj.(*batchv1.Job)
				if len(job.OwnerReferences) == 0 {
					DeleteVPA(*client, env.TargetNamespaces, job.Name, job.Namespace)
				}
			},
		})
	}

	if env.EnablePods {
		podInformer := factory.Core().V1().Pods().Informer()
		podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				if len(pod.OwnerReferences) == 0 {
					CreateVPA(*client, env.TargetNamespaces, "Pod", pod.Name, pod.Namespace)
				}
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				if len(pod.OwnerReferences) == 0 {
					DeleteVPA(*client, env.TargetNamespaces, pod.Name, pod.Namespace)
				}
			},
		})
	}

	return factory
}
