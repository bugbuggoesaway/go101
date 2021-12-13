package server

type spexServer struct{}

func (s *spexServer) Start() error {
	return nil
}

func (s *spexServer) Stop() error {
	return nil
}

func NewSPEXServer(srv interface{}, address string) Server {
	return &spexServer{}
}
