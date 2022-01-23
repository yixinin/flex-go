package client

import (
	"context"
	"net"
	"time"

	"github.com/yixinin/flex/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	etcdClient *clientv3.Client
	app        string
	connMgr    ConnManager
	event      chan Event
}

func (c *Client) AddAddr(ctx context.Context, id, addr string) error {
	ip, port, err := parseIp(addr)
	if err != nil {
		logger.Errorf(ctx, "parse addr error:%v", err)
		return err
	}
	tcpAddr := &net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	c.event <- Event{
		EventType: EventAdd,
		Id:        id,
		Addr:      tcpAddr,
	}
	return nil
}

func (c *Client) DelAddr(ctx context.Context, id string) error {
	c.event <- Event{
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
	return &Client{
		etcdClient: client,
		app:        appName,
		event:      make(chan Event, 5),
	}
}

func (c *Client) Run(ctx context.Context) {
	go c.Watch(ctx)
	go c.recvEvent(ctx)
}

func (c *Client) Close() {
	close(c.event)
}
