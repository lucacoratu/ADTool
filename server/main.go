package main

import "github.com/lucacoratu/ADTool/server/server"

func main() {
	server := server.APIServer{}
	err := server.Init()
	if err != nil {
		return
	}
	server.Run()
}
