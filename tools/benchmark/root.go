package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ydb-platform/etcd-ydb/pkg/etcd"
)

var RootCmd = &cobra.Command{
	Use: "benchmark",
}

var (
	endpoints    []string
	totalConns   uint
	totalClients uint
)

func init() {
	RootCmd.PersistentFlags().StringSliceVar(&endpoints, "endpoints", []string{"127.0.0.1:2379"}, "gRPC endpoints")
	RootCmd.PersistentFlags().UintVar(&totalConns, "conns", 1, "Total number of gRPC connections")
	RootCmd.PersistentFlags().UintVar(&totalClients, "clients", 1, "Total number of gRPC clients")
}

func newClients() ([]*etcd.Client, error) {
	if totalConns < uint(len(endpoints)) {
		return nil, fmt.Errorf("conns < len(endpoints)")
	}
	conns := make([]*etcd.Client, totalConns)
	for i := range conns {
		conn, err := etcd.NewClient(endpoints[i%len(endpoints)])
		if err != nil {
			return nil, err
		}
		conns[i] = conn
	}
	clients := make([]*etcd.Client, totalClients)
	for i := range clients {
		clients[i] = conns[i%len(conns)]
	}
	return clients, nil
}
