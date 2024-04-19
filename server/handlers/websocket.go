package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lucacoratu/ADTool/server/logging"
	"github.com/lucacoratu/ADTool/server/websocket"
)

type WebsocketHandler struct {
	logger logging.ILogger
}

func NewWebsocketHandler(logger logging.ILogger) *WebsocketHandler {
	return &WebsocketHandler{logger: logger}
}

/*
 * This function will handle when a client connects to the websocket endpoint
 */
func (wsh *WebsocketHandler) ServeAgentWs(pool *websocket.Pool, rw http.ResponseWriter, r *http.Request) {
	//Get the agent UUID from the mux variables
	vars := mux.Vars(r)
	agent_id, _ := strconv.Atoi(vars["id"])

	//Upgrade the connection to a Websocket connection
	ws, err := websocket.Upgrade(rw, r)
	//Check if an error occured
	if err != nil {
		//Log the error
		wsh.logger.Error(err.Error())
		return
	}

	//Create the client structure which will be saved in the pool
	client := &websocket.AgentClient{
		Conn:   ws,
		Pool:   pool,
		Status: "Offline",
		Id:     int64(agent_id),
	}

	//Call the client register function
	pool.RegisterAgent <- client
	//Start reading data from the connection
	//go client.Write()
	go client.Read()
}
