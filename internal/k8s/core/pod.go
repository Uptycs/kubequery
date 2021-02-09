/**
 * Copyright (c) 2020-present, The kubequery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package core

import (
	"context"

	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/kubequery/internal/k8s"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type pod struct {
	k8s.CommonNamespacedFields
	k8s.CommonPodFields
	v1.PodStatus
}

// PodColumns returns kubernetes pod fields as Osquery table columns.
func PodColumns() []table.ColumnDefinition {
	return k8s.GetSchema(&pod{})
}

// PodsGenerate generates the kubernetes pods as Osquery table data.
func PodsGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	options := metav1.ListOptions{}
	results := make([]map[string]string, 0)

	for {
		pods, err := k8s.GetClient().CoreV1().Pods(metav1.NamespaceAll).List(ctx, options)
		if err != nil {
			return nil, err
		}

		for _, p := range pods.Items {
			item := &pod{
				CommonNamespacedFields: k8s.GetCommonNamespacedFields(p.ObjectMeta),
				CommonPodFields:        k8s.GetCommonPodFields(p.Spec),
				PodStatus:              p.Status,
			}
			results = append(results, k8s.ToMap(item))
		}

		if pods.Continue == "" {
			break
		}
		options.Continue = pods.Continue
	}

	return results, nil
}

type podContainer struct {
	k8s.CommonNamespacedFields
	k8s.CommonContainerFields
	PodName              string
	ContainerType        string
	State                v1.ContainerState
	LastTerminationState v1.ContainerState
	Ready                bool
	RestartCount         int32
	ImageID              string
	ContainerID          string
	Started              *bool
}

// PodContainerColumns returns kubernetes pod container fields as Osquery table columns.
func PodContainerColumns() []table.ColumnDefinition {
	return k8s.GetSchema(&podContainer{})
}

func createPodContainer(p v1.Pod, c v1.Container, cs v1.ContainerStatus, containerType string) *podContainer {
	item := &podContainer{
		CommonNamespacedFields: k8s.GetCommonNamespacedFields(p.ObjectMeta),
		CommonContainerFields:  k8s.GetCommonContainerFields(c),
		PodName:                p.Name,
		ContainerType:          containerType,
		State:                  cs.State,
		LastTerminationState:   cs.LastTerminationState,
		Ready:                  cs.Ready,
		RestartCount:           cs.RestartCount,
		ImageID:                cs.ImageID,
		ContainerID:            cs.ContainerID,
		Started:                cs.Started,
	}
	item.Name = c.Name
	return item
}

func createPodEphemeralContainer(p v1.Pod, c v1.EphemeralContainer, cs v1.ContainerStatus) *podContainer {
	item := &podContainer{
		CommonNamespacedFields: k8s.GetCommonNamespacedFields(p.ObjectMeta),
		CommonContainerFields:  k8s.GetCommonEphemeralContainerFields(c),
		PodName:                p.Name,
		ContainerType:          "ephemeral",
		State:                  cs.State,
		LastTerminationState:   cs.LastTerminationState,
		Ready:                  cs.Ready,
		RestartCount:           cs.RestartCount,
		ImageID:                cs.ImageID,
		ContainerID:            cs.ContainerID,
		Started:                cs.Started,
	}
	item.Name = c.Name
	return item
}

// PodContainersGenerate generates the kubernetes pod containers as Osquery table data.
func PodContainersGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	options := metav1.ListOptions{}
	results := make([]map[string]string, 0)

	for {
		pods, err := k8s.GetClient().CoreV1().Pods(metav1.NamespaceAll).List(ctx, options)
		if err != nil {
			return nil, err
		}

		for _, p := range pods.Items {
			for i, c := range p.Spec.InitContainers {
				item := createPodContainer(p, c, p.Status.InitContainerStatuses[i], "init")
				results = append(results, k8s.ToMap(item))
			}
			for i, c := range p.Spec.Containers {
				item := createPodContainer(p, c, p.Status.ContainerStatuses[i], "container")
				results = append(results, k8s.ToMap(item))
			}
			for i, c := range p.Spec.EphemeralContainers {
				item := createPodEphemeralContainer(p, c, p.Status.EphemeralContainerStatuses[i])
				results = append(results, k8s.ToMap(item))
			}
		}

		if pods.Continue == "" {
			break
		}
		options.Continue = pods.Continue
	}

	return results, nil
}

type podVolume struct {
	k8s.CommonNamespacedFields
	k8s.CommonVolumeFields
	PodName string
}

// PodVolumeColumns returns kubernetes pod volume fields as Osquery table columns.
func PodVolumeColumns() []table.ColumnDefinition {
	return k8s.GetSchema(&podVolume{})
}

// PodVolumesGenerate generates the kubernetes pod volumes as Osquery table data.
func PodVolumesGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	options := metav1.ListOptions{}
	results := make([]map[string]string, 0)

	for {
		pods, err := k8s.GetClient().CoreV1().Pods(metav1.NamespaceAll).List(ctx, options)
		if err != nil {
			return nil, err
		}

		for _, p := range pods.Items {
			for _, v := range p.Spec.Volumes {
				item := &podVolume{
					CommonNamespacedFields: k8s.GetCommonNamespacedFields(p.ObjectMeta),
					CommonVolumeFields:     k8s.GetCommonVolumeFields(v),
					PodName:                p.Name,
				}
				item.Name = v.Name
				results = append(results, k8s.ToMap(item))
			}
		}

		if pods.Continue == "" {
			break
		}
		options.Continue = pods.Continue
	}

	return results, nil
}
