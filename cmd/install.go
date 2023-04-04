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
		k8s := cmd.Flags().Lookup("k8s").Value.String()

		region := cmd.Flags().Lookup("region").Value.String()
		if region == "" {
			region, _ = utils.GetDefaultAwsRegion()
		}

		gcpProjectId := cmd.Flags().Lookup("gcp_project_id").Value.String()
		if k8s == "gke" && gcpProjectId == "" {
			fmt.Println("Error: You need to provide the --gcp_project_id if using GKE")
			return
		}

		// Let's do the setup
		releaseIdentifer := utils.GenerateReleaseIdentifier()

		inputData := utils.SetupData{
			Identifier:   releaseIdentifer,
			ReleaseName:  name,
			Namespace:    namespace,
			Region:       region,
			K8s:          k8s,
			GCPProjectId: gcpProjectId,
		}

		if namespace == "" { // if namespace is not provided, use the current namespace
			namespace, _ := utils.GetCurrentNamespace()
			inputData.Namespace = namespace

			setupData, err := utils.Setup(inputData)
			if err != nil {
				fmt.Println("Error: ", err)
				panic(err)
			}

			setupData.Namespace = namespace

			utils.CreateConfigMap(setupData)
		} else { // if namespace is provided, use that namespace
			setupData, err := utils.Setup(inputData)
			if err != nil {
				fmt.Println("Error: ", err)
				panic(err)
			}

			utils.CreateConfigMap(setupData)
		}

	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	namespace := ""
	installCmd.Flags().StringVar(&namespace, "namespace", "", "namespace to install the helm chart1")

	region := ""
	installCmd.Flags().StringVar(&region, "region", "", "region to install the installation in.")

	gcpProjectId := ""
	installCmd.Flags().StringVar(&gcpProjectId, "gcp_project_id", "", "GCP Project ID to install the installation in.")

}
