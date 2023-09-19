package server

import (
	"fmt"
	// "time"
	logrus "github.com/sirupsen/logrus"
	logger "github.com/0187773933/FireC2Server/v1/logger"
	fiber "github.com/gofiber/fiber/v2"
	fiber_cookie "github.com/gofiber/fiber/v2/middleware/encryptcookie"
	fiber_cors "github.com/gofiber/fiber/v2/middleware/cors"
	favicon "github.com/gofiber/fiber/v2/middleware/favicon"
	types "github.com/0187773933/FireC2Server/v1/types"
	bolt_api "github.com/boltdb/bolt"
	media_player "github.com/0187773933/FireC2Server/v1/media_player"
)

var GlobalServer *Server
var log *logrus.Logger

type Server struct {
	FiberApp *fiber.App `yaml:"fiber_app"`
	Config types.ConfigFile `yaml:"config"`
	DB *bolt_api.DB `yaml:"-"`
	MediaPlayer *media_player.MediaPlayer `yaml:"-"`
}

func ( s *Server ) SetupRoutes() {
	s.SetupPublicRoutes()
	s.SetupAdminRoutes()
}

func ( s *Server ) Start() {
	log.Printf( "Listening on http://localhost:%s" , s.Config.ServerPort )
	fmt.Printf( "Admin Login @ http://localhost:%s/%s\n" , s.Config.ServerPort , s.Config.ServerLoginUrlPrefix )
	fmt.Printf( "Admin Username === %s\n" , s.Config.AdminUsername )
	fmt.Printf( "Admin Password === %s\n" , s.Config.AdminPassword )
	fmt.Printf( "Admin API Key === %s\n" , s.Config.ServerAPIKey )
	s.FiberApp.Listen( fmt.Sprintf( ":%s" , s.Config.ServerPort ) )
}

func New( db *bolt_api.DB , config types.ConfigFile ) ( server Server ) {
	server.FiberApp = fiber.New()
	server.Config = config
	server.DB = db
	GlobalServer = &server

	log = logger.Log
	log.Debug( "Server Starting" )
	server.MediaPlayer = media_player.New( db , &config )
	server.FiberApp.Use( server.LogRequest )
	server.FiberApp.Use( favicon.New() )
	server.FiberApp.Use( fiber_cookie.New( fiber_cookie.Config{
		Key: server.Config.ServerCookieSecret ,
	}))
	server.FiberApp.Use( fiber_cors.New( fiber_cors.Config{
		AllowOrigins: fmt.Sprintf( "%s, %s" , server.Config.ServerBaseUrl , server.Config.ServerLiveUrl ) ,
		AllowHeaders:  "Origin, Content-Type, Accept, key",
	}))
	server.SetupRoutes()
	server.FiberApp.Get( "/*" , func( context *fiber.Ctx ) ( error ) { return context.Redirect( "/" ) } )
	return
}