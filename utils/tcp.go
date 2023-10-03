package utils

import (
	"bufio"
	"context"
	"errors"
	"net"
)

func WriteConn(conn net.Conn, b []byte) error {
	_, err := conn.Write(append(b, '\n'))
	return err
}

func ReadConn(conn net.Conn, ch chan []byte, errChan chan error) {
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		ch <- scanner.Bytes()
	} else {
		errChan <- errors.New("scan false")
	}
	return
}

func ReadConnWithCtx(ctx context.Context, conn net.Conn) ([]byte, error) {
	dataChan := make(chan []byte)
	errChan := make(chan error)
	defer close(dataChan)
	defer close(errChan)
	//go ReadConn(conn, size, dataChan, errChan)
	go ReadConn(conn, dataChan, errChan)
	for {
		select {
		case buf := <-dataChan:
			return buf, nil
		case err := <-errChan:
			return nil, err
		case <-ctx.Done():
			return nil, errors.New("context done")
		}
	}
}
