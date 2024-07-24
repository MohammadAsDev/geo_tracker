package interfaces

import "github.com/MohammadAsDev/geo_tracker/src/entities/messages"

type SystemCommandId uint8

type SystemCommand interface {
	GetSystemCommandId() SystemCommandId
}

type ConnServer interface {
	GetMessagesChannel() chan messages.TrackingMessage
	GetErrsChan() chan error
	StartServer() error
	GetCommandsChannel() chan SystemCommand
}
