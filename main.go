package main

import (
	"fmt"
	"os"
	"os/signal"
	"context"
	"syscall"
	"path/filepath"
	redis "github.com/redis/go-redis/v9"
	logger "github.com/0187773933/FireC2Server/v1/logger"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	types "github.com/0187773933/FireC2Server/v1/types"
	server "github.com/0187773933/FireC2Server/v1/server"
)

var s server.Server
var DB *redis.Client

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

// https://pkg.go.dev/github.com/redis/go-redis/v9#Options
func SetupDB( config *types.ConfigFile ) {
	DB = redis.NewClient( &redis.Options{
		Addr: config.RedisAddress ,
		Password: config.RedisPassword ,
		DB: config.RedisDBNumber ,
	})
	var ctx = context.Background()
	ping_result , err := DB.Ping( ctx ).Result()
	logger.Log.Printf( "DB Connected : PING = %s" , ping_result )
	if err != nil { panic( err ) }
}

func main() {
	defer utils.SetupStackTraceReport()
	var config_file_path string
	if len( os.Args ) > 1 {
		config_file_path , _ = filepath.Abs( os.Args[ 1 ] )
	} else {
		config_file_path , _ = filepath.Abs( "./config.yaml" )
		if _ , err := os.Stat( config_file_path ); os.IsNotExist( err ) {
			config_file_path , _ = filepath.Abs( "./SAVE_FILES/config.yaml" )
			if _ , err := os.Stat( config_file_path ); os.IsNotExist( err ) {
				panic( "Config File Not Found" )
			}
		}
	}
	config := utils.ParseConfig( config_file_path )
	logger.Log.Printf( "Loaded Config File From : %s" , config_file_path )
	utils.WriteLoginURLPrefix( config.ServerLoginUrlPrefix )
	utils.FingerPrint( &config )

	SetupCloseHandler()
	SetupDB( &config )
	// utils.GenerateNewKeys()
	logger.Log.Printf( "Loading Server" )
	s = server.New( DB , config )
	logger.Log.Printf( "Starting Server" )
	s.Start()
}