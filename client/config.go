package client

type Config struct {
	App       string   `mapstructure:"app"`
	Topic     string   `mapstructure:"topic"`
	Pubsub    string   `mapstructure:"pubsub"` // pub or sub
	Endpoints []string `mapstructure:"endpoints"`
}
