/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zinclabs/zctl/pkg/utils"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs ZincObserve",
	Long: `Installs ZincObserve.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Install Run called")
		// name := cmd.Flags().Lookup("name").Value.String()
		// namespace := cmd.Flags().Lookup("namespace").Value.String()
		namespace := cmd.Flags().String("namespace", viper.GetString("metadata.namespace"), "namespace to install the helm chart")
		k8s := cmd.Flags().String("k8s", viper.GetString("spec.k8s"), "k8s to install the helm chart")
		// k8s := cmd.Flags().Lookup("k8s").Value.String()
		install_minio := cmd.Flags().String("install_minio", viper.GetString("spec.install_minio"), "Specify if you want to install minio. Default is false.")
		// install_minio := cmd.Flags().Lookup("install_minio").Value.String()
		storage_provider := cmd.Flags().String("storage_provider", viper.GetString("spec.storage_provider"), "Valid values are s3, gcs, minio, swift.")
		// storage_provider := cmd.Flags().Lookup("storage_provider").Value.String()
		s3_access_key := cmd.Flags().String("s3_access_key", viper.GetString("spec.s3_access_key"), "S3 access key")
		// s3_access_key := cmd.Flags().Lookup("s3_access_key").Value.String()
		s3_secret_key := cmd.Flags().String("s3_secret_key", viper.GetString("spec.s3_secret_key"), "S3 secret key")
		// s3_secret_key := cmd.Flags().Lookup("s3_secret_key").Value.String()
		s3_server_url := cmd.Flags().String("s3_server_url", viper.GetString("spec.s3_server_url"), "S3 server url")
		// s3_server_url := cmd.Flags().Lookup("s3_server_url").Value.String()
		s3_bucket_name := cmd.Flags().String("s3_bucket_name", viper.GetString("spec.s3_bucket_name"), "S3 bucket name")
		// s3_bucket_name := cmd.Flags().Lookup("s3_bucket_name").Value.String()
		region := cmd.Flags().String("region", viper.GetString("spec.region"), "region to install the installation in.")
		// region := cmd.Flags().Lookup("region").Value.String()
		gcp_project_id := cmd.Flags().String("gcp_project_id", viper.GetString("spec.gcp_project_id"), "GCP Project ID to install the installation in.")
		// gcp_project_id := cmd.Flags().Lookup("gcp_project_id").Value.String()

		// get value for config
		name := viper.GetString("metadata.name")

		fmt.Println("name is: ", name)

		// convert install_minio to bool
		install_minio_bool, err := strconv.ParseBool(*install_minio)
		if err != nil {
			install_minio_bool = false
		}

		// Let's do the setup
		release_identifer := utils.GenerateReleaseIdentifier()

		inputData := utils.SetupData{
			Identifier:      release_identifer,
			ReleaseName:     name,
			Namespace:       *namespace,
			Region:          *region,
			K8s:             *k8s,
			GCPProjectId:    *gcp_project_id,
			S3AccessKey:     *s3_access_key,
			S3SecretKey:     *s3_secret_key,
			InstallMinIO:    install_minio_bool,
			StorageProvider: *storage_provider,
			S3ServerURL:     *s3_server_url,
			BucketName:      *s3_bucket_name,
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
	fmt.Println("init installCmd ")
	rootCmd.AddCommand(installCmd)

	fmt.Println("config file in installCmd: ", viper.ConfigFileUsed())

	installCmd.Flags().String("namespace", viper.GetString("metadata.namespace"), "namespace to install the helm chart")
	installCmd.Flags().String("region", viper.GetString("spec.region"), "region to install the installation in.")
	installCmd.Flags().String("gcp_project_id", viper.GetString("spec.gcp_project_id"), "GCP Project ID to install the installation in.")
	installCmd.Flags().String("install_minio", viper.GetString("spec.install_minio"), "Specify if you want to install minio. Default is false.")
	installCmd.Flags().String("storage_provider", viper.GetString("spec.storage_provider"), "Valid values are s3, gcs, minio, swift.")
	installCmd.Flags().String("s3_bucket_name", viper.GetString("spec.s3_bucket_name"), "s3 compatible bucket.")
	installCmd.Flags().String("s3_server_url", viper.GetString("spec.s3_server_url"), "s3 compatible server url.")
	installCmd.Flags().String("s3_access_key", viper.GetString("spec.s3_access_key"), "s3_access_key to use.")
	installCmd.Flags().String("s3_secret_key", viper.GetString("spec.s3_secret_key"), "s3_secret_key to use.")

	// Bind the flags to the configuration keys
	viper.BindPFlag("metadata.namespace", installCmd.Flags().Lookup("namespace"))
	viper.BindPFlag("spec.region", installCmd.Flags().Lookup("region"))
	viper.BindPFlag("spec.gcp_project_id", installCmd.Flags().Lookup("gcp_project_id"))
	viper.BindPFlag("spec.install_minio", installCmd.Flags().Lookup("install_minio"))
	viper.BindPFlag("spec.storage_provider", installCmd.Flags().Lookup("storage_provider"))
	viper.BindPFlag("spec.s3_bucket_name", installCmd.Flags().Lookup("s3_bucket_name"))
	viper.BindPFlag("spec.s3_server_url", installCmd.Flags().Lookup("s3_server_url"))
	viper.BindPFlag("spec.s3_access_key", installCmd.Flags().Lookup("s3_access_key"))
	viper.BindPFlag("spec.s3_secret_key", installCmd.Flags().Lookup("s3_secret_key"))

	// Bind the flags to the command
	installCmd.MarkFlagRequired("namespace")
	// installCmd.MarkFlagRequired("region")
	// installCmd.MarkFlagRequired("storage_provider")
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
