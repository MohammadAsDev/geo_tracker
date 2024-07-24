package ws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MohammadAsDev/geo_tracker/src/config"
	"github.com/MohammadAsDev/geo_tracker/src/entities/messages"
	"github.com/MohammadAsDev/geo_tracker/src/interfaces"
	"github.com/gorilla/websocket"
)

type CheckOriginMethod func(r *http.Request) bool

const _CHAN_SIZE = 1024

const (
	DEACTIVATE interfaces.SystemCommandId = 0
)

type WsCommand struct {
	CommandId interfaces.SystemCommandId
	DriverId  uint64
}

func (command WsCommand) GetSystemCommandId() interfaces.SystemCommandId {
	return command.CommandId
}

type WebSocketServer struct {
	_Ctx context.Context

	_Addr     string
	_Clients  map[uint64]*websocket.Conn
	_Upgrader *websocket.Upgrader

	_CommandsChan chan interfaces.SystemCommand
	_MessagesChan chan messages.TrackingMessage
	_ErrsChan     chan error

	_SystemCache interfaces.Cache
}

type WsBuilder struct {
	ws *WebSocketServer
}

func NewWsBuilder(ctx context.Context, config *config.WsConfig) *WsBuilder {
	builder := &WsBuilder{}
	builder.ws = &WebSocketServer{
		_Ctx: ctx,

		_Addr:    config.Addr,
		_Clients: map[uint64]*websocket.Conn{},
		_Upgrader: &websocket.Upgrader{
			ReadBufferSize:  config.ReadBufferSize,
			WriteBufferSize: config.WriteBufferSize,
		},

		_CommandsChan: make(chan interfaces.SystemCommand),
		_MessagesChan: make(chan messages.TrackingMessage, _CHAN_SIZE),
		_ErrsChan:     make(chan error),
	}
	return builder
}

func (builder *WsBuilder) WithSystemCache(cache interfaces.Cache) *WsBuilder {
	builder.ws._SystemCache = cache
	return builder
}

func (builder *WsBuilder) WithDefaultCheckOrigin() *WsBuilder {
	builder.ws._Upgrader.CheckOrigin = builder.ws._DefaultUpgradingMethod
	return builder
}

func (builder *WsBuilder) WithNopCheckOrigin() *WsBuilder {
	builder.ws._Upgrader.CheckOrigin = builder.ws._NopUpgradingMethod
	return builder
}

func (builder *WsBuilder) Build() interfaces.ConnServer {
	return builder.ws
}

func (server *WebSocketServer) StartServer() error {
	fmt.Printf("The server is running [%s]...\n", server._Addr)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Tracking server is running..."))
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/track", func(w http.ResponseWriter, r *http.Request) {
		conn, err := server._Upgrader.Upgrade(w, r, nil)
		if err != nil {
			server._ErrsChan <- err
			return
		}
		defer conn.Close()

		var driver_id uint64
		queries := r.URL.Query()

		if id, err := strconv.ParseUint(queries.Get("driver_id"), 0, 64); err == nil {
			driver_id = uint64(id)
		} else {
			server._ErrsChan <- err
			return
		}

		server._Clients[driver_id] = conn

		for {
			var message messages.TrackingMessage
			if err := conn.ReadJSON(&message); err != nil {
				server._ErrsChan <- errors.New("testing: " + err.Error())
				break
			}
			message.DriverId = driver_id
			server._MessagesChan <- message
		}
	})

	go func() {
		for {
			cmd := <-server._CommandsChan
			err := server._ExecuteCommand(cmd.(WsCommand))
			if err != nil {
				server._ErrsChan <- err
			}
		}
	}()

	err := http.ListenAndServe(server._Addr, nil)

	return err
}

func (server *WebSocketServer) _ExecuteCommand(cmd WsCommand) error {
	switch cmd.GetSystemCommandId() {
	case DEACTIVATE:
		driver_id := cmd.DriverId

		if _, found := server._Clients[driver_id]; !found {
			return errors.New("ws: can't find ws for user")
		}
		if err := server._Clients[driver_id].Close(); err != nil {
			return err
		}
		delete(server._Clients, driver_id)
	}
	return nil
}

func (server *WebSocketServer) GetMessagesChannel() chan messages.TrackingMessage {
	return server._MessagesChan
}

func (server *WebSocketServer) GetErrsChan() chan error {
	return server._ErrsChan
}

func (server *WebSocketServer) GetCommandsChannel() chan interfaces.SystemCommand {
	return server._CommandsChan
}

func (ws *WebSocketServer) _DefaultUpgradingMethod(r *http.Request) bool {
	queries := r.URL.Query()
	if ws._SystemCache == nil {
		panic("cache is not supported for ws server")
	}
	driver_id := queries.Get("driver_id")
	token, err := ws._SystemCache.GetDriverToken(driver_id)
	if err != nil {
		log.Print(err)
		return false
	}
	request_token := queries.Get("token")
	return request_token == token
}

func (ws *WebSocketServer) _NopUpgradingMethod(r *http.Request) bool {
	return true
}
