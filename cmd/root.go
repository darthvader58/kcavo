package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	namespace     string
	allNamespaces bool
	output        string
)

var rootCmd = &cobra.Command{
	Use:   "kcavo",
	Short: "Kubernetes cost analysis and optimization tool",
	Long: `kubectl-cost is a kubectl plugin that helps you:
  • Visualize resource usage across your cluster
  • Analyze costs and spending patterns
  • Optimize GPU scheduling and allocation
  • Get cost-saving recommendations
  
Example usage:
  kubectl cost analyze                    # Analyze costs in current namespace
  kubectl cost analyze --all-namespaces   # Analyze cluster-wide costs
  kubectl cost visualize                  # Visualize resources
  kubectl cost gpu                        # Analyze GPU usage
  kubectl cost optimize                   # Get optimization recommendations`,
	Version: "1.0.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubectl-cost.yaml)")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "kubernetes namespace (default is current context namespace)")
	rootCmd.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "analyze across all namespaces")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "output format: table, json, yaml")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kcavo")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
func getNamespace() string {
	if allNamespaces {
		return ""
	}
	if namespace != "" {
		return namespace
	}
	return "default"
}
