package registry

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/yixinin/flex/addrs"
	"github.com/yixinin/flex/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterAddr(ctx context.Context, port uint16) {
	client, err := clientv3.New(clientv3.Config{
		DialTimeout: 5 * time.Second,
		Endpoints:   etcdConfig.Endpoints,
	})
	if err != nil {
		panic(err)
	}

	resp, err := client.Grant(ctx, 1)
	if err != nil {
		panic(err)
	}

	var addr = &net.TCPAddr{
		IP:   GetLocalIP(),
		Port: int(port),
	}

	var key = fmt.Sprintf("flex/%s/%s", etcdConfig.App, primitive.NewObjectID().Hex())
	_, err = client.Put(ctx, key, string(addrs.Marshal(addr)), clientv3.WithLease(resp.ID))
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
		}
	}
}
