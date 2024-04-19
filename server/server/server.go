package server

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/lucacoratu/ADTool/server/configuration"
	"github.com/lucacoratu/ADTool/server/database"
	"github.com/lucacoratu/ADTool/server/handlers"
	"github.com/lucacoratu/ADTool/server/logging"
	"github.com/lucacoratu/ADTool/server/websocket"
)

type APIServer struct {
	srv           *http.Server
	logger        logging.ILogger
	configFile    string
	configuration configuration.Configuration
	dbConnection  database.IConnection
}

func (api *APIServer) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.logger.Info(r.Method, "-", r.URL.Path, r.RemoteAddr)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		// compare the return-value to the authMW
		next.ServeHTTP(w, r)
	})
}

// Initialize the api http server based on the configuration file
func (api *APIServer) Init() error {
	//Initialize the logger
	api.logger = logging.NewDefaultDebugLogger()
	api.logger.Debug("Logger initialized")

	//Define command line arguments of the agent
	flag.StringVar(&api.configFile, "config", "", "The path to the configuration file")
	//Parse command line arguments
	flag.Parse()

	//Load the configuration from file
	err := api.configuration.LoadConfigurationFromFile(api.configFile)
	if err != nil {
		api.logger.Fatal("Error occured when loading the config from file,", err.Error())
		return err
	}
	api.logger.Debug("Loaded configuration from file")

	//Initialize the database connection
	api.dbConnection = database.NewMysqlConnection(api.logger, api.configuration)
	err = api.dbConnection.Init()
	if err != nil {
		api.logger.Error("Error occured when initializing database connection", err.Error())
		return err
	}
	api.logger.Debug("Connection to the database has been initialized")

	//Create the pool
	pool := websocket.NewPool(api.logger, api.dbConnection)
	//Start the pool in a goroutine
	go pool.Start()

	//Create the router
	r := mux.NewRouter()
	//Use the logging middleware
	r.Use(api.LoggingMiddleware)

	//Create the handlers
	wsHandler := handlers.NewWebsocketHandler(api.logger)
	agentHandler := handlers.NewAgentsHandler(api.logger, api.configuration, api.dbConnection, pool)

	//Add the routes
	//Create the subrouter for the API path
	apiGetSubrouter := r.PathPrefix("/api/v1/").Methods("GET").Subrouter()
	apiPostSubrouter := r.PathPrefix("/api/v1/").Methods("POST").Subrouter()
	//apiDeleteSubrouter := r.PathPrefix("/api/v1/").Methods("DELETE").Subrouter()
	//apiPutSubrouter := r.PathPrefix("/api/v1/").Methods("PUT").Subrouter()

	//Create the route for healthcheck
	apiGetSubrouter.HandleFunc("/healthcheck", handlers.Healthcheck)
	//Create the route for agents
	apiGetSubrouter.HandleFunc("/agents", agentHandler.GetAgents)
	//Create the route to get commands of the agent
	apiGetSubrouter.HandleFunc("/agents/{id:[0-9]+}/cmd", agentHandler.GetCommands)

	//Create the route for registering an agent
	apiPostSubrouter.HandleFunc("/agents", agentHandler.CreateAgent)
	//Create the route to execute a command on an agent
	apiPostSubrouter.HandleFunc("/agents/{id:[0-9]+}/cmd", agentHandler.ExecuteCommandOnAgent)
	//Create the route to execute a recurring command on an agent
	apiPostSubrouter.HandleFunc("/agents/{id:[0-9]+}/reccmd", agentHandler.ExecuteRecurringCommandOnAgent)

	//Create the route which will handle websocket agent connections
	apiGetSubrouter.HandleFunc("/agents/{id:[0-9]+}/ws", func(rw http.ResponseWriter, r *http.Request) {
		wsHandler.ServeAgentWs(pool, rw, r)
	})

	api.srv = &http.Server{
		Addr: api.configuration.ListeningAddress + ":" + strconv.Itoa(api.configuration.ListeningPort),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	return nil
}

// Start the api server
func (api *APIServer) Run() {
	var wait time.Duration = 5
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := api.srv.ListenAndServe(); err != nil {
			api.logger.Error(err.Error())
		}
	}()

	api.logger.Info("Started server on port", api.configuration.ListeningPort)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	api.srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	api.logger.Info("shutting down")
	os.Exit(0)
}
