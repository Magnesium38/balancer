package Balancer

type Server struct {
}

func (node *Server) Register() error {
	return nil
}

func (node *Server) GetHost() string {
	return nil
}

func (node *Server) GetPort() int {
	return nil
}

func (node *Server) ListenAndServe() error {
	return nil
}
