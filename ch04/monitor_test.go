package ch04

import (
	"io"
	"log"
	"net"
	"os"
)

type Monitor struct {
	*log.Logger
}

func (m *Monitor) Write(p []byte) (int, error) {
	err := m.Output(2, string(p))
	if err != nil {
		log.Println(err)
	}
	return len(p), nil
}

func ExampleMonitor() {
	monitor := &Monitor{Logger: log.New(os.Stdout, "monitor: ", 0)}
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		monitor.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)

		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer func() { _ = conn.Close() }()

		b := make([]byte, 1024)
		r := io.TeeReader(conn, monitor)
		w := io.MultiWriter(conn, monitor)

		for {
			n, err := r.Read(b)

			if err != nil {
				if err != io.EOF {
					monitor.Println(err)
					return
				}
				break
			}

			_, err = w.Write(b[:n])
			if err != nil {
				if err != io.EOF {
					monitor.Println(err)
					return
				}
				break
			}
		}

	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		monitor.Fatal(err)
	}

	_, err = conn.Write([]byte("Test\n"))
	if err != nil {
		monitor.Fatal(err)
	}

	_ = conn.Close()
	<-done

	//Output:
	//monitor: Test
	//monitor: Test
}
