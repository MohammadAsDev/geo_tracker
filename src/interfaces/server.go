package interfaces

import (
	"github.com/MohammadAsDev/geo_tracker/src/entities/messages"
)

type TrackingSever interface {
	StartServer() error
	GetErrsChan() chan error
	GetMessagesChannel() chan messages.TrackingMessage
}
