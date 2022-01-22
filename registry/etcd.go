package registry

import (
	"context"
	"fmt"

	"github.com/yixinin/flex/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterAddr(ctx context.Context) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: etcdConfig.Endpoints,
	})
	if err != nil {
		panic(err)
	}

	resp, err := client.Grant(ctx, 1)
	if err != nil {
		panic(err)
	}

	var addr string
	var key = fmt.Sprintf("%s/addr/%s", etcdConfig.App, primitive.NewObjectID().Hex())
	_, err = client.Put(ctx, key, addr, clientv3.WithLease(resp.ID))
	if err != nil {
		panic(err)
	}
	logger.Infof(ctx, "register app key:%s addr:%s", key, addr)

	ch, err := client.KeepAlive(ctx, resp.ID)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case c := <-ch:
			if c == nil {
				return
			}
			logger.Infof(ctx, "lease keep alive:%v", c.TTL)
		}
	}
}
