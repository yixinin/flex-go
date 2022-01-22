package registry

type EtcdConfig struct {
	App       string
	Endpoints []string
}

var etcdConfig EtcdConfig

func Init(cfg EtcdConfig) {
	etcdConfig = cfg
}
