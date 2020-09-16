package config

import (
	JuderModels "github.com/Rabbit-OJ/Rabbit-OJ-Judger/models"
)

var (
	SupportLanguage []JuderModels.SupportLanguage
	CompileObject   map[string]JuderModels.CompileInfo
	LocalImages     map[string]bool
)
