# Vertical Pod Autoscaler Manager
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/jcluppnow/vpa-manager/release.yaml)

VPA Manager is designed to automate the creation and deletion of VerticalPodAutoscaler (VPA) resources by monitoring the state of your workloads.
VPA Manager can be configured to monitor any number of namespaces and resources.

## Usage
[Helm](https://helm.sh/) must be installed to use the charts. Please refer to Helm's documentation to get started.

Once Helm is set up properly, add the repo as follows:

```console
helm repo add jcluppnow https://jcluppnow.github.io/vpa-manager
helm repo update
```

### Installing the Chart

To install the chart with the release name `my-release`:

```console
helm install my-release jcluppnow/vpa-manager
```

Chart documentation is available in `chart` directory.

## Contributing

We'd love to have you contribute! Please refer to our [contribution guidelines](https://github.com/jcluppnow/vpa-manager/blob/main/CONTRIBUTING.md) for details.

## License
VPA Manager is available under the Apache-2.0 license. See the [LICENSE](https://github.com/jcluppnow/vpa-manager/blob/main/LICENSE) file for more info.
