package errors

import "fmt"

type Server struct {
	code    int
	message string
}

func NewServer(code int, message string) *Server {
	return &Server{
		code:    code,
		message: message,
	}
}

func (s *Server) Error() string {
	return fmt.Sprintf("HTTP server error with code: %d and message: %s", s.code, s.message)
}
