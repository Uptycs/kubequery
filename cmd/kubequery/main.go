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
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	osquery "github.com/Uptycs/basequery-go"
	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/kubequery/internal/k8s"
	"github.com/Uptycs/kubequery/internal/k8s/event"
	"github.com/Uptycs/kubequery/internal/k8s/tables"
)

var (
	verbose  = flag.Bool("verbose", false, "Whether to enable verbose logging")
	socket   = flag.String("socket", "", "Path to the extensions UNIX domain socket")
	timeout  = flag.Int("timeout", 5, "Seconds to wait for autoloaded extensions")
	interval = flag.Int("interval", 5, "Seconds delay between connectivity checks")
)

func main() {
	flag.Parse()
	if *socket == "" {
		panic("Missing required --socket argument")
	}

	err := k8s.Init()
	if err != nil {
		panic(fmt.Sprintf("Error connecting to kubernetes API server: %s", err))
	}

	// TODO: Version
	server, err := osquery.NewExtensionManagerServer(
		"kubequery",
		*socket,
		osquery.ServerVersion("1.0.0"),
		osquery.ServerTimeout(time.Second*time.Duration(*timeout)),
		osquery.ServerPingInterval(time.Second*time.Duration(*interval)),
	)
	if err != nil {
		panic(fmt.Sprintf("Error launching kubequery: %s", err))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	for _, t := range tables.GetTables() {
		server.RegisterPlugin(table.NewPlugin(t.Name, t.Columns, t.GenFunc))
	}

	go func() {
		if err := server.Run(); err != nil {
			panic(fmt.Sprintf("Failed to start extension manager server: %s", err))
		}
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	// Wait for the extension manager to start before sending events
	time.Sleep(time.Second * 3)

	watcher, err := event.CreateEventWatcher(*socket, time.Second*time.Duration(*timeout))
	if err != nil {
		fmt.Println("Failed to create kubernetes event watcher: ", err)
	} else {
		watcher.Start()
	}

	<-quit

	if watcher != nil {
		watcher.Stop()
	}
	server.Shutdown(context.Background())
}
