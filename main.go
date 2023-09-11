package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"path/filepath"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	server "github.com/0187773933/FireC2Server/v1/server"
)

var s server.Server

func SetupCloseHandler() {
	c := make( chan os.Signal )
	signal.Notify( c , os.Interrupt , syscall.SIGTERM , syscall.SIGINT )
	go func() {
		<-c
		fmt.Println( "\r- Ctrl+C pressed in Terminal" )
		fmt.Printf( "Shutting Down %s Server\n" , s.Config.ServerName )
		s.DB.Close()
		s.FiberApp.Shutdown()
		os.Exit( 0 )
	}()
}

func main() {

	config_file_path := "./config.yaml"
	if len( os.Args ) > 1 { config_file_path , _ = filepath.Abs( os.Args[ 1 ] ) }
	config := utils.ParseConfig( config_file_path )

	SetupCloseHandler()
	// utils.GenerateNewKeys()
	utils.WriteLoginURLPrefix( config.ServerLoginUrlPrefix )
	s = server.New( config )
	fmt.Printf( "Loaded Config File From : %s\n" , config_file_path )
	s.Start()

}