package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTcpTransport(t *testing.T) {
	config := TcpTransportConfig{
		ListenAddr: ":8080",	
		ShakeHands: DefaultHandShake,
		Decoder: DefaultDecoder{},
	}

	tr := NewTcpTransport(config)
	assert.Equal(t, tr.ListenAddr, config.ListenAddr)

	assert.Nil(t, tr.ListenAndAccept())
}
