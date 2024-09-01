# VPA Manager Helm Chart

## Get Repo Info

```console
helm repo add jcluppnow https://jcluppnow.github.io/vpa-manager
helm repo update
```

_See [helm repo](https://helm.sh/docs/helm/helm_repo/) for command documentation._

## Installing the Chart

To install the chart with the release name `my-release`:

```console
helm install my-release jcluppnow/vpa-manager
```

## Uninstalling the Chart

To uninstall the `my-release` deployment:

```console
helm uninstall my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

| Parameter                                 | Description                                   | Default                                                 |
|-------------------------------------------|-----------------------------------------------|---------------------------------------------------------|
| `affinity`                                | Affinity settings for pod assignment          | `{}`                                                    |
| `deploymentAnnotations`                     | Deployment annotations (can be templated)        | `{}`                                                    |
| `image.pullPolicy`          | VPA Manager image pull policy   | `IfNotPresent`                                          |
| `image.repository`          | VPA Manager container image repository    | `vpa-manager`                                               |
| `image.tag`                 | VPA Manager container image tag           | `1.31.1`                                                |
| `imageRenderer.nodeSelector`               | Node labels for pod assignment                | `{}`                                                    |
| `podAnnotations`                     | Pod annotations (can be templated)        | `{}`                                                    |
| `priorityClassName`                       | Name of Priority Class to assign pods         | `nil`                                                   |
| `resources`                               | CPU/Memory resource requests/limits           | `{}`                                                    |
| `resourcesToManage.cronjobs`                               | Toggle for whether this resource type is in scope for VPA creation     | `True`                                                    |
| `resourcesToManage.deployments`                               | Toggle for whether this resource type is in scope for VPA creation           | `True`                                                    |
| `resourcesToManage.jobs`                               | Toggle for whether this resource type is in scope for VPA creation           | `True`                                                    |
| `resourcesToManager.pods`                               | Toggle for whether this resource type is in scope for VPA creation           | `True`                                                    |
| `revisionHistoryLimit`                    | Number of old ReplicaSets to retain           | `10`                                                    |
| `serviceAccountAnnotations`                     | Service Account annotations (can be templated)        | `{}`                                                    |
| `tolerations`                | Toleration labels for pod assignment          | `[]`                                                    |
| `updateMode`                | Parameter to control which mode each Vertical Pod Autoscaler resource will have          | `Off`                                                    |
| `watchedNamespaces`                | List of namespaces to monitor for VPA creation, empty defaults to ALL namespaces       | `[]`       
