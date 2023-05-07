package structure

// MainServer 主逻辑服 特有配置
type MainServer struct {
	GameDataDir        string `mapstructure:"game_data_dir" yaml:"game_data_dir" json:"game_data_dir"`                      // 游戏配置表目录地址
	SensitiveWordsFile string `mapstructure:"sensitive_words_file" yaml:"sensitive_words_file" json:"sensitive_words_file"` // 游戏屏蔽字文件地址
}
