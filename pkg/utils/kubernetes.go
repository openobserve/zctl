package utils

import (
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func Client(context string) (*kubernetes.Clientset, error) {
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		},
	)

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(restConfig)
}

func DynamicClient(context string) (dynamic.Interface, error) {
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		},
	)

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return dynamic.NewForConfig(restConfig)
}

func KubeContextForCluster(clusterEndpoint string) (string, error) {
	raw, err := Kubeconfig()
	if err != nil {
		return "", err
	}

	found := ""

	for name, context := range raw.Contexts {
		if _, ok := raw.Clusters[context.Cluster]; ok {
			if raw.Clusters[context.Cluster].Server == clusterEndpoint {
				found = name
				break
			}
		}
	}

	return found, nil
}

func Kubeconfig() (*clientcmdapi.Config, error) {
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	raw, err := config.RawConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	return &raw, nil
}

func GetCurrentKubeContextAPIEndpoint() (string, error) {
	// Use the default kubeconfig file to create a Config object.
	kubeconfig, err := clientcmd.LoadFromFile(clientcmd.RecommendedHomeFile)
	if err != nil {
		return "", err
	}

	// Use the client to retrieve the current context.
	currentContext := kubeconfig.CurrentContext
	context := kubeconfig.Contexts[currentContext]

	// Use the context to retrieve the API server endpoint.
	cluster := kubeconfig.Clusters[context.Cluster]
	return cluster.Server, nil
}
