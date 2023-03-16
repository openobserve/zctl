package utils

func Setup(releaseName string, namespace string) {
	SetupAWS(releaseName)
	SetupHelm(releaseName, namespace)

}
