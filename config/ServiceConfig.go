package config

import "strings"

type serviceConfig struct {
	Name        string `toml:"name"`
	DisplayName string `toml:"displayName"`
	Description string `toml:"description"`
}

//格式化
func (sc *serviceConfig) FormatConfig() {
	sc.Name = strings.Trim(sc.Name, " ")
	sc.DisplayName = strings.Trim(sc.DisplayName, " ")
	sc.Description = strings.Trim(sc.Description, " ")
	if sc.Name == "" {
		sc.Name = "GoAgentWin"
	}
	if sc.DisplayName == "" {
		sc.DisplayName = "GoAgentWin"
	}
	if sc.Description == "" {
		sc.Description = sc.Name
	}
}
