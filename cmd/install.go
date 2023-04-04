/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zinclabs/zctl/pkg/utils"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs ZincObserve",
	Long: `Installs ZincObserve.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		name := cmd.Flags().Lookup("name").Value.String()
		namespace := cmd.Flags().Lookup("namespace").Value.String()
		k8s := cmd.Flags().Lookup("k8s").Value.String()
		install_minio := cmd.Flags().Lookup("install_minio").Value.String()
		storage_provider := cmd.Flags().Lookup("storage_provider").Value.String()
		s3_access_key := cmd.Flags().Lookup("s3_access_key").Value.String()
		s3_secret_key := cmd.Flags().Lookup("s3_secret_key").Value.String()
		s3_server_url := cmd.Flags().Lookup("s3_server_url").Value.String()
		s3_bucket_name := cmd.Flags().Lookup("s3_bucket_name").Value.String()
		region := cmd.Flags().Lookup("region").Value.String()
		gcpProjectId := cmd.Flags().Lookup("gcp_project_id").Value.String()

		// convert install_minio to bool
		install_minio_bool, err := strconv.ParseBool(install_minio)
		if err != nil {
			fmt.Println("Error: ", err)
			panic(err)
		}

		// Let's do the setup
		releaseIdentifer := utils.GenerateReleaseIdentifier()

		inputData := utils.SetupData{
			Identifier:      releaseIdentifer,
			ReleaseName:     name,
			Namespace:       namespace,
			Region:          region,
			K8s:             k8s,
			GCPProjectId:    gcpProjectId,
			S3AccessKey:     s3_access_key,
			S3SecretKey:     s3_secret_key,
			InstallMinIO:    install_minio_bool,
			StorageProvider: storage_provider,
			S3ServerURL:     s3_server_url,
			BucketName:      s3_bucket_name,
		}

		inputData, err = ValidateAndFix(inputData)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		setupData, err := utils.Setup(inputData)
		if err != nil {
			fmt.Println("Error: ", err)
			panic(err)
		}

		utils.CreateConfigMap(setupData)
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

	install_minio := "false"
	installCmd.Flags().StringVar(&install_minio, "install_minio", "", "Specify if you want to install minio. Default is false.")

	storage_provider := ""
	installCmd.Flags().StringVar(&storage_provider, "storage_provider", "", "Valid values are s3, gcs, minio, swift.")

	s3_bucket_name := ""
	installCmd.Flags().StringVar(&s3_bucket_name, "s3_bucket_name", "", "s3 compatible bucket.")

	s3_server_url := ""
	installCmd.Flags().StringVar(&s3_server_url, "s3_server_url", "", "s3 compatible server url.")

	s3_access_key := ""
	installCmd.Flags().StringVar(&s3_access_key, "s3_access_key", "", "s3_access_key to use.")

	s3_secret_key := ""
	installCmd.Flags().StringVar(&s3_secret_key, "s3_secret_key", "", "s3_secret_key to use.")
}

// ValidateAndFix validates the input data and fixes it if possible
func ValidateAndFix(setupData utils.SetupData) (utils.SetupData, error) {
	if setupData.K8s == "eks" && setupData.Region == "" {
		setupData.Region, _ = utils.GetDefaultAwsRegion()
	}

	if setupData.K8s == "gke" && setupData.GCPProjectId == "" {
		return setupData, fmt.Errorf("error: You need to provide the --gcp_project_id if using GKE")
	}

	if setupData.Namespace == "" {
		namespace, _ := utils.GetCurrentNamespace()
		setupData.Namespace = namespace
	}

	if setupData.K8s == "plain" {
		if setupData.StorageProvider == "minio" && !setupData.InstallMinIO {
			if setupData.S3ServerURL == "" || setupData.S3AccessKey == "" || setupData.S3SecretKey == "" {
				return setupData, fmt.Errorf("error: You need to provide the --s3_server_url, --s3_access_key and --s3_secret_key if using minio and --install_minio=false")
			}
		}

		if setupData.StorageProvider == "minio" && setupData.InstallMinIO {
			if setupData.S3ServerURL != "" || setupData.S3AccessKey != "" || setupData.S3SecretKey != "" || setupData.BucketName != "" {
				return setupData, fmt.Errorf("error: You cannot provide --s3_access_key, --s3_secret_key and --s3_server_url when installing minio. It will automatically get picked from the minio installation. Please remove these flags and try again")
			}
		}

		if setupData.StorageProvider != "minio" {
			if setupData.S3ServerURL == "" || setupData.S3AccessKey == "" || setupData.S3SecretKey == "" {
				return setupData, fmt.Errorf("error: You need to provide --s3_access_key, --s3_secret_key and --s3_server_url when using a non minio storage provider. Please use these flags and try again")
			}
		}
	}

	return setupData, nil
}
