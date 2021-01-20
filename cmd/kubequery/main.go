/**
 * Copyright (c) 2020-present, The kubequery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/Uptycs/kubequery/internal/k8s"
	"github.com/Uptycs/kubequery/internal/k8s/admissionregistration"
	"github.com/Uptycs/kubequery/internal/k8s/apps"
	"github.com/Uptycs/kubequery/internal/k8s/autoscaling"
	"github.com/Uptycs/kubequery/internal/k8s/batch"
	"github.com/Uptycs/kubequery/internal/k8s/core"
	"github.com/Uptycs/kubequery/internal/k8s/discovery"
	"github.com/Uptycs/kubequery/internal/k8s/networking"
	"github.com/Uptycs/kubequery/internal/k8s/policy"
	"github.com/Uptycs/kubequery/internal/k8s/rbac"
	"github.com/Uptycs/kubequery/internal/k8s/storage"

	"github.com/kolide/osquery-go"
	"github.com/kolide/osquery-go/plugin/table"
)

var (
	socket   = flag.String("socket", "", "Path to the extensions UNIX domain socket")
	timeout  = flag.Int("timeout", 3, "Seconds to wait for autoloaded extensions")
	interval = flag.Int("interval", 3, "Seconds delay between connectivity checks")
)

func registerTables(server *osquery.ExtensionManagerServer) {
	server.RegisterPlugin(
		// Admission Registration
		table.NewPlugin("kubernetes_mutating_webhooks", admissionregistration.MutatingWebhookColumns(), admissionregistration.MutatingWebhooksGenerate),
		table.NewPlugin("kubernetes_validating_webhooks", admissionregistration.ValidatingWebhookColumns(), admissionregistration.ValidatingWebhooksGenerate),

		// Apps
		table.NewPlugin("kubernetes_daemon_sets", apps.DaemonSetColumns(), apps.DaemonSetsGenerate),
		table.NewPlugin("kubernetes_daemon_set_containers", apps.DaemonSetContainerColumns(), apps.DaemonSetContainersGenerate),
		table.NewPlugin("kubernetes_daemon_set_volumes", apps.DaemonSetVolumeColumns(), apps.DaemonSetVolumesGenerate),
		table.NewPlugin("kubernetes_deployments", apps.DeploymentColumns(), apps.DeploymentsGenerate),
		table.NewPlugin("kubernetes_deployments_containers", apps.DeploymentContainerColumns(), apps.DeploymentContainersGenerate),
		table.NewPlugin("kubernetes_deployments_volumes", apps.DeploymentVolumeColumns(), apps.DeploymentVolumesGenerate),
		table.NewPlugin("kubernetes_replica_sets", apps.ReplicaSetColumns(), apps.ReplicaSetsGenerate),
		table.NewPlugin("kubernetes_replica_set_containers", apps.ReplicaSetContainerColumns(), apps.ReplicaSetContainersGenerate),
		table.NewPlugin("kubernetes_replica_set_volumes", apps.ReplicaSetVolumeColumns(), apps.ReplicaSetVolumesGenerate),
		table.NewPlugin("kubernetes_stateful_sets", apps.StatefulSetColumns(), apps.StatefulSetsGenerate),
		table.NewPlugin("kubernetes_stateful_set_containers", apps.StatefulSetContainerColumns(), apps.StatefulSetContainersGenerate),
		table.NewPlugin("kubernetes_stateful_set_volumes", apps.StatefulSetVolumeColumns(), apps.StatefulSetVolumesGenerate),

		// Autoscaling
		table.NewPlugin("kubernetes_horizontal_pod_autoscalers", autoscaling.HorizontalPodAutoscalersColumns(), autoscaling.HorizontalPodAutoscalerGenerate),

		// Batch
		table.NewPlugin("kubernetes_cron_jobs", batch.CronJobColumns(), batch.CronJobsGenerate),
		table.NewPlugin("kubernetes_jobs", batch.JobColumns(), batch.JobsGenerate),

		// Core
		table.NewPlugin("kubernetes_component_statuses", core.ComponentStatusColumns(), core.ComponentStatusesGenerate),
		table.NewPlugin("kubernetes_config_maps", core.ConfigMapColumns(), core.ConfigMapsGenerate),
		table.NewPlugin("kubernetes_endpoint_subsets", core.EndpointSubsetColumns(), core.EndpointSubsetsGenerate),
		table.NewPlugin("kubernetes_limit_ranges", core.LimitRangeColumns(), core.LimitRangesGenerate),
		table.NewPlugin("kubernetes_namespaces", core.NamespaceColumns(), core.NamespacesGenerate),
		table.NewPlugin("kubernetes_nodes", core.NodeColumns(), core.NodesGenerate),
		table.NewPlugin("kubernetes_persistent_volume_claims", core.PersistentVolumeClaimColumns(), core.PersistentVolumeClaimsGenerate),
		table.NewPlugin("kubernetes_persistent_volumes", core.PersistentVolumeColumns(), core.PersistentVolumesGenerate),
		table.NewPlugin("kubernetes_pod_templates", core.PodTemplateColumns(), core.PodTemplatesGenerate),
		table.NewPlugin("kubernetes_pods", core.PodColumns(), core.PodsGenerate),
		table.NewPlugin("kubernetes_pod_containers", core.PodContainerColumns(), core.PodContainersGenerate),
		table.NewPlugin("kubernetes_pod_volumes", core.PodVolumeColumns(), core.PodVolumesGenerate),
		table.NewPlugin("kubernetes_resource_quotas", core.ResourceQuotaColumns(), core.ResourceQuotasGenerate),
		table.NewPlugin("kubernetes_secrets", core.SecretColumns(), core.SecretsGenerate),
		table.NewPlugin("kubernetes_service_accounts", core.ServiceAccountColumns(), core.ServiceAccountsGenerate),
		table.NewPlugin("kubernetes_services", core.ServiceColumns(), core.ServicesGenerate),

		// Discovery
		table.NewPlugin("kubernetes_api_resources", discovery.APIResourceColumns(), discovery.APIResourcesGenerate),
		table.NewPlugin("kubernetes_info", discovery.InfoColumns(), discovery.InfoGenerate),

		// Networking
		table.NewPlugin("kubernetes_ingress_classes", networking.IngressClassColumns(), networking.IngressClassesGenerate),
		table.NewPlugin("kubernetes_ingresses", networking.IngressColumns(), networking.IngressesGenerate),
		table.NewPlugin("kubernetes_network_policies", networking.NetworkPolicyColumns(), networking.NetworkPoliciesGenerate),

		// Policy
		table.NewPlugin("kubernetes_pod_disruption_budget", policy.PodDisruptionBudgetColumns(), policy.PodDisruptionBudgetsGenerate),
		table.NewPlugin("kubernetes_pod_security_policies", policy.PodSecurityPolicyColumns(), policy.PodSecurityPoliciesGenerate),

		// RBAC
		table.NewPlugin("kubernetes_cluster_role_binding_subjects", rbac.ClusterRoleBindingSubjectColumns(), rbac.ClusterRoleBindingSubjectsGenerate),
		table.NewPlugin("kubernetes_cluster_role_policy_rule", rbac.ClusterRolePolicyRuleColumns(), rbac.ClusterRolePolicyRulesGenerate),
		table.NewPlugin("kubernetes_role_binding_subjects", rbac.RoleBindingSubjectColumns(), rbac.RoleBindingSubjectsGenerate),
		table.NewPlugin("kubernetes_role_policy_rule", rbac.RolePolicyRuleColumns(), rbac.RolePolicyRulesGenerate),

		// Storage
		table.NewPlugin("kubernetes_csi_drivers", storage.CSIDriverColumns(), storage.CSIDriversGenerate),
		table.NewPlugin("kubernetes_csi_node_drivers", storage.CSINodeDriverColumns(), storage.CSINodeDriversGenerate),
		table.NewPlugin("kubernetes_storage_capacities", storage.CSIStorageCapacityColumns(), storage.CSIStorageCapacitiesGenerate),
		table.NewPlugin("kubernetes_storage_classes", storage.SGClassColumns(), storage.SGClassesGenerate),
		table.NewPlugin("kubernetes_volume_attachments", storage.VolumeAttachmentColumns(), storage.VolumeAttachmentsGenerate),
	)
}

func main() {
	flag.Parse()
	if *socket == "" {
		panic("Missing required --socket argument")
	}

	err := k8s.Init()
	if err != nil {
		panic(err.Error())
	}

	// TODO: Version and SDK version
	server, err := osquery.NewExtensionManagerServer(
		"kubequery",
		*socket,
		osquery.ServerTimeout(time.Second*time.Duration(*timeout)),
		osquery.ServerPingInterval(time.Second*time.Duration(*interval)),
	)
	if err != nil {
		panic(fmt.Sprintf("Error launching kubequery: %s\n", err))
	}

	registerTables(server)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
