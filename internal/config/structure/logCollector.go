package structure

// LogCollector 日志收集器配置
type LogCollector struct {
	KafkaAddr                     string `mapstructure:"kafka_addr" yaml:"kafka_addr" json:"kafka_addr"`                                                                         // kafka 地址 集群使用 | 进行分割
	KafkaGameLogTopic             string `mapstructure:"kafka_game_log_topic" yaml:"kafka_game_log_topic" json:"kafka_game_log_topic"`                                           // kafka 游戏日志主题
	KafkaGameLogConsumerGroupName string `mapstructure:"kafka_game_log_consumer_group_name" yaml:"kafka_game_log_consumer_group_name" json:"kafka_game_log_consumer_group_name"` // kafka 游戏日志 消费者组名字
}
