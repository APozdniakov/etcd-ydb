package main

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "benchmark",
}

var (
	endpoint     string
	totalConns   uint
	totalClients uint
)

func init() {
	RootCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "127.0.0.1:2379", "gRPC endpoints")
	RootCmd.PersistentFlags().UintVar(&totalConns, "conns", 1, "Total number of gRPC connections")
	RootCmd.PersistentFlags().UintVar(&totalClients, "clients", 1, "Total number of gRPC clients")
}
