package internal

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"github.com/kataras/iris/v12"
	"todo/internal/controller"
	"todo/internal/repository"
	"todo/pkg/db"
)

type Server struct {
	config     *Config
	httpServer *iris.Application
	wsServer   *socketio.Server
	mysql      *db.MySQL
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,
	}
}

func (s *Server) Start() error {
	s.mysql = db.NewMySQL(s.config.DbDsn, s.config.DbMaxOpenConns, s.config.DbMaxIdleConns)

	s.setupWsServer()
	s.setupHttpServer()

	go s.wsServer.Serve()
	defer s.wsServer.Close()

	return s.httpServer.Listen(
		fmt.Sprintf("%s:%d", s.config.ServerHost, s.config.ServerPort),
		iris.WithoutPathCorrection,
	)
}

func (s *Server) setupWsServer() {
	s.wsServer = socketio.NewServer(nil)

	// https://github.com/googollee/go-socket.io/blob/master/_examples/iris/main.go
	s.wsServer.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	//wsServer.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
	//	fmt.Println("notice:", msg)
	//	s.Emit("reply", "have "+msg)
	//})
	//wsServer.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
	//	s.SetContext(msg)
	//	return "recv " + msg
	//})
	//wsServer.OnEvent("/", "bye", func(s socketio.Conn) string {
	//	last := s.Context().(string)
	//	s.Emit("bye", last)
	//	s.Close()
	//	return last
	//})
	//wsServer.OnError("/", func(s socketio.Conn, e error) {
	//	fmt.Println("meet error:", e)
	//})
	s.wsServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
}

func (s *Server) setupHttpServer() {
	s.httpServer = iris.New()
	s.httpServer.HandleMany("GET POST", "/socket.io/{any:path}", iris.FromStd(s.wsServer))
	s.httpServer.Favicon("./web/favicon.ico")
	s.httpServer.HandleDir("/", iris.Dir("./web"))

	// http://localhost:8080/api
	api := s.httpServer.Party("/api")
	{
		api.Use(iris.Compression)
	}

	userRepository := repository.NewTicketRepository(s.mysql)
	userController := controller.NewTicketController(s.wsServer, userRepository)

	userController.SetupRoutes(api)
}
