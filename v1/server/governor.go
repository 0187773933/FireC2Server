package server

import (
	"fmt"
	// "wg"
	// "sync"
	"time"
	utils "github.com/0187773933/FireC2Server/v1/utils"
)

func ( s *Server ) Governor() {
	ticker := time.NewTicker( 30 * time.Second )
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:

			fmt.Println( "Mutex acquired by governor, performing checks" )

			// TODO
			tv_status := s.TV.Status()
			utils.PrettyPrint( tv_status )


			fmt.Println( "Mutex released by governor" )
		}
	}
}