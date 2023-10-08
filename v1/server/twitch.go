package server

import (
	fmt "fmt"
	time "time"
	context "context"
	sort "sort"
	// url "net/url"
	// "math"
	// "image/color"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	redis "github.com/redis/go-redis/v9"
	circular_set "github.com/0187773933/RedisCircular/v1/set"
)

const TWITCH_ACTIVITY = "tv.twitch.android.viewer/tv.twitch.starshot64.app.StarshotActivity"
const TWITCH_APP_NAME = "tv.twitch.android.viewer"

func ( s *Server ) TwitchReopenApp() {
	log.Debug( "TwitchReopenApp()" )
	s.ADB.StopAllApps()
	s.ADB.Brightness( 0 )
	s.ADB.CloseAppName( TWITCH_APP_NAME )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.OpenAppName( TWITCH_APP_NAME )
	log.Debug( "Done" )
}

func ( s *Server ) TwitchContinuousOpen() {
	start_time_string , _ := utils.GetFormattedTimeStringOBJ()
	log.Debug( "TwitchContinuousOpen()" )
	s.GetStatus()
	log.Debug( s.Status )
	s.Set( "active_player_name" , "twitch" )
	s.Set( "active_player_command" , "play" )
	s.Set( "active_player_start_time" , start_time_string )
	log.Debug( fmt.Sprintf( "Top Window Activity === %s" , s.Status.ADBTopWindow ) )
	if s.Status.ADBTopWindow != TWITCH_ACTIVITY {
		log.Debug( "twitch was NOT already open" )
		s.TwitchReopenApp()
		time.Sleep( 500 * time.Millisecond )
	} else {
		log.Debug( "twitch was already open" )
	}
}

func ( s *Server ) TwitchLiveNext( c *fiber.Ctx ) ( error ) {

	s.TwitchContinuousOpen()

	next_stream := circular_set.Next( s.DB , "STATE.TWITCH.FOLLOWING.LIVE" )
	log.Debug( "Next === " , next_stream )
	if next_stream == "" {
		log.Debug( "Empty , Refreshing" )
		s.TwitchLiveUpdate()
		next_stream = circular_set.Next( s.DB , "STATE.TWITCH.FOLLOWING.LIVE" )
		if next_stream == "" {
			log.Debug( "nobody is live ...." )
			return c.JSON( fiber.Map{
				"url": "/twitch/live/next" ,
				"stream": "nobody is live ...." ,
				"result": false ,
			})
		}
	}
	// force refresh on last
	first_in_set_z , _ := s.DB.ZRangeWithScores( context.Background() , "STATE.TWITCH.FOLLOWING.LIVE" , 0 , 0 ).Result()
	first_in_set := first_in_set_z[0].Member.(string)
	if first_in_set == next_stream {
		log.Debug( "recycled list" )
		s.TwitchLiveUpdate()
		// next_stream = circular_set.Next( s.DB , "STATE.TWITCH.FOLLOWING.LIVE" )
	}
	// var context = context.Background()
	// next_stream , _ := s.DB.LPop( context , "STATE.TWITCH.FOLLOWING.LIVE" ).Result()
	// if next_stream == "" {
	// 	log.Debug( "Live Stream List Empty , Refreshing" )
	// 	s.TwitchLiveUpdate()
	// 	next_stream , _ = s.DB.LPop( context , "STATE.TWITCH.FOLLOWING.LIVE" ).Result()
	// 	if next_stream == "" {
	// 		log.Debug( "nobody is live ...." )
	// 		return c.JSON( fiber.Map{
	// 			"url": "/twitch/live/next" ,
	// 			"stream": "nobody is live ...." ,
	// 			"result": false ,
	// 		})
	// 	}
	// }
	log.Debug( fmt.Sprintf( "TwitchLiveNext( %s )" , next_stream ) )
	fmt.Println( next_stream )
	// // TODO === Need to Add TV Mute and Unmute
	// // TODO === If Same Playlist don't open , just press next ? depends
	// s.ADB.SetVolume( 0 )
	uri := fmt.Sprintf( "twitch://stream/%s" , next_stream )
	s.ADB.OpenURI( uri )
	s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	return c.JSON( fiber.Map{
		"url": "/twitch/live/next" ,
		"stream": next_stream ,
		"result": true ,
	})
}


// Auth Option - 1 = Client Credentials
// curl -X POST 'https://id.twitch.tv/oauth2/token' \
// -H 'Content-Type: application/x-www-form-urlencoded' \
// -d 'client_id=asdf&client_secret=asdf&grant_type=client_credentials'

// Auth Option - 2 - Step - 1 = Create a Toekn that Can Automatically be Refreshed
// 2 steps , first step is creating a "code"
// https://id.twitch.tv/oauth2/authorize
//     ?response_type=code
//     &client_id=asdf
//     &redirect_uri=http://localhost:9371
//     &scope=user:read:follows

// Auth Option - 2 - Step - 2 = Second step is exchaning the "code" for a "token"
// This HAS to be a post request
// curl -X POST 'https://id.twitch.tv/oauth2/token' \
// -H 'Content-Type: application/x-www-form-urlencoded' \
// -d 'client_id=asdf' \
// -d 'client_secret=asdf' \
// -d 'code=asdf' \
// -d 'grant_type=authorization_code' \
// -d 'redirect_uri=http://localhost:9371'
	// returns :
	// {"access_token":"asdf","expires_in":13832,"refresh_token":"asdf","scope":["user:read:follows"],"token_type":"bearer"}

// Auth Option - 2 - Repeated Step - Refresh = Now we can just call refresh every time we need to make a batch of api calls
// curl -X POST https://id.twitch.tv/oauth2/token \
// -H 'Content-Type: application/x-www-form-urlencoded' \
// -d 'grant_type=refresh_token' \
// -d 'refresh_token=asdf' \
// -d 'client_id=asdf' \
// -d 'client_secret=asdf'
	// So you HAVE to store this refresh_token
	// update it every time

// Auth Option - 3 = Create Specific Auth Tokens
// These are temporary , and can't
// https://dev.twitch.tv/docs/authentication/scopes/
// https://id.twitch.tv/oauth2/authorize
//     ?response_type=token
//     &client_id=asdf
//     &redirect_uri=http://localhost:9371
//     &scope=user:read:follows

// Token Validation
// curl -X GET 'https://id.twitch.tv/oauth2/validate' \
// -H 'Authorization: OAuth asdf'

// Get User-ID of Username
// curl -X GET 'https://api.twitch.tv/helix/users?login=asdf' \
// -H 'Client-Id: asdf' \
// -H 'Authorization: Bearer Client-Credentials-Auth-Token'

// Get the channels you follow which are currently live
// curl -X GET 'https://api.twitch.tv/helix/streams/followed?user_id=asdf' \
// -H 'Authorization: Bearer asdf' \
// -H 'Client-Id: asdf'

func ( s *Server ) TwitchRefreshAuthToken() ( access_token string ) {
	access_token = s.Get( "STATE.TWITCH.ACCESS_TOKEN" )
	refresh_token := s.Get( "STATE.TWITCH.REFRESH_TOKEN" )
	if access_token == "" || refresh_token == "" {
		access_token = s.Config.TwitchAccessToken
		refresh_token = s.Config.TwitchRefreshToken
	}
	headers := map[string]string{}
	data := map[string]string{
		"grant_type": "refresh_token" ,
		"refresh_token": refresh_token ,
		"client_id": s.Config.TwitchClientID ,
		"client_secret": s.Config.TwitchClientSecret ,
	}
	new := utils.PostJSON( "https://id.twitch.tv/oauth2/token" , headers , data );
	new_access_token := new.(map[string]interface{})[ "access_token" ]
	new_refresh_token := new.(map[string]interface{})[ "refresh_token" ]
	if new_access_token != nil && new_refresh_token != nil {
		access_token = new_access_token.(string)
		refresh_token = new_refresh_token.(string)
	}
	s.Set( "STATE.TWITCH.ACCESS_TOKEN" , access_token )
	s.Set( "STATE.TWITCH.REFRESH_TOKEN" , refresh_token )
	return
}


// Update DB With List of Currated Live Followers
func ( s *Server ) TwitchLiveUpdate() ( result []string ) {
	// 1.) Get Currated List
	var context = context.Background()
	currated_followers , _ := s.DB.ZRangeByScore(
		context ,
		"LIBRARY.TWITCH.FOLLOWING.CURRATED" ,
		&redis.ZRangeBy{
			Min: "-inf" ,
			Max: "+inf" ,
		} ,
	).Result()

	// 2.) Get Live Followers
	access_token := s.TwitchRefreshAuthToken()
	headers := map[string]string{
		"Authorization": fmt.Sprintf( "Bearer %s" , access_token ) ,
		"Client-Id": s.Config.TwitchClientID ,
	}
	params := map[string]string{
		"user_id": s.Config.TwitchUserID ,
	}
	live_followers_json := utils.GetJSON( "https://api.twitch.tv/helix/streams/followed" , headers , params )
	live_followers := live_followers_json.(map[string]interface{})[ "data" ].( []interface{} )

	// 3.) Get Live Currated List
	live_index_map := make(map[string]int)
	for i , user := range live_followers {
		user_name := user.(map[string]interface{})["user_login"].(string)
		live_index_map[user_name] = i;

	} // lookup table
	s.DB.Del( context , "STATE.TWITCH.FOLLOWING.LIVE" )
	for _ , user := range currated_followers {
		if _ , exists := live_index_map[ user ]; exists {
			result = append( result , user )
			// s.DB.RPush( context , "STATE.TWITCH.FOLLOWING.LIVE" , user )
			circular_set.Add( s.DB , "STATE.TWITCH.FOLLOWING.LIVE" , user )
		}
	}


	currated_map := make(map[string]int)
	for i, user := range currated_followers {
		currated_map[user] = i
	}
	sort.Slice(result, func(i, j int) bool {
		return currated_map[result[i]] < currated_map[result[j]]
	})


	log.Debug( result )
	return
}

func ( s *Server ) GetTwitchLiveUpdate( c *fiber.Ctx ) ( error ) {
	live := s.TwitchLiveUpdate()
	return c.JSON( fiber.Map{
		"url": "/twitch/live/update" ,
		"currated": live ,
		"result": true ,
	})
}