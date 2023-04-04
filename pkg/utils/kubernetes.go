package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateConfigMap(sData SetupData) error {
	name := "zincobserve-setup"

	// var data = make(map[string]string)

	dataBytes, err := json.Marshal(sData)
	if err != nil {
		return err
	}

	// convert the dataBytes to map[string]string
	data := map[string]string{
		"data": string(dataBytes),
	}

	// Use the default kubeconfig file to create a Config object.
	kubeconfig, err := clientcmd.LoadFromFile(clientcmd.RecommendedHomeFile)
	if err != nil {
		return err
	}

	// create a configmap using the kubeconfig retrieved earlier
	config, err := clientcmd.NewDefaultClientConfig(*kubeconfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// create the ConfigMap object
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: sData.Namespace,
		},
		Data: data,
	}

	// create the ConfigMap
	_, err = clientset.CoreV1().ConfigMaps(sData.Namespace).Create(context.Background(), cm, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Println("ConfigMap created successfully")
	return nil
}

func ReadConfigMap(name string, namespace string) (SetupData, error) {

	setupData := SetupData{}
	// Use the default kubeconfig file to create a Config object.
	kubeconfig, err := clientcmd.LoadFromFile(clientcmd.RecommendedHomeFile)
	if err != nil {
		return setupData, err
	}

	// create a configmap using the kubeconfig retrieved earlier
	config, err := clientcmd.NewDefaultClientConfig(*kubeconfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return setupData, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return setupData, err
	}

	// Read the ConfigMap
	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return setupData, err
	}

	// marshal the configma data object into setupData
	err = json.Unmarshal([]byte(cm.Data["data"]), &setupData)
	if err != nil {
		return setupData, err
	}

	fmt.Println("ConfigMap read successfully")

	return setupData, nil
}

func DeleteConfigMap(name string, namespace string) error {
	// Use the default kubeconfig file to create a Config object.
	kubeconfig, err := clientcmd.LoadFromFile(clientcmd.RecommendedHomeFile)
	if err != nil {
		return err
	}

	// create a configmap using the kubeconfig retrieved earlier
	config, err := clientcmd.NewDefaultClientConfig(*kubeconfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// Delete the ConfigMap
	err = clientset.CoreV1().ConfigMaps(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	fmt.Println("ConfigMap deleted successfully")

	return nil
}

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

func GetCurrentNamespace() (string, error) {
	// Use the default kubeconfig file to create a Config object.
	kubeconfig, err := clientcmd.LoadFromFile(clientcmd.RecommendedHomeFile)
	if err != nil {
		return "", err
	}

	// Use the client to retrieve the current context.

	return kubeconfig.Contexts[kubeconfig.CurrentContext].Namespace, nil

}

func GetReleaseIdentifierFromReleaseName(releaseName string) string {
	return releaseName + "-release-identifier"
}
