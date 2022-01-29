package etcd

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/yixinin/flex/addrs"
	"github.com/yixinin/flex/client/addr"
	"github.com/yixinin/flex/client/config"
	"github.com/yixinin/flex/client/event"
	"github.com/yixinin/flex/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Config struct {
	App       string   `mapstructure:"app"`
	Endpoints []string `mapstructure:"endpoints"`
}

func (c Config) Check() bool {
	if c.App == "" || len(c.Endpoints) == 0 {
		return false
	}
	return true
}

type EtcdAddrManager struct {
	client *clientv3.Client
	app    string
	event  chan event.AddrEvent
	addrs  map[string]*net.TCPAddr
	cancel func()
}

func NewEtcdAddrManager(conf config.Config) addr.AddrManager {
	c, _ := conf.(Config)
	m := &EtcdAddrManager{
		app:   c.App,
		event: make(chan event.AddrEvent, 255),
	}
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   c.Endpoints,
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	m.client = client
	return m
}

func (m *EtcdAddrManager) Event() chan event.AddrEvent {
	return m.event
}

func (c *EtcdAddrManager) close(ctx context.Context) {
	if c.cancel != nil {
		c.cancel()
	}
	c.client.Close()
	close(c.event)
}

func (c *EtcdAddrManager) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		recover()
		cancel()
		c.close(ctx)
	}()

	c.cancel = cancel
	c.watch(ctx)
	return nil
}

func (c *EtcdAddrManager) watch(ctx context.Context) {
	var keyPrefix = fmt.Sprintf("flex/%s", c.app)

	resp, err := c.client.Get(ctx, keyPrefix, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	for _, kv := range resp.Kvs {
		id, addr := addrs.ParseKv(ctx, kv.Key, kv.Value)
		if id == "" || addr == nil {
			logger.Warnf(ctx, "get wrong kv:%v", kv)
			continue
		}
		c.addrs[id] = addr
		c.event <- event.AddrEvent{
			EventType: event.EventAdd,
			Id:        id,
			Addr:      addr,
		}
	}

	ch := c.client.Watch(ctx, keyPrefix, clientv3.WithPrefix())
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			for _, ev := range msg.Events {
				id, addr := addrs.ParseKv(ctx, ev.Kv.Key, ev.Kv.Value)
				if id == "" || addr == nil {
					logger.Warnf(ctx, "recv wrong event:%v", ev)
					continue
				}
				switch ev.Type {
				case clientv3.EventTypePut:
					c.addrs[id] = addr
					c.event <- event.AddrEvent{
						EventType: event.EventAdd,
						Id:        id,
						Addr:      addr,
					}
				case clientv3.EventTypeDelete:
					delete(c.addrs, id)
					c.event <- event.AddrEvent{
						EventType: event.EventDel,
						Id:        id,
						Addr:      addr,
					}
				}
			}
		}
	}
}
