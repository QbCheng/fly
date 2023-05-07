package structure

// ConnServer 连接服务 特有配置
type ConnServer struct {
	ListenPort       string `mapstructure:"listen_port" yaml:"listen_port" json:"listen_port"`                      // 连接服监听端口
	MetadataConnAddr string `mapstructure:"metadata_conn_addr" yaml:"metadata_conn_addr" json:"metadata_conn_addr"` // 连接服节点元数据 -- 连接地址. 会在节点信息中保存, 并且转发到所有的节点中
	ConnTyp          int    `mapstructure:"conn_typ" yaml:"conn_typ" json:"conn_typ"`                               // 连接类型
	CertFile         string `mapstructure:"cert_file" yaml:"cert_file" json:"cert_file"`
	KeyFile          string `mapstructure:"key_file" yaml:"key_file" json:"key_file"`
}
