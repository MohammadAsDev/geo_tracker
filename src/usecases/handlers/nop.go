package handlers

type NopCommand struct {
	_Handler
}

func (handler NopCommand) Hanlde() error {
	// do nothing
	return nil
}
