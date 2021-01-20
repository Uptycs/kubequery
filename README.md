[![Go Report Card](https://goreportcard.com/badge/github.com/Uptycs/kubequery)](https://goreportcard.com/report/github.com/Uptycs/kubequery)

# kubequery powered by Osquery

kubequery is a [Osquery](https://osquery.io) extension that provides SQL based analytics for [Kubernetes](https://kubernetes.io) clusters

kubequery will be packaged as docker image available from [dockerhub](https://hub.docker.com/r/uptycs/kubequery). It is expected to be deployed as a [Kubernetes Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment) per cluster. A sample deployment template is available [here](kubequery.yaml)


kubequery tables [schema is available here](docs/schema.md)

## Build

Go 1.15 and make are required to build kubequery. Run:

`make`

## FAQ

* Kubernetes events support
* Why are some columns JSON?
