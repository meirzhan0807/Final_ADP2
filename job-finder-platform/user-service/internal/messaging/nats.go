package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type NATSClient struct{ conn *nats.Conn }

func NewNATSClient(url string) (*NATSClient, error) {
	conn, err := nats.Connect(url,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(20),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Printf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS reconnected to %s", nc.ConnectedUrl())
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}
	log.Printf("NATS connected to %s", url)
	return &NATSClient{conn: conn}, nil
}

func (c *NATSClient) Publish(subject string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("NATS marshal error: %v", err)
		return
	}
	if err := c.conn.Publish(subject, b); err != nil {
		log.Printf("NATS publish error: %v", err)
	}
}

func (c *NATSClient) Subscribe(subject string, fn func([]byte)) {
	c.conn.Subscribe(subject, func(msg *nats.Msg) { fn(msg.Data) })
}

func (c *NATSClient) Close() {
	if c.conn != nil {
		c.conn.Drain()
	}
}
