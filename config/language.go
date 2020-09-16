package config

import (
	JuderModels "Rabbit-OJ-Backend/services/judger/models"
)

var (
	SupportLanguage []JuderModels.SupportLanguage
	CompileObject   map[string]JuderModels.CompileInfo
	LocalImages     map[string]bool
)
