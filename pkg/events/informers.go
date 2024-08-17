package events

import (
	"log/slog"
	"time"
	"vpa-manager/pkg/config"
	"vpa-manager/pkg/vpa"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func CreateInformers(env config.ControllerEnv, config *rest.Config, clientset kubernetes.Interface) informers.SharedInformerFactory {
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
				vpa.CreateVPA(client, env.WatchedNamespaces, "CronJob", cronJob.Name, cronJob.Namespace, env.UpdateMode)
			},
			DeleteFunc: func(obj interface{}) {
				cronJob := obj.(*batchv1.CronJob)
				vpa.DeleteVPA(client, env.WatchedNamespaces, cronJob.Name, cronJob.Namespace)
			},
		})
	}

	if env.EnableDeployments {
		deploymentInformer := factory.Apps().V1().Deployments().Informer()
		deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				deployment := obj.(*appsv1.Deployment)
				vpa.CreateVPA(client, env.WatchedNamespaces, "Deployment", deployment.Name, deployment.Namespace, env.UpdateMode)
			},
			DeleteFunc: func(obj interface{}) {
				deployment := obj.(*appsv1.Deployment)
				vpa.DeleteVPA(client, env.WatchedNamespaces, deployment.Name, deployment.Namespace)
			},
		})
	}

	if env.EnableJobs {
		jobInformer := factory.Batch().V1().Jobs().Informer()
		jobInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				job := obj.(*batchv1.Job)
				if len(job.OwnerReferences) == 0 {
					vpa.CreateVPA(client, env.WatchedNamespaces, "Job", job.Name, job.Namespace, env.UpdateMode)
				}
			},
			DeleteFunc: func(obj interface{}) {
				job := obj.(*batchv1.Job)
				if len(job.OwnerReferences) == 0 {
					vpa.DeleteVPA(client, env.WatchedNamespaces, job.Name, job.Namespace)
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
					vpa.CreateVPA(client, env.WatchedNamespaces, "Pod", pod.Name, pod.Namespace, env.UpdateMode)
				}
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				if len(pod.OwnerReferences) == 0 {
					vpa.DeleteVPA(client, env.WatchedNamespaces, pod.Name, pod.Namespace)
				}
			},
		})
	}

	return factory
}
