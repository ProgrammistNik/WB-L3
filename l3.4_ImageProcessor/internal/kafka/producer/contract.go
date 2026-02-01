package producer

import "context"

type ServiceProducer interface {
	Send(context.Context, []byte, []byte) error
}
