/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zctl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	fmt.Println("init rootCmd called")
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.zctl.yaml)")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.myapp.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	// Bind viper values to the root command flags
	rootCmd.PersistentFlags().String("name", viper.GetString("metadata.name"), "name of the installation for underlying helm chart")
	viper.BindPFlag("metadata.name", rootCmd.Flags().Lookup("name"))
	rootCmd.MarkPersistentFlagRequired("name")

	rootCmd.PersistentFlags().String("k8s", viper.GetString("spec.k8s"), "k8s cluster type. eks, gke, plain")
	viper.BindPFlag("spec.k8s", rootCmd.Flags().Lookup("k8s"))
	rootCmd.MarkPersistentFlagRequired("k8s")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// if cfgFile != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// 	// Find home directory.
	// 	home, err := os.UserHomeDir()
	// 	cobra.CheckErr(err)

	// 	// Search config in home directory with name ".zctl" (without extension).
	// 	viper.AddConfigPath(home)
	// 	viper.SetConfigType("yaml")
	// 	viper.SetConfigName(".zctl")
	// }

	// viper.AutomaticEnv() // read in environment variables that match

	// // If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	// }
	if viper.GetString("config") != "" {
		// Use config file from the flag.
		viper.SetConfigFile(viper.GetString("config"))
	} else {
		// Search config in home directory with name ".myapp" (without extension).
		viper.AddConfigPath("$HOME")
		viper.SetConfigName(".myapp")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
