package spec

type Parameter struct {
	APIVersion int
	Type       ServiceType
	Module     string
}

type ServiceType string

const (
	ServiceHost    ServiceType = "host"
	ServicePlugin  ServiceType = "plugin"
	ServiceUnknown ServiceType = "unknown"
	ServiceNone    ServiceType = "none"
	EnvModuleName              = "env"
)
