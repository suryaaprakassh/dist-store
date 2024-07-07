package main

type Action int

const (
	Save Action = iota
	Delete
)

// is a type that is transferred between
// nodes in a tcp network
type Payload struct {
	Action Action
	Key    string
	Size   int64
}
