package p2p

import (
	"encoding/gob"
	"io"
)

// decode any message between peers over the network
type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type DefaultDecoder struct {
}

func (d DefaultDecoder) Decode(reader io.Reader, data *RPC) error {
	return gob.NewDecoder(reader).Decode(&data.Payload)
}
