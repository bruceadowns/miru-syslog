package lib

// MakoJSON holds mako structured json
type MakoJSON struct {
	Timestamp          string `json:"@timestamp,omitempty"`
	Version            int    `json:"@version,omitempty"`
	Message            string `json:"message,omitempty"`
	LoggerName         string `json:"logger_name,omitempty"`
	ThreadName         string `json:"thread_name,omitempty"`
	Level              string `json:"level,omitempty"`
	LevelValue         int    `json:"level_value,omitempty"`
	ServiceName        string `json:"service_name,omitempty"`
	ServiceEnvironment string `json:"service_environment,omitempty"`
	ServicePipeline    string `json:"service_pipeline,omitempty"`
	ServiceVersion     string `json:"service_version,omitempty"`
}
