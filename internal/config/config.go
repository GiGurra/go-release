package config

type AppConfig struct {
	Module  string
	Version string
}

func GetDefaultAppConfig() AppConfig {
	return AppConfig{
		Module:  "",
		Version: "",
	}
}

func (config *AppConfig) Validate() {

}
