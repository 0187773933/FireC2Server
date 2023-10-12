package server

import (
	fmt "fmt"
	time "time"
	// url "net/url"
	// "math"
	// "image/color"
	"strings"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	// redis "github.com/redis/go-redis/v9"
	circular_set "github.com/0187773933/RedisCircular/v1/set"
)



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
	headers := map[string]string{
		"Accept": "application/json" ,
	}
	params := map[string]string{
		"part": "snippet" ,
		"forUsername": channel_name ,
		"key": s.Config.YouTubeAPIKeyOne ,
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
	base_url := "https://youtube.googleapis.com/youtube/v3/search"
	params := url.Values{}
	params.Add( "part" , "snippet" )
	params.Add( "channelId" , channel_id )
	params.Add( "eventType" , "live" )
	params.Add( "maxResults" , "50" )
	params.Add( "type" , "video" )
	params.Add( "key" , s.Config.YouTubeAPIKeyOne )
	full_url := fmt.Sprintf( "%s?%s" , base_url , params.Encode() )
	max_retries := 1
	for i := 0; i < max_retries; i++ {
		req, _ := http.NewRequest( "GET" , full_url , nil )
		req.Header.Add( "Accept" , "application/json" )
		resp, _ := http.DefaultClient.Do( req )
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll( resp.Body )
		var live_videos YoutubeResponse
		json.Unmarshal( body, &live_videos )
		if len( live_videos.Items ) >= 1 {
			for _, item := range live_videos.Items {
				video := YoutubeVideo{
					Id:   item.Id.VideoId,
					Name: item.Snippet.Title,
				}
				result = append( result , video )
			}
			return result
        } else {
            fmt.Println( "failed. Retrying..." )
            // fmt.Println( string( body ) )
            time.Sleep( 2 * time.Second )
        }
	}
	fmt.Println( "Max retries reached. No live videos found." )
	return
}

// Update DB With List of Currated Live Followers
func ( s *Server ) YouTubeLiveUpdate() ( result []YoutubeVideo ) {
	var ctx = context.Background()
	s.DB.Del( ctx , "STATE.YOUTUBE.LIVE.VIDEOS" )
	for channel_id , _ := range s.Config.Library.YouTube.Following.Live {
		fmt.Println( "\n" , channel_id , s.Config.Library.YouTube.Following.Live[ channel_id ].Name )
		live_videos := s.YouTubeGetChannelsLiveVideos( channel_id )
		for _ , video_item := range live_videos {
			for _ , video_name := range s.Config.Library.YouTube.Following.Live[ channel_id ].Videos {
				if strings.Contains( strings.ToLower( video_item.Name ) , video_name ) {
					fmt.Println( "adding" , video_item.Id , video_item.Name )
					circular_set.Add( s.DB , "STATE.YOUTUBE.LIVE.VIDEOS" , video_item.Id )
				}
			}
		}
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