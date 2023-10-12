package server

import (
	// "fmt"
	// "strings"
	// fiber "github.com/gofiber/fiber/v2"
)

func ( s *Server ) SetupAdminRoutes() {

	// Generic ADB
	s.FiberApp.Get( "play" , s.Play )
	s.FiberApp.Get( "pause" , s.Pause )
	s.FiberApp.Get( "stop" , s.Stop )
	s.FiberApp.Get( "next" , s.Next )
	s.FiberApp.Get( "previous" , s.Previous )

	// Spotify
	spotify := s.FiberApp.Group( "/spotify" )
	spotify.Use( validate_admin_mw )
	spotify.Get( "/shuffle/:state" , s.SpotifySetShuffle )
	spotify.Get( "/song/:song_id" , s.SpotifySong )
	spotify.Get( "/playlist/:playlist_id" , s.SpotifyPlaylist )
	spotify.Get( "/playlist-shuffle/:playlist_id" , s.SpotifyPlaylistWithShuffle )
	// spotify.Get( "/next/song" , SpotifyPlaylistWithShuffle )
	// spotify.Get( "/next/playlist" , SpotifyNextPlaylist )
	spotify.Get( "/next/playlist-shuffle" , s.SpotifyNextPlaylistWithShuffle )
	// spotify.Get( "/previous/song" , SpotifyPlaylistWithShuffle )
	// spotify.Get( "/previous/playlist" , SpotifyPlaylistWithShuffle )
	// spotify.Get( "/previous" , SpotifyPressPreviousButton ) // needs a custom previous , requires 2 clicks if in shuffle-mode

	// Twitch
	twitch := s.FiberApp.Group( "/twitch" )
	twitch.Use( validate_admin_mw )
	twitch.Get( "/next" , s.TwitchLiveNext )
	twitch.Get( "/previous" , s.TwitchLivePrevious )
	twitch.Get( "/update" , s.GetTwitchLiveUpdate )

	// Disney
	disney := s.FiberApp.Group( "/disney" )
	disney.Use( validate_admin_mw )
	disney.Get( "/next" , s.DisneyMovieNext )
	disney.Get( "/previous" , s.DisneyMoviePrevious )
	disney.Get( "/movie/:movie_id" , s.DisneyMovie )


	// YouTube
	youtube := s.FiberApp.Group( "/youtube" )
	youtube.Use( validate_admin_mw )
	youtube.Get( "/update/live" , s.GetYouTubeLiveUpdate )
	// s.SetupMediaPlayerRoutes( youtube , "youtube" )

}