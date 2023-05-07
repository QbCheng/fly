package structure

// GateServer 网关服务配置
type GateServer struct {
	ListenPort string `mapstructure:"listen_port" yaml:"listen_port" json:"listen_port"` // 网关服监听端口
}
