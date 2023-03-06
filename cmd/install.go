/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs ZincObserve",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
observe provides easy yet sophisticated observability for your Kubernetes clusters
Zinc.`,
	Run: func(cmd *cobra.Command, args []string) {
		name := cmd.Flags().Lookup("name").Value.String()
		fmt.Println("install called with name: ", name)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	K8s := "eks"
	installCmd.Flags().StringVar(&K8s, "k8s", "eks", "k8s cluster type. eks, gke, plain")
	installCmd.MarkFlagRequired("k8s")

}
