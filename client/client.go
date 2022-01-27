package client

import (
	"context"
	"net"
	"time"

	"github.com/yixinin/flex/logger"
	"github.com/yixinin/flex/message"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	etcdClient *clientv3.Client
	app        string
	connMgr    ConnManager
	event      chan ConnEvent
	cancel     func()
}

func (c *Client) addAddr(ctx context.Context, id, addr string) error {
	ip, port, err := parseIp(addr)
	if err != nil {
		logger.Errorf(ctx, "parse addr error:%v", err)
		return err
	}
	tcpAddr := &net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	c.event <- ConnEvent{
		EventType: EventAdd,
		Id:        id,
		Addr:      tcpAddr,
	}
	return nil
}

func (c *Client) delAddr(ctx context.Context, id string) error {
	c.event <- ConnEvent{
		EventType: EventAdd,
		Id:        id,
	}
	return nil
}

func NewClient(endpoints []string, appName string) *Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	var ctx, cancel = context.WithCancel(context.Background())
	c := &Client{
		etcdClient: client,
		app:        appName,
		event:      make(chan ConnEvent, 5),
		cancel:     cancel,
	}
	c.run(ctx)
	return c
}

func (c *Client) run(ctx context.Context) {
	go c.Watch(ctx)
	go c.onConnEvent(ctx)
}

func (c *Client) Close() {
	c.cancel()
	close(c.event)
}

func (c *Client) Publish(ctx context.Context, key, groupKey string, payload []byte) error {
	return c.connMgr.Send(ctx, key, groupKey, payload)
}

func (c *Client) Recv(ctx context.Context, timeout time.Duration) (message.Message, error) {
	return nil, nil
}
