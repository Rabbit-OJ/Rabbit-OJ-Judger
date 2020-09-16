package models

type JudgerConfigType struct {
	Kafka       kafkaConfig `json:"kafka"`
	Rpc         string      `json:"rpc"`
	AutoRemove  autoRemove  `json:"auto_remove"`
	Concurrent  concurrent  `json:"concurrent"`
	LocalImages []string    `json:"local_images"`
	Languages   []language  `json:"languages"`
	Extensions  extensions  `json:"extensions"`
}

type extensions struct {
	AutoPull   bool       `json:"auto_pull"`
	CheckJudge checkJudge `json:"check_judge"`
	Expire     expire     `json:"expire"`
}

type expire struct {
	Enabled  bool  `json:"enabled"`
	Interval int64 `json:"interval"` // interval: minutes
}

type checkJudge struct {
	Enabled  bool  `json:"enabled"`
	Interval int64 `json:"interval"` // interval: minutes
	Requeue  bool  `json:"requeue"`
}

type language struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Enabled bool        `json:"enabled"`
	Args    CompileInfo `json:"args"`
}

type concurrent struct {
	Judge uint `json:"judge"`
}

type autoRemove struct {
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

type kafkaConfig struct {
	Brokers []string `json:"brokers"`
}
