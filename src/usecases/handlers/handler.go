package handlers

import (
	"github.com/MohammadAsDev/geo_tracker/src/entities/events"
	"github.com/MohammadAsDev/geo_tracker/src/interfaces"
	"github.com/MohammadAsDev/geo_tracker/src/usecases/tracker"
)

type _Handler struct {
	Event   events.Event
	Tracker tracker.Tracker
}

func BuildHandler(event events.Event, system tracker.Tracker) (handler interfaces.Handler) {
	switch event.EventId {
	case events.DEACTIVATE_EVENT:
		handler = DeactivateCommand{
			_Handler: _Handler{
				Tracker: system,
				Event:   event,
			},
		}
	case events.START_EVNET:
		handler = StartCommand{
			_Handler: _Handler{
				Tracker: system,
				Event:   event,
			},
		}
	default:
		handler = NopCommand{}
	}
	return
}
