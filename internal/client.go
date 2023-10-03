package internal

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/xdg-go/scram"
	"net"
	"scram_tcp/utils"
	"time"
)

type Client struct {
	*scram.Client
	address     string
	ctx         context.Context
	logger      *log.Logger
	authTimeout time.Duration
}

func NewClient(address, username, password string) (*Client, error) {
	scramClt, err := scram.SHA256.NewClient(username, password, "")
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:      scramClt,
		address:     address,
		ctx:         context.Background(),
		logger:      log.New(),
		authTimeout: time.Second * AUTH_TIMEOUT,
	}, nil
}

func (c *Client) MakeRequest() {
	ctx, cancelFunc := context.WithTimeout(c.ctx, c.authTimeout)
	defer cancelFunc()

	conn, err := c.configureConn()
	if err != nil {
		c.logger.Error("Dial failed:", err.Error())
		return
	}
	defer conn.Close()

	err = c.authorize(ctx, conn)
	if err != nil {
		c.logger.Error("SCRAM auth err:", err.Error())
		return
	}

	quoteBytes, err := utils.ReadConnWithCtx(ctx, conn)
	if err != nil {
		c.logger.Error("read conn err:", err)
		return
	}

	c.logger.Println(string(quoteBytes))
}

func (c *Client) configureConn() (net.Conn, error) {
	tcpServer, err := net.ResolveTCPAddr(TYPE, c.address)

	if err != nil {
		c.logger.Error("ResolveTCPAddr failed:", err.Error())
		return nil, err
	}

	return net.DialTCP(TYPE, nil, tcpServer)
}

func (c *Client) authorize(ctx context.Context, conn net.Conn) error {
	conv := c.NewConversation()
	got, err := conv.Step("")
	if err != nil {
		return err
	}
	err = utils.WriteConn(conn, []byte(got))
	if err != nil {
		return err
	}

	for {
		if conv.Done() {
			break
		}

		var received []byte
		received, err = utils.ReadConnWithCtx(ctx, conn)
		if err != nil {
			break
		}

		got, err = conv.Step(string(received))
		if err != nil {
			break
		}

		err = utils.WriteConn(conn, []byte(got))
		if err != nil {
			break
		}
	}

	if err != nil {
		return err
	}

	if !conv.Valid() {
		return errors.New("conv is not valid")
	}

	if !conv.Done() {
		return errors.New("conv is not finished")
	}
	return nil
}
