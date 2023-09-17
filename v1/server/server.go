package server

import (
	"fmt"
	"time"
	fiber "github.com/gofiber/fiber/v2"
	fiber_cookie "github.com/gofiber/fiber/v2/middleware/encryptcookie"
	fiber_cors "github.com/gofiber/fiber/v2/middleware/cors"
	favicon "github.com/gofiber/fiber/v2/middleware/favicon"
	types "github.com/0187773933/FireC2Server/v1/types"
	bolt_api "github.com/boltdb/bolt"
	media_player "github.com/0187773933/FireC2Server/v1/media_player"
)

var GlobalServer *Server

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
	s.Printf( "Listening on http://localhost:%s\n" , s.Config.ServerPort )
	s.Printf( "Admin Login @ http://localhost:%s/%s\n" , s.Config.ServerPort , s.Config.ServerLoginUrlPrefix )
	s.Printf( "Admin Username === %s\n" , s.Config.AdminUsername )
	s.Printf( "Admin Password === %s\n" , s.Config.AdminPassword )
	s.Printf( "Admin API Key === %s\n" , s.Config.ServerAPIKey )
	s.FiberApp.Listen( fmt.Sprintf( ":%s" , s.Config.ServerPort ) )
}

func New( config types.ConfigFile ) ( server Server ) {
	server.FiberApp = fiber.New()
	server.Config = config
	GlobalServer = &server
	db , _ := bolt_api.Open( config.BoltDBPath , 0600 , &bolt_api.Options{ Timeout: ( 3 * time.Second ) } )
	server.DB = db
	// tx , err := server.DB.Begin( true )
	// tx.CreateBucketIfNotExists( []byte( "state" ) );
	// tx.Commit();
	// fmt.Println( "err ===" , err )
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