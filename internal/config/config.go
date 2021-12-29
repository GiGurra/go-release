package config

type AppConfig struct {
	Module                   string
	Version                  string
	IgnoreUncommittedChanges bool
}

func GetDefaultAppConfig() AppConfig {
	return AppConfig{
		Module:                   "",
		Version:                  "",
		IgnoreUncommittedChanges: false,
	}
}
