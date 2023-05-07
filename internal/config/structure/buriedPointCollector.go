package structure

// BuriedPointCollector 数据埋点收集器
type BuriedPointCollector struct {
	KafkaAddr string `mapstructure:"kafka_addr" yaml:"kafka_addr" json:"kafka_addr"` // kafka 地址 集群使用 | 进行分割
}
