package structure

// HttpServer Http服务 特有配置
type HttpServer struct {
	ListenPort  string `yaml:"listen_port" mapstructure:"listen_port" json:"listen_port"`       // HTTP服监听端口
	GameDataDir string `yaml:"game_data_dir" mapstructure:"game_data_dir" json:"game_data_dir"` // 游戏配置表目录地址
}
