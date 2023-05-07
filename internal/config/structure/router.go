package structure

// Router 路由配置. 包含服务发现和服务注册
type Router struct {
	ZookeeperAddr string `mapstructure:"zookeeper_addr" yaml:"zookeeper_addr" json:"zookeeper_addr"` // zookeeper 地址, 集群使用 | 进行分割. 用来作为服务发现和服务注册
	BusAddr       string `mapstructure:"bus_addr" yaml:"bus_addr" json:"bus_addr"`                   // BusAddr, 用来作为数据总线. 当前系统支持 Ants.  集群使用 | 进行分割
}
