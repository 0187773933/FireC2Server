package server

import (
	"fmt"
	"time"
	"sync"
	// context "context"
	logger "github.com/0187773933/FireC2Server/v1/logger"
	fiber "github.com/gofiber/fiber/v2"
	fiber_cookie "github.com/gofiber/fiber/v2/middleware/encryptcookie"
	fiber_cors "github.com/gofiber/fiber/v2/middleware/cors"
	favicon "github.com/gofiber/fiber/v2/middleware/favicon"
	types "github.com/0187773933/FireC2Server/v1/types"
	redis "github.com/redis/go-redis/v9"
	// adb_wrapper "ADBWrapper/v1/wrapper"
	adb_wrapper "github.com/0187773933/ADBWrapper/v1/wrapper"
	// tv "github.com/0187773933/FireC2Server/v1/tv"
	tv_controller "github.com/0187773933/TVController/v1/controller"
	utils "github.com/0187773933/FireC2Server/v1/utils"
)

var GlobalServer *Server
var log = logger.GetLogger()
// var LocalIPS = utils.GetLocalIPAddresses()

const ACTIVITY_PROFILE_PICKER = "com.amazon.ftv.profilepicker/com.amazon.ftv.profilepicker.ui.PickerActivity"
const ADS_PACKAGE = "com.amazon.wallpaper.LockscreenWallpaperService"

type Status struct {
	StartTime string `json:"start_time"`
	StartTimeOBJ time.Time `json:"-"`
	PlayerName string `json:"player_name"`
	PlayerCommand string `json:"player_command"`
	NowPlayingID string `json:"now_playing_id"`
	NowPlayingText string `json:"now_playing_text"`
	PreviousPlayerName string `json:"previous_player_name"`
	PreviousPlayerCommand string `json:"previous_player_command"`
	PreviousStartTime string `json:"previous_start_time"`
	PreviousStartTimeOBJ time.Time `json:"-"`
	PreviousStartTimeDuration time.Duration `json:"-"`
	PreviousStartTimeDurationSeconds float64 `json:"previous_start_time_duration_seconds"`
	ADB adb_wrapper.Status `json:"adb"`
	TV tv_controller.Status `json:"tv"`
}

type Server struct {
	FiberApp *fiber.App `yaml:"fiber_app"`
	Config types.ConfigFile `yaml:"config"`
	DB *redis.Client `yaml:"-"`
	ADB adb_wrapper.Wrapper `json:"-"`
	TV *tv_controller.Controller `json:"-"`
	Status Status `json:"-"`
	StatusMutex sync.Mutex `json:"-"`
}

func ( s *Server ) SetupRoutes() {
	s.SetupPublicRoutes()
	s.SetupAdminRoutes()
}

func ( s *Server ) Start() {
	log.Printf( "Listening on http://localhost:%s" , s.Config.ServerPort )
	// log.Printf( "Listening on http://%s:%s" , s.Config.ServerPort )
	fmt.Printf( "Admin Login @ http://localhost:%s/%s\n" , s.Config.ServerPort , s.Config.ServerLoginUrlPrefix )
	fmt.Printf( "Admin Username === %s\n" , s.Config.AdminUsername )
	fmt.Printf( "Admin Password === %s\n" , s.Config.AdminPassword )
	fmt.Printf( "Admin API Key === %s\n" , s.Config.ServerAPIKey )
	// go s.Governor()
	// go s.ADB.WatchLog()
	s.FiberApp.Listen( fmt.Sprintf( ":%s" , s.Config.ServerPort ) )
}

func New( db *redis.Client , config types.ConfigFile ) ( server Server ) {
	server.FiberApp = fiber.New()
	server.Config = config
	server.DB = db
	server.ADB = server.ADBConnect()
	utils.PrettyPrint( server.ADB )
	server.TV = tv_controller.New( &config.TV )
	server.StatusMutex = sync.Mutex{}
	GlobalServer = &server
	server.StoreLibrary()
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