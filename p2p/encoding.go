package p2p

import (
	"io"
)

// decode any message between peers over the network
type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type DefaultDecoder struct {
}

func (d DefaultDecoder) Decode(reader io.Reader, data *RPC) error {
	buf := make([]byte, 2048)
	n, err := reader.Read(buf)
	if err != nil {
		return err
	}
	data.Payload = buf[:n]
	return nil
}
