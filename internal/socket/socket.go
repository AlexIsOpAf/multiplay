package socket

import (
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

type Soc struct {
	srv socketio.Server
}

func NewSocket() *Soc {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		s.Emit("Reply", "Communication Open")
		return nil
	})

	go server.Serve()

	http.Handle("/socket/", server)

	return &Soc{
		srv: *server,
	}
}

func (s *Soc) Close() error {
	err := s.srv.Close()

	if err != nil {
		return err
	}

	return nil
}
