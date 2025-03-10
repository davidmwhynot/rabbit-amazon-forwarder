package config

const (
	// MappingFile mapping file environment variable
	MappingFile      = "MAPPING_FILE"
	CaCertFile       = "CA_CERT_FILE"
	CertFile         = "CERT_FILE"
	KeyFile          = "KEY_FILE"
	LambdaEndpoint   = "LAMBDA_ENDPOINT"
	LambdaDisableSSL = "DISABLE_SSL"
)

// RabbitEntry RabbitMQ mapping entry
type RabbitEntry struct {
	Type          string   `json:"type"`
	Name          string   `json:"name"`
	ConnectionURL string   `json:"connection"`
	ExchangeName  string   `json:"topic"`
	QueueName     string   `json:"queue"`
	RoutingKey    string   `json:"routing"`
	RoutingKeys   []string `json:"routingKeys"`
}

// AmazonEntry SQS/SNS mapping entry
type AmazonEntry struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Target string `json:"target"`
}
