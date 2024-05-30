package events

type Consumer interface {
	Start() error
}
