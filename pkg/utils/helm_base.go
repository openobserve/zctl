package utils

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/postrender"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"helm.sh/helm/v3/pkg/strvals"
)

type Helm struct {
	AppVersion    string
	ChartName     string
	ChartVersion  string
	Namespace     string
	PostRenderer  postrender.PostRenderer
	ReleaseName   string
	RepositoryURL string
	Wait          bool

	SetValues  []string
	ValuesFile string
}

// initialize creates and initializes a new Helm action configuration object with the specified Kubernetes context and namespace.
// It returns a pointer to the Configuration object or an error if one occurs.
func initialize(kubeContext, namespace string) (*action.Configuration, error) {
	// Workaround for https://github.com/helm/helm/issues/7430.
	_ = os.Setenv("HELM_KUBECONTEXT", kubeContext)
	_ = os.Setenv("HELM_NAMESPACE", namespace)

	// Create a new CLI settings object.
	settings := cli.New()

	// Initialize the action configuration.
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "secret", log.Printf); err != nil {
		return nil, fmt.Errorf("failed to initialize helm action config: %w", err)
	}

	return actionConfig, nil
}

// DownloadChart downloads the specified Helm chart from the repository URL and chart version, and returns a pointer to the Chart object or an error if one occurs.
func (h *Helm) DownloadChart() (*chart.Chart, error) {
	// Create a new list of getters for fetching remote chart repositories.
	getters := getter.All(&cli.EnvSettings{})

	// Parse the repository URL.
	u, err := url.Parse(h.RepositoryURL)
	if err != nil {
		return nil, err
	}

	var chartPath string
	if u.Scheme == "oci" {
		// Use the repository URL and chart version for OCI artifacts.
		chartPath = h.RepositoryURL + ":" + h.ChartVersion
	} else {
		// Find the chart in the repository.
		chartPath, err = repo.FindChartInRepoURL(h.RepositoryURL, h.ChartName, h.ChartVersion, "", "", "", getters)
		if err != nil {
			return nil, err
		}
	}

	// Print a log message for the chart being downloaded.
	fmt.Printf("Downloading Chart: %s\n", chartPath)

	// Get the appropriate getter for the repository URL scheme.
	g, err := getters.ByScheme(u.Scheme)
	if err != nil {
		return nil, err
	}

	// Download the chart archive into memory.
	data, err := g.Get(chartPath)
	if err != nil {
		if strings.HasPrefix(h.RepositoryURL, "oci://public.ecr.aws") {
			msg := "Please review: https://docs.aws.amazon.com/AmazonECR/latest/public/public-troubleshooting.html"
			err = fmt.Errorf("%w\n%s", err, msg)
		}
		return nil, err
	}

	// Decompress the chart archive.
	files, err := loader.LoadArchiveFiles(data)
	if err != nil {
		return nil, err
	}

	// Load the chart.
	chart, err := loader.LoadFiles(files)
	if err != nil {
		return nil, err
	}

	return chart, nil
}

// Install deploys the specified Helm chart with the given parameters, and returns an error if one occurs.
func (h *Helm) Install(chart *chart.Chart, kubeContext string) error {
	// Parse the values file.
	values := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(h.ValuesFile), &values); err != nil {
		return fmt.Errorf("failed to parse values file: %w", err)
	}

	// Parse the set values and add them to the values map.
	for _, v := range h.SetValues {
		if err := strvals.ParseInto(v, values); err != nil {
			return fmt.Errorf("failed parsing --set data: %w", err)
		}
	}

	// Initialize the Helm action configuration.
	actionConfig, err := initialize(kubeContext, h.Namespace)
	if err != nil {
		return err
	}

	// Configure the Helm install options.
	instAction := action.NewInstall(actionConfig)
	instAction.Namespace = h.Namespace
	instAction.ReleaseName = h.ReleaseName
	instAction.CreateNamespace = true
	instAction.IsUpgrade = true
	instAction.PostRenderer = h.PostRenderer
	instAction.Wait = h.Wait
	instAction.Timeout = 300 * time.Second
	chart.Metadata.AppVersion = h.AppVersion

	// Install the chart.
	fmt.Println("Installing using helm chart...")
	rel, err := instAction.Run(chart, values)
	if err != nil {
		return fmt.Errorf("helm install failed: %s", err)
	}

	// Print the chart installation details.
	fmt.Printf("Using chart version %q, installed %q version %q in namespace %q\n",
		rel.Chart.Metadata.Version, rel.Name, rel.Chart.Metadata.AppVersion, rel.Namespace)

	// Print the chart NOTES.
	if len(rel.Info.Notes) > 0 {
		fmt.Printf("NOTES:\n%s\n", strings.TrimSpace(rel.Info.Notes))
	}

	return nil
}

func (h *Helm) UnInstall(releaseName, namespace string) error {

	kubeConfig := cli.New()

	// Set up the Helm action configuration.
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(kubeConfig.RESTClientGetter(), kubeConfig.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Fatalf("Failed to initialize Helm action configuration: %v", err)
	}

	if err := deleteHelmRelease(releaseName, actionConfig); err != nil {
		log.Fatalf("Error deleting Helm release: %v", err)
	}

	return nil
}

func deleteHelmRelease(releaseName string, actionConfig *action.Configuration) error {
	deleteAction := action.NewUninstall(actionConfig)
	res, err := deleteAction.Run(releaseName)
	if err != nil {
		return fmt.Errorf("failed to delete Helm release: %w", err)
	}
	fmt.Printf("Release %s deleted: %s\n", releaseName, res.Release.Info.Description)
	return nil
}

// List returns a list of all installed releases in the specified Kubernetes cluster.
func List(kubeContext string) ([]*release.Release, error) {
	// Initialize the Helm action configuration.
	actionConfig, err := initialize(kubeContext, "")
	if err != nil {
		return nil, err
	}

	// Configure the Helm list options.
	client := action.NewList(actionConfig)
	client.AllNamespaces = true

	// List all installed releases.
	releases, err := client.Run()
	if (err) != nil {
		return nil, err
	}

	return releases, nil
}

// Status returns the status of the specified release in the specified Kubernetes cluster.
func Status(kubeContext, releaseName, namespace string) (string, error) {
	// Initialize the Helm action configuration.
	actionConfig, err := initialize(kubeContext, namespace)
	if err != nil {
		return "", err
	}

	// Configure the Helm status options.
	status := action.NewStatus(actionConfig)

	// Get the status of the specified release.
	rel, err := status.Run(releaseName)
	if (err) != nil {
		return "", err
	}

	// Strip chart metadata from the output.
	rel.Chart = nil

	return "", nil
}

// Uninstall uninstalls the specified release from the specified Kubernetes cluster and namespace.
func Uninstall(kubeContext, releaseName, namespace string) error {
	// Initialize the Helm action configuration.
	actionConfig, err := initialize(kubeContext, namespace)
	if err != nil {
		return err
	}

	// Configure the Helm uninstall options.
	uninstall := action.NewUninstall(actionConfig)

	// Uninstall the specified release.
	_, err = uninstall.Run(releaseName)
	if err != nil {
		return fmt.Errorf("failed uninstalling chart: %w", err)
	}

	return nil
}
