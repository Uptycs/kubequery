/**
 * Copyright (c) 2020-present, The kubequery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package k8s

import (
	"context"
	"sync"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	lock       sync.Mutex
	clientset  kubernetes.Interface
	clusterUID types.UID
)

func initClientset(config *rest.Config) error {
	if config == nil {
		// Get in-cluster configuration if one is not provided
		conf, err := rest.InClusterConfig()
		if err != nil {
			return err
		}
		config = conf
	}

	var err error
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	return nil
}

func initUID() error {
	ks, err := GetClient().CoreV1().Namespaces().Get(context.TODO(), "kube-system", v1.GetOptions{})
	if err != nil {
		return err
	}
	clusterUID = ks.UID
	return nil
}

// Init creates in-cluster kubernetes configuration and a client set using the configuration.
// This returns error if KUBERNETES_SERVICE_HOST or KUBERNETES_SERVICE_PORT environment variables are not set.
func Init() error {
	lock.Lock()
	defer lock.Unlock()

	err := initClientset(nil)
	if err != nil {
		return err
	}
	err = initUID()
	if err != nil {
		return err
	}

	return nil
}

// GetClient returns kubernetes interface that can be used to communicate with API server.
func GetClient() kubernetes.Interface {
	return clientset
}

// GetClusterUID returns unique identifier for the current kubernetes cluster.
// This is same as the kube-system namespace UID.
func GetClusterUID() types.UID {
	return clusterUID
}

// SetClient is helper function to override the kubernetes interface with fake one for testing.
func SetClient(c kubernetes.Interface, u types.UID) {
	lock.Lock()
	defer lock.Unlock()

	clientset = c
	clusterUID = u
}
