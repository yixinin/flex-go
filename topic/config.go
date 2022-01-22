package topic

type Config struct {
	Topic      string `mapstructure:"topic"`
	RouterName string `mapstructure:"router"`
	BufferName string `mapstructure:"buffer"`
}
