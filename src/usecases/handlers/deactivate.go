package handlers

import (
	"encoding/json"
	"errors"

	"github.com/MohammadAsDev/geo_tracker/src/entities/events"
	"github.com/MohammadAsDev/geo_tracker/src/infrastructure/ws"
)

type DeactivateCommand struct {
	_Handler
}

func (handler DeactivateCommand) Hanlde() error {
	payload := handler.Event.Payload
	var deactivate_payload events.DeactivatePayload

	if err := json.Unmarshal([]byte(payload), &deactivate_payload); err != nil {
		return errors.New("error while unmarshaling the deactivate event")
	}

	ws_command := ws.WsCommand{
		CommandId: ws.DEACTIVATE,
		DriverId:  deactivate_payload.DriverId,
	}

	cmd_chan := handler.Tracker.ConnServer().GetCommandsChannel()
	cmd_chan <- ws_command

	return nil
}
