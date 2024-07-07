package p2p

import (
	"errors"
	"io"
	"log/slog"
)

// decode any message between peers over the network
type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type DefaultDecoder struct {
}

func (d DefaultDecoder) Decode(reader io.Reader, data *RPC) error {

	buf := make([]byte, 1)

	_, err := reader.Read(buf)

	if err != nil {
		return errors.New("Stream Mode Not Found")
	}

	if buf[0] == StreamingMode {
		data.Stream = true
		slog.Debug("Started streaming mode")
		return nil
	}

	buf = make([]byte, 1028)

	n, err := reader.Read(buf)

	if err != nil {
		slog.Warn("Error Read Value from connection ", err)
		return err
	}

	data.Payload = buf[:n]

	return nil
}
