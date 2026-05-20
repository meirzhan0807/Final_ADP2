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
	)
	if err != nil {
		return nil, fmt.Errorf("nats: %w", err)
	}
	log.Printf("NATS connected: %s", url)
	return &NATSClient{conn: conn}, nil
}

func (c *NATSClient) Publish(subject string, data interface{}) {
	b, _ := json.Marshal(data)
	c.conn.Publish(subject, b)
}

func (c *NATSClient) Close() { c.conn.Drain() }
