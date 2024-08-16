package config

import (
	"github.com/marcodd23/go-micro-core/pkg/configmgr"
	"log"
)

// ServiceConfig - Application Level config.
// embed configmgr.BaseConfig
type ServiceConfig struct {
	configmgr.BaseConfig `mapstructure:",squash"`
	Rest                 Rest `yaml:"rest"`
}

// Rest configuration
type Rest struct {
	Endpoints map[string]Endpoint `yaml:"endpoints"`
}

// Endpoint configuration
type Endpoint struct {
	Method string `yaml:"method"`
	Path   string `yaml:"path"`
}

// LoadConfiguration - It load the property-<ENV>.yaml into the ServiceConfig struct.
func LoadConfiguration() *ServiceConfig {
	var cfg ServiceConfig

	err := configmgr.LoadConfigForEnv(&cfg)
	if err != nil {
		log.Panicf("error loading property files: %+v", err)
	}

	return &cfg
}
