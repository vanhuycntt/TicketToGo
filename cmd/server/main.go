package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	rootCmd := &cobra.Command{
		Use:  "server",
		Long: "Mobile APIs",
	}
	rootCmd.AddCommand(&cobra.Command{
		Use:  "start",
		Long: "Start Mobile APIs",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("starting")
		},
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to run the service: %v", err)
	}
}
