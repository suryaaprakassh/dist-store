package main

// Payload represents anything
// that is transferred between peers in a network
type Payload struct {
	Key  string
	Data []byte
}
