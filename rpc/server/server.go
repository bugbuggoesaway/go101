package server

import "log"

type Server interface {
	Start() error
	Stop() error
}

type Servers []Server

func (s Servers) Start() error {
	ch := make(chan error)
	for _, server := range s {
		go func(server Server) {
			err := server.Start()
			if err != nil {
				log.Printf("Start. err=[%v]", err)
				ch <- err
			}
		}(server)
	}
	return <-ch
}

func (s Servers) Stop() error {
	var stopErr error
	for _, server := range s {
		err := server.Stop()
		if err != nil {
			stopErr = err
		}
	}
	return stopErr
}

func NewServer(srv interface{}, address string, protocol Protocol) Server { //FIXME servers listen to the same address?
	var servers Servers
	if protocol.IsGRPCEnabled() {
		servers = append(servers, NewGRPCServer(srv, address))
	}
	if protocol.IsSPEXEnabled() {
		servers = append(servers, NewSPEXServer(srv, address))
	}
	return servers
}
