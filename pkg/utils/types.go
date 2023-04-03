package utils

type ZincObserveValues struct {
	Image              Image              `yaml:"image"`
	ImagePullSecrets   []string           `yaml:"imagePullSecrets"`
	NameOverride       string             `yaml:"nameOverride"`
	FullnameOverride   string             `yaml:"fullnameOverride"`
	ServiceAccount     ServiceAccount     `yaml:"serviceAccount"`
	PodSecurityContext PodSecurityContext `yaml:"podSecurityContext"`
	SecurityContext    SecurityContext    `yaml:"securityContext"`
	ReplicaCount       ReplicaCount       `yaml:"replicaCount"`
	Auth               Auth               `yaml:"auth"`
	Config             Config             `yaml:"config"`
	Service            Service            `yaml:"service"`
	Ingress            Ingress            `yaml:"ingress"`
	Resources          Resources          `yaml:"resources"`
	Autoscaling        Autoscaling        `yaml:"autoscaling"`
	CertIssuer         CertIssuer         `yaml:"certIssuer"`
	Ingester           Ingester           `yaml:"ingester"`
	Etcd               Etcd               `yaml:"etcd"`
	MinIO              MinIO              `yaml:"minio"`
}

type Ingester struct {
	Persistence Persistence `yaml:"persistence"`
}

type CertIssuer struct {
	Enabled bool `yaml:"enabled"`
}

type Image struct {
	Registry   string `yaml:"registry"`
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
	PullPolicy string `yaml:"pullPolicy"`
}

type ServiceAccount struct {
	Create      bool              `yaml:"create"`
	Name        string            `yaml:"name"`
	Annotations map[string]string `yaml:"annotations"`
}
type PodSecurityContext struct {
	FSGroup      int  `yaml:"fsGroup"`
	RunAsUser    int  `yaml:"runAsUser"`
	RunAsGroup   int  `yaml:"runAsGroup"`
	RunAsNonRoot bool `yaml:"runAsNonRoot"`
}

type SecurityContext struct {
	ReadOnlyRootFilesystem bool `yaml:"readOnlyRootFilesystem"`
	RunAsNonRoot           bool `yaml:"runAsNonRoot"`
	RunAsUser              int  `yaml:"runAsUser"`
	Capabilities           struct {
		Drop []string `yaml:"drop"`
	} `yaml:"capabilities"`
}

type ReplicaCount struct {
	Ingester     int `yaml:"ingester"`
	Querier      int `yaml:"querier"`
	Router       int `yaml:"router"`
	Alertmanager int `yaml:"alertmanager"`
	Compactor    int `yaml:"compactor"`
}

type Auth struct {
	ZO_ROOT_USER_EMAIL    string `yaml:"ZO_ROOT_USER_EMAIL"`
	ZO_ROOT_USER_PASSWORD string `yaml:"ZO_ROOT_USER_PASSWORD"`
	ZOS3ACCESSKEY         string `yaml:"ZO_S3_ACCESS_KEY"`
	ZOS3SECRETKEY         string `yaml:"ZO_S3_SECRET_KEY"`
}

type Config struct {
	ZOLOCALMODE                     string `yaml:"ZO_LOCAL_MODE"`
	ZOHTTPPORT                      string `yaml:"ZO_HTTP_PORT"`
	ZOGRPCPORT                      string `yaml:"ZO_GRPC_PORT"`
	ZOGRPCTIMEOUT                   string `yaml:"ZO_GRPC_TIMEOUT"`
	ZOGRPCORGHEADERKEY              string `yaml:"ZO_GRPC_ORG_HEADER_KEY"`
	ZOROUTETIMEOUT                  string `yaml:"ZO_ROUTE_TIMEOUT"`
	ZOLOCALMODESTORAGE              string `yaml:"ZO_LOCAL_MODE_STORAGE"`
	ZONODEROLE                      string `yaml:"ZO_NODE_ROLE"`
	ZOINSTANCENAME                  string `yaml:"ZO_INSTANCE_NAME"`
	ZODATADIR                       string `yaml:"ZO_DATA_DIR"`
	ZODATAWALDIR                    string `yaml:"ZO_DATA_WAL_DIR"`
	ZODATASTREAMDIR                 string `yaml:"ZO_DATA_STREAM_DIR"`
	ZOWALMEMORYMODEENABLED          string `yaml:"ZO_WAL_MEMORY_MODE_ENABLED"`
	ZOFILEEXTJSON                   string `yaml:"ZO_FILE_EXT_JSON"`
	ZOFILEEXTPARQUET                string `yaml:"ZO_FILE_EXT_PARQUET"`
	ZOPARQUETCOMPRESSION            string `yaml:"ZO_PARQUET_COMPRESSION"`
	ZOTIMESTAMPCOL                  string `yaml:"ZO_TIME_STAMP_COL"`
	ZOWIDENINGSCHEMAEVOLUTION       string `yaml:"ZO_WIDENING_SCHEMA_EVOLUTION"`
	ZOFEATUREPERTHREADLOCK          string `yaml:"ZO_FEATURE_PER_THREAD_LOCK"`
	ZOFEATUREFULLTEXTONALLFIELDS    string `yaml:"ZO_FEATURE_FULLTEXT_ON_ALL_FIELDS"`
	ZOUIENABLED                     string `yaml:"ZO_UI_ENABLED"`
	ZOMETRICSDEDUPENABLED           string `yaml:"ZO_METRICS_DEDUP_ENABLED"`
	ZOTRACINGENABLED                string `yaml:"ZO_TRACING_ENABLED"`
	OTELOTLPHTTPENDPOINT            string `yaml:"OTEL_OTLP_HTTP_ENDPOINT"`
	ZOTRACINGHEADERKEY              string `yaml:"ZO_TRACING_HEADER_KEY"`
	ZOTRACINGHEADERVALUE            string `yaml:"ZO_TRACING_HEADER_VALUE"`
	ZOTELEMETRY                     string `yaml:"ZO_TELEMETRY"`
	ZOTELEMETRYURL                  string `yaml:"ZO_TELEMETRY_URL"`
	ZOJSONLIMIT                     string `yaml:"ZO_JSON_LIMIT"`
	ZOPAYLOADLIMIT                  string `yaml:"ZO_PAYLOAD_LIMIT"`
	ZOMAXFILESIZEONDISK             string `yaml:"ZO_MAX_FILE_SIZE_ON_DISK"`
	ZOMAXFILERETENTIONTIME          string `yaml:"ZO_MAX_FILE_RETENTION_TIME"`
	ZOFILEPUSHINTERVAL              string `yaml:"ZO_FILE_PUSH_INTERVAL"`
	ZOFILEMOVETHREADNUM             string `yaml:"ZO_FILE_MOVE_THREAD_NUM"`
	ZOQUERYTHREADNUM                string `yaml:"ZO_QUERY_THREAD_NUM"`
	ZOTSALLOWEDUPTO                 string `yaml:"ZO_TS_ALLOWED_UPTO"`
	ZOMETRICSLEADERPUSHINTERVAL     string `yaml:"ZO_METRICS_LEADER_PUSH_INTERVAL"`
	ZOMETRICSLEADERELECTIONINTERVAL string `yaml:"ZO_METRICS_LEADER_ELECTION_INTERVAL"`
	ZOHEARTBEATINTERVAL             string `yaml:"ZO_HEARTBEAT_INTERVAL"`
	ZOCOMPACTENABLED                string `yaml:"ZO_COMPACT_ENABLED"`
	ZOCOMPACTINTERVAL               string `yaml:"ZO_COMPACT_INTERVAL"`
	ZOCOMPACTMAXFILESIZE            string `yaml:"ZO_COMPACT_MAX_FILE_SIZE"`
	ZOMEMORYCACHEENABLED            string `yaml:"ZO_MEMORY_CACHE_ENABLED"`
	ZOMEMORYCACHECACHELATESTFILES   string `yaml:"ZO_MEMORY_CACHE_CACHE_LATEST_FILES"`
	ZOMEMORYCACHEMAXSIZE            string `yaml:"ZO_MEMORY_CACHE_MAX_SIZE"`
	ZOMEMORYCACHERELEASESIZE        string `yaml:"ZO_MEMORY_CACHE_RELEASE_SIZE"`
	RUSTLOG                         string `yaml:"RUST_LOG"`
	ZOCOLSPERRECORDLIMIT            string `yaml:"ZO_COLS_PER_RECORD_LIMIT"`
	ZOETCDPREFIX                    string `yaml:"ZO_ETCD_PREFIX"`
	ZOETCDCONNECTTIMEOUT            string `yaml:"ZO_ETCD_CONNECT_TIMEOUT"`
	ZOETCDCOMMANDTIMEOUT            string `yaml:"ZO_ETCD_COMMAND_TIMEOUT"`
	ZOETCDLOCKWAITTIMEOUT           string `yaml:"ZO_ETCD_LOCK_WAIT_TIMEOUT"`
	ZOETCDUSER                      string `yaml:"ZO_ETCD_USER"`
	ZOETCDPASSWORD                  string `yaml:"ZO_ETCD_PASSWORD"`
	ZOETCDCLIENTCERTAUTH            string `yaml:"ZO_ETCD_CLIENT_CERT_AUTH"`
	ZOETCDTRUSTEDCAFILE             string `yaml:"ZO_ETCD_TRUSTED_CA_FILE"`
	ZOETCDCERTFILE                  string `yaml:"ZO_ETCD_CERT_FILE"`
	ZOETCDKEYFILE                   string `yaml:"ZO_ETCD_KEY_FILE"`
	ZOETCDDOMAINNAME                string `yaml:"ZO_ETCD_DOMAIN_NAME"`
	ZOETCDLOADPAGESIZE              string `yaml:"ZO_ETCD_LOAD_PAGE_SIZE"`
	ZOSLEDDATADIR                   string `yaml:"ZO_SLED_DATA_DIR"`
	ZOSLEDPREFIX                    string `yaml:"ZO_SLED_PREFIX"`
	ZOS3SERVERURL                   string `yaml:"ZO_S3_SERVER_URL"`
	ZOS3REGIONNAME                  string `yaml:"ZO_S3_REGION_NAME"`
	ZOS3BUCKETNAME                  string `yaml:"ZO_S3_BUCKET_NAME"`
	ZOS3PROVIDER                    string `yaml:"ZO_S3_PROVIDER"`
	ZODATALIFECYCLE                 string `yaml:"ZO_DATA_LIFECYCLE"`
}

type Service struct {
	Type string `yaml:"type"`
	Port string `yaml:"port"`
}

type Ingress struct {
	Enabled     bool              `yaml:"enabled"`
	ClassName   string            `yaml:"className"`
	Annotations map[string]string `yaml:"annotations"`
	Hosts       []Host            `yaml:"hosts"`
	TLS         []TLS             `yaml:"tls"`
}

type Host struct {
	Host  string `yaml:"host"`
	Paths []Path `yaml:"paths"`
}

type Path struct {
	Path     string  `yaml:"path"`
	PathType string  `yaml:"pathType"`
	Backend  Backend `yaml:"backend"`
}

type Backend struct {
	Service string `yaml:"service"`
}

type BackendService struct {
	Name string      `yaml:"name"`
	Port BackendPort `yaml:"port"`
}

type BackendPort struct {
	Port string `yaml:"port"`
}

type TLS struct {
	Hosts      []string `yaml:"hosts"`
	SecretName string   `yaml:"secretName"`
}

type Resources struct {
	Limits   map[string]string `yaml:"limits"`
	Requests map[string]string `yaml:"requests"`
}

type Autoscaling struct {
	Enabled                           bool `yaml:"enabled"`
	MinReplicas                       int  `yaml:"minReplicas"`
	MaxReplicas                       int  `yaml:"maxReplicas"`
	TargetCPUUtilizationPercentage    int  `yaml:"targetCPUUtilizationPercentage"`
	TargetMemoryUtilizationPercentage int  `yaml:"targetMemoryUtilizationPercentage"`
}
type MinIO struct {
	Enabled       bool        `yaml:"enabled"`
	Region        string      `yaml:"region"`
	RootUser      string      `yaml:"rootUser"`
	RootPassword  string      `yaml:"rootPassword"`
	DrivesPerNode int         `yaml:"drivesPerNode"`
	Replicas      int         `yaml:"replicas"`
	Mode          string      `yaml:"mode"`
	Image         Image       `yaml:"image"`
	MCImage       Image       `yaml:"mcImage"`
	Buckets       []Bucket    `yaml:"buckets"`
	Resources     Resources   `yaml:"resources"`
	Persistence   Persistence `yaml:"persistence"`
}

type Bucket struct {
	Name   string `yaml:"name"`
	Policy string `yaml:"policy"`
	Purge  bool   `yaml:"purge"`
}

type Etcd struct {
	Enabled      bool        `yaml:"enabled"`
	ExternalUrl  string      `yaml:"externalUrl"`
	ReplicaCount int         `yaml:"replicaCount"`
	Image        Image       `yaml:"image"`
	ExtraEnvVars []NameValue `yaml:"extraEnvVars"`
	Persistence  Persistence `yaml:"persistence"`
	Auth         EtcdAuth    `yaml:"auth"`
	LogLevel     string      `yaml:"logLevel"`
}

type EtcdAuth struct {
	RBAC RBAC `yaml:"rbac"`
}

type RBAC struct {
	Create                  bool   `yaml:"create"`
	AllowNoneAuthentication bool   `yaml:"allowNoneAuthentication"`
	RootPassword            string `yaml:"rootPassword"`
}

type NameValue struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Persistence struct {
	Enabled           bool              `yaml:"enabled"`
	Size              string            `yaml:"size"`
	StorageClass      string            `yaml:"storageClass"`
	AccessModes       []string          `yaml:"accessModes"`
	Annotations       map[string]string `yaml:"annotations"`
	VolumePermissions VolumePermission  `yaml:"volumePermissions"`
}

type VolumePermission struct {
	Enabled bool `yaml:"enabled"`
}

type IAMPolicy struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

type Statement struct {
	Sid      string   `json:"Sid"`
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource []string `json:"Resource"`
}

type SetupData struct {
	Identifier     string `json:"identifier"`   // unique identifier generated randomly to avoid conflicts
	BucketName     string `json:"bucket_name"`  // s3 bucket name
	ReleaseName    string `json:"release_name"` // helm release name
	IamRole        string `json:"iam_role"`     // role name
	K8s            string `json:"k8s"`          // k8s cluster name eks, gke, plain
	S3AccessKey    string `json:"s3_access_key"`
	S3SecretKey    string `json:"s3_secret_key"`
	Namespace      string `json:"namespace"`
	Region         string `json:"region"`
	GCPProjectId   string `json:"gcp_project_id"`
	ClusterName    string `json:"cluster_name"`
	ServiceAccount string `json:"service_account"`
}
