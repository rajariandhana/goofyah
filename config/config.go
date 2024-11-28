// config.go
package config

// Define a structure for the config with Templates as a nested struct
type Config struct {
	Templates TemplatesConfig `json:"templates"`
}

// Define the structure for the Templates config
type TemplatesConfig struct {
	Path string `json:"path"`
}

// Optionally, you could load the config from a JSON file or environment variables

func LoadConfig() {

}
