package models

type JudgerConfigType struct {
	Kafka       KafkaConfig    `json:"kafka"`
	Rpc         string         `json:"rpc"`
	AutoRemove  AutoRemoveType `json:"auto_remove"`
	Concurrent  ConcurrentType `json:"concurrent"`
	BuildImages []string       `json:"local_images"`
	Languages   []LanguageType `json:"languages"`
	Extensions  ExtensionsType `json:"extensions"`
}

type ExtensionsType struct {
	HostBind   bool           `json:"host_bind"`
	AutoPull   bool           `json:"auto_pull"`
	CheckJudge CheckJudgeType `json:"check_judge"`
	Expire     ExpireType     `json:"expire"`
}

type ExpireType struct {
	Enabled  bool  `json:"enabled"`
	Interval int64 `json:"interval"` // interval: minutes
}

type CheckJudgeType struct {
	Enabled  bool  `json:"enabled"`
	Interval int64 `json:"interval"` // interval: minutes
	Requeue  bool  `json:"requeue"`
}

type LanguageType struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Enabled bool        `json:"enabled"`
	Args    CompileInfo `json:"args"`
}

type ConcurrentType struct {
	Judge uint `json:"judge"`
}

type AutoRemoveType struct {
	Containers bool `json:"containers"`
	Files      bool `json:"files"`
}

type CompileInfo struct {
	BuildArgs   []string    `json:"build_args"`
	Source      string      `json:"source"`
	NoBuild     bool        `json:"no_build"`
	BuildTarget string      `json:"build_target"`
	BuildImage  string      `json:"build_image"`
	Constraints Constraints `json:"constraints"`
	RunArgs     []string    `json:"run_args"`
	RunArgsJSON string      `json:"-"`
	RunImage    string      `json:"run_image"`
}

type Constraints struct {
	BuildTimeout int   `json:"build_timeout"` // unit:seconds
	RunTimeout   int   `json:"run_timeout"`   // unit: seconds
	CPU          int64 `json:"cpu"`           // unit: COREs / 1e9
	Memory       int64 `json:"memory"`        // unit: bytes
}

type KafkaConfig struct {
	Brokers []string `json:"brokers"`
}
