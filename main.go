package main

import (
	"fmt"
	"runtime"
	"os"
	"os/signal"
	"time"
	"syscall"
	"path/filepath"
	bolt_api "github.com/boltdb/bolt"
	logger "github.com/0187773933/FireC2Server/v1/logger"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	types "github.com/0187773933/FireC2Server/v1/types"
	server "github.com/0187773933/FireC2Server/v1/server"
)

var s server.Server
var DB *bolt_api.DB

func SetupCloseHandler() {
	c := make( chan os.Signal )
	signal.Notify( c , os.Interrupt , syscall.SIGTERM , syscall.SIGINT )
	go func() {
		<-c
		// logger.Log.Println( "\r- Ctrl+C pressed in Terminal" )
		fmt.Println( "\r" )
		logger.Log.Println( "Ctrl+C pressed in Terminal" )
		logger.Log.Printf( "Shutting Down %s Server" , s.Config.ServerName )
		DB.Close()
		s.FiberApp.Shutdown()
		os.Exit( 0 )
	}()
}

func SetupDB( config *types.ConfigFile ) {
	db , _ := bolt_api.Open( config.BoltDBPath , 0600 , &bolt_api.Options{ Timeout: ( 3 * time.Second ) } )
	DB = db
	tx , err := DB.Begin( true )
	tx.CreateBucketIfNotExists( []byte( "state" ) );
	tx.CreateBucketIfNotExists( []byte( "logs" ) );
	tx.Commit();
	if err != nil { panic( err ) }
}

func SetupStackTraceReport() {
	if r := recover(); r != nil {
		stacktrace := make([]byte, 1024)
		runtime.Stack(stacktrace, true)
		fmt.Printf("%s\n", stacktrace)
	}
}

func main() {

	config_file_path := "./config.yaml"
	if len( os.Args ) > 1 { config_file_path , _ = filepath.Abs( os.Args[ 1 ] ) }
	config := utils.ParseConfig( config_file_path )

	SetupCloseHandler()
	SetupDB( &config )
	defer SetupStackTraceReport()
	// var log =
	// logger.Init()
	// utils.GenerateNewKeys()
	utils.WriteLoginURLPrefix( config.ServerLoginUrlPrefix )
	s = server.New( DB , config )
	logger.Log.Printf( "Loaded Config File From : %s" , config_file_path )
	s.Start()

}