package handlers

import "net/http"

//Endpoint to check if the server is up and running
func Healthcheck(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("{\"status\":\"alive\"}"))
}
