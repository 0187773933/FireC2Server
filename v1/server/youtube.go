package server

import (
	fmt "fmt"
	time "time"
	"math/rand"
	"strings"
	// "context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	// redis "github.com/redis/go-redis/v9"
	circular_set "github.com/0187773933/RedisCircular/v1/set"
)

const YOUTUBE_ACTIVITY = "com.amazon.firetv.youtube/dev.cobalt.app.MainActivity"
const YOUTUBE_APP_NAME = "com.amazon.firetv.youtube"

func ( s *Server ) YouTubeReopenApp() {
	log.Debug( "YouTubeReopenApp()" )
	s.ADB.StopAllApps()
	s.ADB.Brightness( 0 )
	s.ADB.CloseAppName( YOUTUBE_APP_NAME )
	time.Sleep( 500 * time.Millisecond )
	s.ADB.OpenAppName( YOUTUBE_APP_NAME )
	log.Debug( "Done" )
}

func ( s *Server ) YouTubeContinuousOpen() {
	start_time_string , _ := utils.GetFormattedTimeStringOBJ()
	log.Debug( "YouTubeContinuousOpen()" )
	s.GetStatus()
	log.Debug( s.Status )
	s.Set( "active_player_name" , "youtube" )
	s.Set( "active_player_command" , "play" )
	s.Set( "active_player_start_time" , start_time_string )
	log.Debug( fmt.Sprintf( "Top Window Activity === %s" , s.Status.ADB.Activity ) )
	if s.Status.ADB.Activity == YOUTUBE_ACTIVITY {
		log.Debug( "youtube was already open" )
	} else {
		log.Debug( "youtube was NOT already open" )
		s.YouTubeReopenApp()
		time.Sleep( 500 * time.Millisecond )
	}
}

func ( s *Server ) YouTubeLiveNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "YouTubeLiveNext()" )
	s.YouTubeContinuousOpen()
	video_id := circular_set.Next( s.DB , "STATE.YOUTUBE.LIVE.VIDEOS" )
	uri := fmt.Sprintf( "https://www.youtube.com/watch?v=%s" , video_id )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
	// s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	return c.JSON( fiber.Map{
		"url": "/youtube/live/next" ,
		"video_id": video_id ,
		"result": true ,
	})
}

func ( s *Server ) YouTubeLivePrevious( c *fiber.Ctx ) ( error ) {
	log.Debug( "YouTubeLivePrevious()" )
	s.YouTubeContinuousOpen()
	video_id := circular_set.Previous( s.DB , "STATE.YOUTUBE.LIVE.VIDEOS" )
	uri := fmt.Sprintf( "https://www.youtube.com/watch?v=%s" , video_id )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
	// s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	return c.JSON( fiber.Map{
		"url": "/youtube/live/previous" ,
		"video_id": video_id ,
		"result": true ,
	})
}

func ( s *Server ) YouTubeVideo( c *fiber.Ctx ) ( error ) {
	video_id := c.Params( "video_id" )
	log.Debug( fmt.Sprintf( "YouTubeVideo( %s )" , video_id ) )
	s.YouTubeContinuousOpen()
	uri := fmt.Sprintf( "https://www.youtube.com/watch?v=%s" , video_id )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
	// s.ADB.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	return c.JSON( fiber.Map{
		"url": "/youtube/:video_id" ,
		"video_id": video_id ,
		"result": true ,
	})
}

// am start -a android.intent.action.VIEW -d "vnd.youtube:PLcW8xNfZoh7cgn3QNojv1ly5NZfg9xVN9?listType=playlist"
// https://www.youtube.com/playlist?list=PLcW8xNfZoh7cgn3QNojv1ly5NZfg9xVN9

type YoutubeVideoInfo struct {
	Items []struct {
		ID string `json:"id"`
	} `json:"items"`
}
func ( s *Server ) YouTubeIsVideoIdAvailable( video_id string ) ( result bool ) {
	result = false
	next_api_key := circular_set.Next( s.DB , "CONFIG.YOUTUBE.API_KEYS" )
	base_url := "https://youtube.googleapis.com/youtube/v3/videos"
	params := url.Values{}
	params.Add( "part" , "id" )
	params.Add( "id" , video_id )
	params.Add( "key" , next_api_key )
	full_url := fmt.Sprintf( "%s?%s" , base_url , params.Encode() )
	req, _ := http.NewRequest( "GET" , full_url , nil )
	req.Header.Add( "Accept" , "application/json" )
	resp , _ := http.DefaultClient.Do( req )
	defer resp.Body.Close()
	var video_info YoutubeVideoInfo
	body , _ := ioutil.ReadAll( resp.Body )
	json.Unmarshal( body , &video_info )
	if len( video_info.Items ) >= 1 {
		result = true
	}
	return
}


// 1.) Get Channels channelID
// curl \
// 'https://youtube.googleapis.com/youtube/v3/channels' \
// --header 'Accept: application/json' \
// --compressed \
// --get \
// --data-urlencode 'part=snippet' \
// --data-urlencode 'forUsername=MontereyBayAquarium' \
// --data-urlencode 'key=asdf'

// 2.) Get Channels Live Videos
// curl \
// 'https://youtube.googleapis.com/youtube/v3/search' \
// --header 'Accept: application/json' \
// --compressed \
// --get \
// --data-urlencode 'part=snippet' \
// --data-urlencode 'channelId=UCnM5iMGiKsZg-iOlIO2ZkdQ' \
// --data-urlencode 'eventType=live' \
// --data-urlencode 'maxResults=50' \
// --data-urlencode 'type=video' \
// --data-urlencode 'key=asdf'

// or just goto dev console on youtube channel's /stream page
// and run `ytInitialData.metadata.channelMetadataRenderer.externalId`
func ( s *Server ) YouTubeGetChannelId( channel_name string ) ( result string ) {
	next_api_key := circular_set.Next( s.DB , "CONFIG.YOUTUBE.API_KEYS" )
	headers := map[string]string{
		"Accept": "application/json" ,
	}
	params := map[string]string{
		"part": "snippet" ,
		"forUsername": channel_name ,
		"key": next_api_key ,
	}
	response_json := utils.GetJSON( "https://youtube.googleapis.com/youtube/v3/channels" , headers , params )
	items , _ := response_json.( map[string]interface{} )[ "items" ].( []interface{} )
	if len( items ) < 1 { fmt.Println( response_json ); return }
	first_result , _ := items[ 0 ].( map[string]interface{} )
	result , _ = first_result[ "id" ].( string )
	fmt.Println( result )
	return
}

type YoutubeVideo struct {
	Id string `json:"id"`
	Name string `json:"name"`
}
type YouTubeResponseItem struct {
	Id struct {
		VideoId string `json:"videoId"`
	} `json:"id"`
	Snippet struct {
		Title string `json:"title"`
	} `json:"snippet"`
}
type YoutubeResponse struct {
	Items []YouTubeResponseItem `json:"items"`
}
// https://developers.google.com/youtube/v3/docs/search/list
func ( s *Server ) YouTubeGetChannelsLiveVideos( channel_id string ) ( result []YoutubeVideo ) {
	next_api_key := circular_set.Next( s.DB , "CONFIG.YOUTUBE.API_KEYS" )
	base_url := "https://youtube.googleapis.com/youtube/v3/search"
	params := url.Values{}
	params.Add( "part" , "snippet" )
	params.Add( "channelId" , channel_id )
	params.Add( "eventType" , "live" )
	params.Add( "maxResults" , "50" )
	params.Add( "type" , "video" )
	params.Add( "key" , next_api_key )
	full_url := fmt.Sprintf( "%s?%s" , base_url , params.Encode() )
	max_retries := 1
	for i := 0; i < max_retries; i++ {
		req , _ := http.NewRequest( "GET" , full_url , nil )
		req.Header.Add( "Accept" , "application/json" )
		resp , _ := http.DefaultClient.Do( req )
		defer resp.Body.Close()
		body , _ := ioutil.ReadAll( resp.Body )
		var live_videos YoutubeResponse
		json.Unmarshal( body, &live_videos )
		if len( live_videos.Items ) >= 1 {
			for _ , item := range live_videos.Items {
				video := YoutubeVideo{
					Id: item.Id.VideoId,
					Name: item.Snippet.Title,
				}
				result = append( result , video )
			}
			return result
        } else {
            fmt.Println( "failed. Retrying..." )
            next_api_key = circular_set.Next( s.DB , "CONFIG.YOUTUBE.API_KEYS" )
            // fmt.Println( string( body ) )
            time.Sleep( 2 * time.Second )
        }
	}
	fmt.Println( "Max retries reached. No live videos found." )
	return
}

// Update DB With List of Currated Live Followers
// fucking idiots with this god damn quota. bro
func ( s *Server ) YouTubeLiveUpdate() ( result []string ) {
	// s.DB.Del( context.Background() , "STATE.YOUTUBE.LIVE.VIDEOS" )
	for channel_id , _ := range s.Config.Library.YouTube.Following.Live {
		fmt.Println( "\n" , channel_id , s.Config.Library.YouTube.Following.Live[ channel_id ].Name )
		live_videos := s.YouTubeGetChannelsLiveVideos( channel_id )
		for _ , video_item := range live_videos {
			for _ , video_name := range s.Config.Library.YouTube.Following.Live[ channel_id ].Videos {
				if strings.Contains( strings.ToLower( video_item.Name ) , video_name ) {
					fmt.Println( "adding" , video_item.Id , video_item.Name )
					result = append( result , video_item.Id )
				}
			}
		}
		time.Sleep( 1 * time.Second )
	}
	rand.Shuffle( len( result ) , func( i , j int ) {
		result[ i ] , result[ j ] = result[ j ] , result[ i ]
	})
	for _ , video_id := range result {
		circular_set.Add( s.DB , "STATE.YOUTUBE.LIVE.VIDEOS" , video_id )
	}
	return
}

func ( s *Server ) GetYouTubeLiveUpdate( c *fiber.Ctx ) ( error ) {
	live := s.YouTubeLiveUpdate()
	return c.JSON( fiber.Map{
		"url": "/youtube/live/update" ,
		"live": live ,
		"result": true ,
	})
}