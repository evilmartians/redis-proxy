package multiproxy

import (
	"bufio"
	"io"

	"github.com/secmask/go-redisproto"
)

// Session contains the reference to the target Redis pool,
// underlying IO object, and is responsible for reading and writing data
type Session struct {
	io io.ReadWriter
}

// HandleCommands reads commands from the client's IO, parses them,
// passes to the Redis client and writes the response back.
// It is also responsible to treat some commands differently:
//    - SELECT — used to route a client to the specified database pool
// (this is the only way to associate a client with a database)
//    - SCRIPT — scripts management
//    - CLIENT — client state commands (e.g., SETNAME)
//    - MULTI/EXEC — accumulate commands and wait for exec to process them together
// TODO: To be implemented
func (s *Session) HandleCommands() error {
	parser := redisproto.NewParser(s.io)
	writer := redisproto.NewWriter(bufio.NewWriter(s.io))

	for {
		command, err := parser.ReadCommand()
		if err != nil {
			_, ok := err.(*redisproto.ProtocolError)
			if ok {
				err = writer.WriteError(err.Error())
				if err != nil {
					return err
				}
			} else {
				return err
			}
			continue
		}

		err = writer.WriteBulkString("OK")

		if command.IsLast() {
			writer.Flush()
		}

		if err != nil {
			return err
		}
	}
}
