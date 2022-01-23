package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/yixinin/flex/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (c *Client) Watch(ctx context.Context) {
	defer c.etcdClient.Close()

	var keyPrefix = fmt.Sprintf("flex/%s", c.app)

	resp, err := c.etcdClient.Get(ctx, keyPrefix, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	for _, kv := range resp.Kvs {
		id, addr := parseKv(ctx, kv.Key, kv.Value)
		if id == "" || addr == "" {
			continue
		}
		c.AddAddr(ctx, id, addr)
	}

	ch := c.etcdClient.Watch(ctx, keyPrefix, clientv3.WithPrefix())
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			for _, ev := range msg.Events {
				id, addr := parseKv(ctx, ev.Kv.Key, ev.Kv.Value)
				if id == "" || addr == "" {
					continue
				}
				switch ev.Type {
				case clientv3.EventTypePut:
					c.AddAddr(ctx, id, addr)
				case clientv3.EventTypeDelete:
					c.DelAddr(ctx, id)
				}
			}
		}
	}
}

func parseKv(ctx context.Context, key, val []byte) (id, addr string) {
	keys := strings.Split(string(key), "/")
	if len(keys) != 2 || len(val) == 0 {
		logger.Warnf(ctx, "unknown key:%s val:%s", key, val)
		return
	}
	id = keys[2]
	addr = string(val)
	return
}
