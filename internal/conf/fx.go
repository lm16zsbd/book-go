package conf

import (
	"os"

	"go.uber.org/fx"
)

var Module = fx.Module("conf",
	fx.Provide(LoadConfig),
)

func LoadConfig() (*Bootstrap, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	return Load(configPath)
}
