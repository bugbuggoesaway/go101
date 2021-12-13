package server

type Protocol uint8

const (
	GRPC Protocol = 1 << 0
	SPEX Protocol = 1 << 1
)

func (p Protocol) IsGRPCEnabled() bool {
	return p&GRPC == GRPC
}

func (p Protocol) IsSPEXEnabled() bool {
	return p&SPEX == SPEX
}
