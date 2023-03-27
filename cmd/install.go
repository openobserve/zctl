/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zinclabs/zctl/pkg/utils"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs ZincObserve",
	Long: `Installs ZincObserve. The subtasks include:
	1. Create a S3 bucket
	2. Create an IAM role with inline policy that is trusted by EKS OIDC provider
	3. Install in the EKS cluster using the helm chart
	3. Create a ConfigMap with the release_name, bucket_name and IAM role
	`,
	Run: func(cmd *cobra.Command, args []string) {
		name := cmd.Flags().Lookup("name").Value.String()

		namespace := cmd.Flags().Lookup("namespace").Value.String()

		region := cmd.Flags().Lookup("region").Value.String()
		if region == "" {
			region, _ = utils.GetDefaultAwsRegion()
		}

		// Let's do the setup
		releaseIdentifer := utils.GenerateReleaseIdentifier()

		if namespace == "" {
			namespace, _ := utils.GetCurrentNamespace()

			setupData, err := utils.Setup(releaseIdentifer, name, namespace, region)
			if err != nil {
				fmt.Println("Error: ", err)
			}

			utils.CreateConfigMap(setupData, namespace)
		} else {
			setupData, err := utils.Setup(releaseIdentifer, name, namespace, region)
			if err != nil {
				fmt.Println("Error: ", err)
			}

			utils.CreateConfigMap(setupData, namespace)
		}

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

	namespace := ""
	installCmd.Flags().StringVar(&namespace, "namespace", "", "namespace to install the helm chart1")
	// installCmd.MarkFlagRequired("namespace")

}
