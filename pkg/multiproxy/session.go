package multiproxy

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/evilmartians/redis-proxy/pkg/redis"
	"github.com/secmask/go-redisproto"
)

type sessionState interface {
	handle(s *Session, cmd *redis.Command) (error, sessionState)
}

type regularState struct{}

func (st regularState) handle(s *Session, cmd *redis.Command) (error, sessionState) {
	if cmd.Name == "SELECT" {
		err := s.writer.WriteError("ERR re-selecting database is not allowed")
		if err != nil {
			return err, nil
		}
		s.writer.Flush()
		return fmt.Errorf("SELECT is called after handshake"), nil
	}

	response, err := s.rdb.Execute(cmd)

	if err != nil {
		return err, nil
	}

	// TODO: implement pipeline
	// if cmd.Last {
	_, err = s.writer.Write([]byte(response + "\r\n"))

	if err != nil {
		return err, nil
	}

	err = s.writer.Flush()
	if err != nil {
		return err, nil
	}

	return nil, nil
}

type handshakeState struct{}

func (st handshakeState) handle(s *Session, cmd *redis.Command) (error, sessionState) {
	if cmd.Name != "SELECT" {
		err := s.writer.WriteError("ERR command is called before `select`")
		if err != nil {
			return err, nil
		}

		err = s.writer.Flush()
		if err != nil {
			return err, nil
		}

		return fmt.Errorf("Command called before SELECT: %s", cmd.Name), nil
	}

	db := cmd.Args[0].(string)
	client := s.proxy.LookupClient(db)

	if client == nil {
		return fmt.Errorf("Database not found: %s", db), nil
	}

	s.dbname = db
	s.rdb = client

	err := s.writer.WriteSimpleString("OK")

	if err != nil {
		return err, nil
	}

	if cmd.Last {
		s.writer.Flush()
	}

	return nil, regularState{}
}

// Session contains the reference to the target Redis pool,
// underlying IO object, and is responsible for reading and writing data
type Session struct {
	io     io.ReadWriter
	dbname string
	rdb    redis.RedisClient

	parser *redisproto.Parser
	writer *redisproto.Writer

	proxy *Proxy
	state sessionState
}

func NewSession(io io.ReadWriter, p *Proxy) *Session {
	parser := redisproto.NewParser(io)
	writer := redisproto.NewWriter(bufio.NewWriter(io))

	return &Session{io: io, parser: parser, writer: writer, proxy: p, state: handshakeState{}}
}

// HandleCommands continuosly reads and executes commands from IO
func (s *Session) HandleCommands() error {
	for {
		err := s.HandleCommand()

		if err != nil {
			return err
		}
	}
}

// HandleCommand reads a command from the client's IO, parses it,
// passes to the Redis client and writes the response back.
// It is also responsible to treat some commands differently:
//    - SELECT — used to route a client to the specified database pool
// (this is the only way to associate a client with a database)
//    - SCRIPT — scripts management
//    - CLIENT — client state commands (e.g., SETNAME)
//    - MULTI/EXEC — accumulate commands and wait for exec to process them together
// TODO: To be implemented
func (s *Session) HandleCommand() error {
	command, err := s.readCommand()

	if err != nil {
		return err
	}

	// Protocol-level error has been already handled
	if command == nil {
		return nil
	}

	// Handle commands which do not require calling a real Redis
	if handler, ok := fastlane[command.Name]; ok {
		return handler(s, command)
	}

	err, newState := s.state.handle(s, command)

	if err != nil {
		return err
	}

	if newState != nil {
		s.state = newState
	}

	return nil
}

func (s *Session) readCommand() (*redis.Command, error) {
	protoCmd, err := s.parser.ReadCommand()
	if err != nil {
		_, ok := err.(*redisproto.ProtocolError)
		if ok {
			err = s.writer.WriteError(err.Error())
			if err != nil {
				return nil, err
			}
			s.writer.Flush()
		} else {
			return nil, err
		}

		return nil, nil
	}

	restArgsCount := protoCmd.ArgCount() - 1

	command := redis.Command{
		Name: strings.ToUpper(string(protoCmd.Get(0))),
		Last: protoCmd.IsLast(),
		Args: make([]interface{}, restArgsCount),
	}

	for i := 1; i < protoCmd.ArgCount(); i++ {
		command.Args[i-1] = string(protoCmd.Get(i))
	}

	return &command, nil
}
