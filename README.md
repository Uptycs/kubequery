[![Go Report Card](https://goreportcard.com/badge/github.com/Uptycs/kubequery)](https://goreportcard.com/report/github.com/Uptycs/kubequery) [![FOSSA Status](https://app.fossa.com/api/projects/custom%2B22616%2Fgit%40github.com%3AUptycs%2Fkubequery.git.svg?type=shield)](https://app.fossa.com/projects/custom%2B22616%2Fgit%40github.com%3AUptycs%2Fkubequery.git?ref=badge_shield) ![CodeQL](https://github.com/Uptycs/kubequery/workflows/CodeQL/badge.svg?branch=master)

# kubequery powered by Osquery

kubequery is a [Osquery](https://osquery.io) extension that provides SQL based analytics for [Kubernetes](https://kubernetes.io) clusters

kubequery will be packaged as docker image available from [dockerhub](https://hub.docker.com/r/uptycs/kubequery). It is expected to be deployed as a [Kubernetes Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment) per cluster. A sample deployment template is available [here](kubequery.yaml)


kubequery tables [schema is available here](docs/schema.md)

## Build

`Go 1.15` and `make` are required to build kubequery. Run: `make`

## FAQ

* Kubernetes events support?

`kubenetes_events` table can be easily implemented in kubequery as traditional table. But ideally it should be a streaming events table similar to `process_events` etc in Osquery. Unfortunately Osquery does not support event tables in extensions currently. Buffering the data in extension and periodically sending it in response to a query is one option, but it is not ideal.

* Why are some columns JSON?

Normalizing nested JSON data like Kubernetes API responses will create an explosion of tables. So some of the columns in kuberenetes tables are left as JSON. Data is eventually processed by [SQLite](https://www.sqlite.org/index.html) with in Osquery. SQLite has very [good JSON](https://www.sqlite.org/json1.html) support. To get the `value` of `rule` in `run_as_user` column from `kubernetes_pod_security_policies` table, the following query can be used:
```sql
  SELECT value FROM kubernetes_pod_security_policies, json_tree(kubernetes_pod_security_policies.run_as_user) WHERE key = 'rule';
```

When streaming data (example: Osquery TLS) from various kubernetes clusters, Lamba like functions can be applied on rows of data. Labmda can extract necessary fields from embedded JSON. If tables are normalized, it will not be trivial to JOIN across them and trigger events/alerts.