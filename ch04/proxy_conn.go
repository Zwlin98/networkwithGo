package ch04

import (
	"io"
	"net"
)

func proxyConn(source, destination string) error {
	connSource, err := net.Dial("tcp", source)
	if err != nil {
		return err
	}
	defer func() { _ = connSource.Close() }()

	connDestination, err := net.Dial("tcp", destination)
	if err != nil {
		return err
	}
	defer func() { _ = connDestination.Close() }()

	go func() {
		_, _ = io.Copy(connSource, connDestination)
	}()

	_, err = io.Copy(connDestination, connSource)

	return err
}
