package interfaces

import "github.com/MohammadAsDev/geo_tracker/src/entities/events"

type ConsumingErrCode uint8

const (
	PARSING_ERR   ConsumingErrCode = 0
	CONSUMING_ERR ConsumingErrCode = 1
	TIMEOUT_ERR   ConsumingErrCode = 2
)

type ConsumingErr struct {
	Err  error
	Code ConsumingErrCode
}

func (err ConsumingErr) Error() string {
	return err.Err.Error()
}

type Consumer interface {
	Start() error
	GetCommandsChannel() chan events.Event
	GetErrorsChannel() chan ConsumingErr
}
