// The server package contains servers implemenations
package server

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tt := []struct {
		name               string
		addr               string
		expectedError      string
		expectedSocketType string
	}{
		{
			"With TCP address",
			"tcp://127.0.0.1:0",
			"",
			"tcp",
		},
		{
			"With Unix address",
			"unix:///tmp/server.sock",
			"",
			"unix",
		},
		{
			"With unsupported protocol",
			"ftp://127.0.0.1:4321",
			"unsupported protocol: ftp",
			"",
		},
		{
			"With implicit TCP",
			"127.0.0.1:0",
			"",
			"tcp",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := New(tc.addr, connectionHandler)

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
				return
			}

			assert.NoError(t, err)

			defer s.Shutdown()

			assert.Equal(t, tc.expectedSocketType, s.Type())
		})
	}
}

func TestRun(t *testing.T) {
	tt := []struct {
		name string
		addr string
		dial string
	}{
		{
			"TCP server",
			"tcp://127.0.0.1:0",
			"tcp",
		},
		{
			"Unix server",
			"unix:///tmp/test.sock",
			"unix",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := New(tc.addr, connectionHandler)
			assert.NoError(t, err)

			defer s.Shutdown()
			go s.Run()

			conn, err := net.Dial(tc.dial, s.Addr())
			if err != nil {
				t.Error("could not connect to server: ", err)
			}

			_, err = conn.Write([]byte("hello\n"))
			if err != nil {
				t.Error("could not write payload to server:", err)
			}

			out := make([]byte, 1024)
			_, err = conn.Read(out)
			if err == nil {
				if bytes.Equal(out, []byte("ECHO: hello\n")) {
					t.Error("response did match expected output")
				}
			} else {
				t.Error("could not read from connection")
			}

			_, err = conn.Write([]byte("STOP\n"))
			if err != nil {
				t.Error("could not write payload to server:", err)
			}

			_, err = conn.Read(out)
			assert.EqualError(t, err, io.EOF.Error())
		})
	}
}

func connectionHandler(conn net.Conn) {
	defer conn.Close()

	for {
		input, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			break
		}

		command := strings.TrimSpace(input)
		if command == "STOP" {
			break
		}

		result := "ECHO: " + command + "\n"
		conn.Write([]byte(result)) // nolint:errcheck
	}
}
