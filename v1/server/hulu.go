package server

import (
	fmt "fmt"
	time "time"
	// url "net/url"
	// "math"
	// "strings"
	color "image/color"
	colorful "github.com/lucasb-eyer/go-colorful"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	// redis "github.com/redis/go-redis/v9"
	adb_wrapper "github.com/0187773933/ADBWrapper/v1/wrapper"
	circular_set "github.com/0187773933/RedisCircular/v1/set"
)

// https://www.hulu.com/series/502bbc34-fa19-48fb-89c6-074da28335d3
// document.querySelectorAll('div[data-testid="details-dropdown-container"]')?.[0]?.querySelectorAll('div[class*="dropdown-indicator"]')?.[0]?.click();
// ( async ()=> {
// 	function sleep( ms ) { return new Promise( resolve => setTimeout( resolve , ms ) ); }
// 	function wait_on_element( query_selector ) {
// 		return new Promise( function( resolve , reject ) {
// 			try {
// 				let READY_CHECK_INTERVAL = setInterval( function() {
// 					let item = document.querySelectorAll( query_selector );
// 					if ( item ) {
// 						if ( item[ 0 ] ) {
// 							clearInterval( READY_CHECK_INTERVAL );
// 							resolve( item[0] );
// 							return;
// 						}
// 					}
// 				} , 500 );
// 				setTimeout( function() {
// 					clearInterval( READY_CHECK_INTERVAL );
// 					resolve();
// 					return;
// 				} , 10000 );
// 			}
// 			catch( error ) { console.log( error ); reject( error ); return; }
// 		});
// 	}
// 	function escape_string( str ) {
// 	  return str
// 	    .replace(/\\/g, '\\\\') // Escape backslashes first
// 	    .replace(/"/g, '\\"') // Escape double quotes
// 	    .replace(/'/g, "\\'") // Escape single quotes
// 	    .replace(/\n/g, '\\n') // Escape newlines
// 	    .replace(/\r/g, '\\r') // Escape carriage returns
// 	    .replace(/\t/g, '\\t'); // Escape tabs
// 	}
// 	function random_int( min , max ) {
// 		min = Math.ceil( min );
// 		max = Math.floor( max );
// 		return Math.floor( Math.random() * ( max - min + 1 ) ) + min;
// 	}
// 	function click( element ) {
// 		[ "mouseenter" , "mouseover" , "mousemove" ].forEach( event_type => {
// 			const m_event = new MouseEvent( event_type , {
// 				view: window ,
// 				bubbles: true ,
// 				cancelable: false ,
// 				clientX: element.getBoundingClientRect().left ,
// 				clientY: element.getBoundingClientRect().top
// 			});
// 			element.dispatchEvent( m_event );
// 		});
// 		const mouse_click = new MouseEvent( "click" , {
// 			view: window ,
// 			bubbles: true ,
// 			cancelable: false
// 		});
// 		element.dispatchEvent( mouse_click );
// 	}
// 	function get_episodes() {
// 		let images = document.querySelectorAll( 'img[class^="StandardEmphasisHorizontalTileThumbnail"]' );
// 		let _ys = "";
// 		for ( let i = 0; i < images.length; ++i ) {
// 			if ( images[ i ].src === "" ) { continue; }
// 			let parts = images[ i ].src.split( "artwork/" );
// 			if ( parts.length < 2 ) { continue; }
// 			let uuid = parts[ 1 ].split( "?" )[ 0 ];
// 			let alt = images[ i ].alt.split( "Cover art for " )[ 1 ];
// 			if ( alt === undefined ) { continue; }
// 			alt = escape_string( alt );
// 			_ys += `          - id: ${uuid}\n`;
// 			_ys += `            name: "${alt}"\n`;
// 		}
// 		return _ys;
// 	}
// 	function get_seasons() {
// 		let seasons_dropdown = document.querySelectorAll( 'div[data-testid="details-dropdown-container"]' );
// 		if ( seasons_dropdown.length < 1 ) { return false; }
// 		seasons_dropdown = seasons_dropdown[ 0 ];
// 		let seasons_dropdown_arrow = seasons_dropdown.querySelectorAll( 'div[class*="dropdown-indicator"]' );
// 		if ( seasons_dropdown_arrow.length < 1 ) { return false; }
// 		seasons_dropdown_arrow = seasons_dropdown_arrow[ 0 ];
// 		seasons_dropdown_arrow.click();
// 		let seasons = seasons_dropdown.querySelectorAll( "ul" );
// 		if ( seasons.length < 1 ) { return false; }
// 		seasons = seasons[ 0 ];
// 		seasons = seasons.querySelectorAll( "li" );
// 		let first_season_number = parseInt( seasons[ 0 ].id.split( "::" )[ 1 ] );
// 		let last_season = parseInt( seasons[ seasons.length - 1 ].id.split( "::" )[ 1 ] );
// 		if ( first_season_number > last_season ) {
// 			seasons = [...seasons].reverse();
// 		}
// 		seasons_dropdown_arrow.click();
// 		return seasons_dropdown_arrow , seasons;
// 	}
// 	async function select_season( season_index ) {
// 		let seasons_dropdown = document.querySelectorAll( 'div[data-testid="details-dropdown-container"]' );
// 		if ( seasons_dropdown.length < 1 ) { return false; }
// 		seasons_dropdown = seasons_dropdown[ 0 ];
// 		let seasons_dropdown_arrow = seasons_dropdown.querySelectorAll( 'div[class*="dropdown-indicator"]' );
// 		if ( seasons_dropdown_arrow.length < 1 ) { return false; }
// 		seasons_dropdown_arrow = seasons_dropdown_arrow[ 0 ];
// 		seasons_dropdown_arrow.click();
// 		let seasons = seasons_dropdown.querySelectorAll( "ul" );
// 		if ( seasons.length < 1 ) { return false; }
// 		seasons = seasons[ 0 ];
// 		seasons = seasons.querySelectorAll( "li" );
// 		let first_season_number = parseInt( seasons[ 0 ].id.split( "::" )[ 1 ] );
// 		let last_season = parseInt( seasons[ seasons.length - 1 ].id.split( "::" )[ 1 ] );
// 		if ( first_season_number > last_season ) {
// 			seasons = [...seasons].reverse();
// 		}
// 		await sleep( random_int( 1000 , 1600 ) );
// 		console.log( "trying to select :" , seasons[ season_index ].innerText );
// 		console.log( seasons[ season_index ] );
// 		click( seasons[ season_index ] );
// 	}
// 	function get_id() {
// 		const url = window.location.href;
// 		const parts = url.split( "/" );
// 		const uuid = parts[ parts.length - 1 ];
// 		return uuid
// 	}
// 	let id = get_id();
// 	let show_name = document.getElementById( "dialog-title" ).innerText;
// 	show_name = escape_string( show_name );
// 	let seasons = get_seasons();
// 	let yaml_string = "";
// 	yaml_string += `  ${id}:\n`;
// 	yaml_string += `    name: "${show_name}"\n`;
// 	yaml_string += `    seasons:\n`;
// 	let total_seasons = seasons.length;
// 	console.log( `total seasons === ${total_seasons}` );
// 	console.log( seasons );
// 	if ( total_seasons === undefined || total_seasons === 0 ) {
// 		console.log( "probably only 1 season in show" );
// 		let episodes = get_episodes();
// 		if ( episodes === "" ) { console.log( yaml_string ); return }
// 		yaml_string += `      - number: "one"\n`;
// 		yaml_string += `        episodes:\n`;
// 		yaml_string += episodes;
// 		console.log( yaml_string );
// 		return;
// 	}
// 	for ( let i = 0; i < seasons.length; ++i ) {
// 		console.log( `getting season [${(i+1)}] of ${total_seasons}` );
// 		let season = seasons[ i ];
// 		console.log( "selecting from dropdown" );
// 		select_season( i );
// 		await sleep( random_int( 2000 , 2600 ) );
// 		await wait_on_element( 'img[class^="StandardEmphasisHorizontalTileThumbnail"]' );
// 		let episodes = get_episodes();
// 		if ( episodes === "" ) { continue; }
// 		yaml_string += `      - number: "${i+1}"\n`;
// 		yaml_string += `        episodes:\n`;
// 		yaml_string += episodes;
// 		await sleep( random_int( 1200 , 1600 ) );
// 	}
// 	console.log( yaml_string );
// })();

// https://community.home-assistant.io/t/androidtv-autoplay-hulu-series-where-you-left-off/531105
// adb shell am start -n com.hulu.livingroomplus/.WKFactivity
// -a hulu.intent.action.PLAY_CONTENT -e content_id 502bbc34-fa19-48fb-89c6-074da28335d3

// am start -n hulu.intent.action.PLAY_CONTENT -e content_id 502bbc34-fa19-48fb-89c6-074da28335d3
// start -n com.hulu.plus/.WKFactivity -a hulu.intent.action.PLAY_CONTENT -e content_id 502bbc34-fa19-48fb-89c6-074da28335d3

// const HULU_ACTIVITY = "com.hulu.plus/com.hulu.plus.MainActivity"
// const HULU_APP_NAME = "com.hulu.plus"

func ( s *Server ) HuluReopenApp() {
	log.Debug( "HuluReopenApp()" )
	s.ADB.StopAllPackages()
	// s.ADB.SetBrightness( 0 )
	s.ADB.ClosePackage( s.Config.ADB.APKS[ "hulu" ][ s.Config.ADB.DeviceType ].Package )
	time.Sleep( 1000 * time.Millisecond )
	s.ADB.OpenPackage( s.Config.ADB.APKS[ "hulu" ][ s.Config.ADB.DeviceType ].Package )
	log.Debug( "Done" )
}

func ( s *Server ) HuluContinuousOpen() {
	start_time_string , _ := utils.GetFormattedTimeStringOBJ()
	log.Debug( "HuluContinuousOpen()" )
	s.GetStatus()
	log.Debug( s.Status )
	s.Set( "active_player_name" , "hulu" )
	s.Set( "active_player_command" , "play" )
	s.Set( "active_player_start_time" , start_time_string )
	log.Debug( fmt.Sprintf( "Top Window Activity === %s" , s.Status.ADB.Activity ) )
	s.ADBWakeup()
	// if strings.Contains( s.Status.ADB.Activity , "hulu" ) {
	// 	log.Debug( fmt.Sprintf( "hulu was already open with activity %s" , s.Status.ADB.Activity ) )
	// 	switch s.Config.ADB.DeviceType {
	// 		case "firecube" , "firestick":
	// 			return;
	// 			break;
	// 		case "firetablet":
	// 			log.Debug( "restarting anyway" )
	// 			break;
	// 	}
	// }
	log.Debug( "hulu was NOT already open" )
	log.Debug( "we have to force app to reopen to get correct deep uri linking" )
	s.HuluReopenApp()
}

func ( s *Server ) HuluSelectProfile() {
	log.Debug( fmt.Sprintf( "HuluSelectProfile( %d )" , s.Config.HuluUserProfileIndex ) )
	for i := 0; i < s.Config.HuluTotalUserProfiles; i++ {
		s.ADB.Down()
		time.Sleep( 200 * time.Millisecond )
	}
	s.ADB.Up()
	time.Sleep( 200 * time.Millisecond )
	s.ADB.Enter()
}

func parse_hulu_sent_id( sent_id string ) ( uri string ) {
	if utils.IsUUID( sent_id ) {
		uri = fmt.Sprintf( "https://www.hulu.com/watch/%s" , sent_id )
		return
	}
	is_url , _ := utils.IsURL( sent_id )
	if is_url {
		uri = sent_id
		return
	}
	return
}

func ( s *Server ) HuluOpenURI( uri string ) {
	log.Debug( fmt.Sprintf( "HuluOpenURI( %s )" , uri ) )
	switch s.Config.ADB.DeviceType {
		case "firecube" , "firestick":
			s.HuluContinuousOpen()
			time.Sleep( 1 * time.Second )
			s.ADB.OpenURI( uri )
			time.Sleep( 4 * time.Second )
			break;
		case "firetablet":
			s.GetStatus()
			s.ADBWakeup()
			s.ADB.StopAllPackages()
			s.ADB.ClosePackage( s.Config.ADB.APKS[ "hulu" ][ s.Config.ADB.DeviceType ].Package )
			time.Sleep( 1 * time.Second )
			s.ADB.OpenURI( uri )
			time.Sleep( 4 * time.Second )
			break;
	}
	// flatten pixels for adb.GetPixelColorsFromImageBytes()
	// otherwise use png.Decode here
	var profile_screen_pixel_colors []color.RGBA
	var profile_screen_pixel_coords []adb_wrapper.Coord
	for _ , coord := range s.Config.ADB.APKS[ "hulu" ][ s.Config.ADB.DeviceType ].Pixels[ "profile_selection" ] {
		profile_screen_pixel_coords = append( profile_screen_pixel_coords , adb_wrapper.Coord{ X: coord.X , Y: coord.Y } )
		c , _ := colorful.Hex( coord.Color )
		r , g , b := c.RGB255()
		profile_screen_pixel_colors = append( profile_screen_pixel_colors , color.RGBA{ R: r , G: g , B: b , A: 255 } )
	}

	verified_now_playing := false
	verified_now_playing_updated_time := 0
	queries := 20
	stage_one_ready := false
	observed_count := 0
	login_observed_count := 0
	normal_observed_count := 0
	for i := 0; i < queries; i++ {
		// status := s.ADB.GetStatus()
		if stage_one_ready == true { break; }
		log.Debug( fmt.Sprintf( "checking [%d] of %d for hulu to be ready" , ( i + 1 ) , queries ) )
		// activity := s.ADB.GetActivity()
		players := s.ADB.FindPlayers( "hulu" )
		if len( players ) > 0 {
			log.Debug( "found hulu player" )
			observed_count = observed_count + 1
			if observed_count > 3 {
				for _ , player := range players {
					if player.Updated > 0 {
						if player.Position > 0 {
							log.Debug( "hulu autostarted playing on it's own" )
							verified_now_playing = true
							verified_now_playing_updated_time = player.Updated
							stage_one_ready = true
							break
						}
					}
				}
			}
		}
		fmt.Println( "observed_count ===" , observed_count )
		fmt.Println( verified_now_playing )
		fmt.Println( verified_now_playing_updated_time )
		switch s.Config.ADB.DeviceType {
			case "firecube":
				screenshot_bytes := s.ADB.ScreenshotToBytes()
				test_colors := s.ADB.GetPixelColorsFromImageBytes( &screenshot_bytes , profile_screen_pixel_coords )
				fmt.Println( test_colors )
				login_screen := false
				for i , test_color := range test_colors {
					if test_color != profile_screen_pixel_colors[ i ] {
						log.Debug( "different color pixel found than on known login screen" )
						login_screen = false
					} else {
						login_screen = true
					}
				}
				if login_screen == true {
					log.Debug( "all test pixels matched , we are on the profile selection screen" )
					login_observed_count += 1
					if login_observed_count >= 3 {
						log.Debug( "3 times saw green pixel ... so we are on the profile selection screen" )
						stage_one_ready = true
						s.HuluSelectProfile()
						log.Debug( "need to double check if it started auto playing still" )
						break;
					}
				} else {
					log.Debug( "we are not on the profile selection screen" )
					normal_observed_count += 1
					if normal_observed_count >= 2 {
						log.Debug( "already observed 2 times , pressing enter , tapping play button" )
						s.ADB.Enter()
						s.ADB.Play()
						stage_one_ready = true
						break;
					}
				}
				break
			case "firestick":
				log.Debug( "sleeping 10 seconds to wait on app init" )
				time.Sleep( 10 * time.Second )
				break;
			case "firetablet":
				if i > 1 && i & 3 == 0 {
					log.Debug( "already observed 3 times , pressing enter , tapping play button" )
					s.ADB.Enter()
					s.ADB.Play()
					s.ADB.Tap( 245 , 367 )
					stage_one_ready = true
					time.Sleep( 2 * time.Second )
					log.Debug( "attempting to enter full screen" )
					s.ADB.Tap( 305 , 178 ) // tap center of video window to activate ui overlay
					time.Sleep( 300 * time.Millisecond )
					s.ADB.Tap( 564 , 60 ) // tap fullscreen button
				}
				break;
		}
		time.Sleep( 1 * time.Second )
	}
	time.Sleep( 4 * time.Second )
	status := s.ADB.GetStatus()
	utils.PrettyPrint( status )
	// log.Debug( "waiting 10 seconds to see if hulu auto starts playing" )
	// playing := s.ADB.WaitOnPlayersPlaying( "hulu" , 10 )
	// if len( playing ) < 1 {
	// 	log.Debug( "never started playing , we might have to try play button" )
	// }
	// utils.PrettyPrint( playing )
	// log.Debug( fmt.Sprintf( "total now playing === %d" , len( playing ) ) )
	// for _ , player := range playing {
	// 	if player.Updated > 0 {
	// 		log.Debug( "hulu autostarted playing on it's own" )
	// 		verified_now_playing = true
	// 		verified_now_playing_updated_time = player.Updated
	// 		break
	// 	}
	// }
	// if verified_now_playing == false {
	// 	log.Debug( "hulu didn't auto start playing , we might have to try play button" )
	// 	return
	// }
	// log.Debug( "waiting now for player progress" )
	// progress := s.ADB.WaitOnPlayersUpdated( "hulu" , verified_now_playing_updated_time , 10 )
	// utils.PrettyPrint( progress )
	// log.Debug( "player progress should be ready" )
	// switch s.Config.ADB.DeviceType {
	// 	case "firecube" , "firestick":
	// 		break;
	// 	case "firetablet":
	// 		log.Debug( "attempting to enter full screen" )
	// 		s.ADB.Tap( 305 , 178 ) // tap center of video window to activate ui overlay
	// 		time.Sleep( 300 * time.Millisecond )
	// 		s.ADB.Tap( 564 , 60 ) // tap fullscreen button
	// 		break;
	// }
}

func ( s *Server ) HuluID( c *fiber.Ctx ) ( error ) {
	sent_id := c.Params( "*" )
	sent_query := c.Request().URI().QueryArgs().String()
	if sent_query != "" { sent_id += "?" + sent_query }
	uri := parse_hulu_sent_id( sent_id )
	log.Debug( fmt.Sprintf( "HuluID( %s )" , uri ) )
	s.HuluOpenURI( uri )
	s.Set( "active_player_now_playing_id" , sent_id )
	s.Set( "active_player_now_playing_uri" , sent_id )
	return c.JSON( fiber.Map{
		"url": "/hulu/:id" ,
		"id": sent_id ,
		"result": true ,
	})
}

func ( s *Server ) HuluMovieNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "HuluMovieNext()" )
	next_movie := circular_set.Next( s.DB , "LIBRARY.HULU.MOVIES" )
	next_movie_name := s.Config.Library.Hulu.Movies[ next_movie ].Name
	uri := fmt.Sprintf( "https://www.hulu.com/movie/%s" , next_movie )
	s.HuluOpenURI( uri )
	s.Set( "STATE.HULU.NOW_PLAYING" , next_movie )
	s.Set( "active_player_now_playing_id" , next_movie )
	s.Set( "active_player_now_playing_text" , next_movie_name )
	return c.JSON( fiber.Map{
		"url": "/hulu/next" ,
		"id": next_movie ,
		"name": next_movie_name ,
		"result": true ,
	})
}

func ( s *Server ) HuluMoviePrevious( c *fiber.Ctx ) ( error ) {
	log.Debug( "HuluMoviePrevious()" )
	next_movie := circular_set.Previous( s.DB , "LIBRARY.HULU.MOVIES" )
	next_movie_name := s.Config.Library.Hulu.Movies[ next_movie ].Name
	uri := fmt.Sprintf( "https://www.hulu.com/movie/%s" , next_movie )
	s.HuluOpenURI( uri )
	s.Set( "STATE.HULU.NOW_PLAYING" , next_movie )
	s.Set( "active_player_now_playing_id" , next_movie )
	s.Set( "active_player_now_playing_text" , next_movie_name )
	return c.JSON( fiber.Map{
		"url": "/hulu/next" ,
		"id": next_movie ,
		"name": next_movie_name ,
		"result": true ,
	})
}

func ( s *Server ) HuluTVID( c *fiber.Ctx ) ( error ) {
	series_id := c.Params( "series_id" )
	// log.Debug( fmt.Sprintf( "HuluTVID( %s )" , series_id ) )
	if series_id == "" {
		return c.JSON( fiber.Map{
			"url": "/hulu/tv/:series_id" ,
			"series_id": series_id ,
			"result": false ,
		})
	}
	_ , series_exists := s.Config.Library.Hulu.TV[ series_id ]
	if series_exists == false {
		return c.JSON( fiber.Map{
			"url": "/hulu/tv/:series_id" ,
			"series_id": series_id ,
			"error": "series doesn't exist in library" ,
			"result": false ,
		})
	}
	// s.Set( "STATE.HULU.NOW_PLAYING.MODE" , "TV" )
	s.Set( "STATE.HULU.NOW_PLAYING.TV.SERIES_ID" , series_id )
	next_episode := circular_set.Next( s.DB , fmt.Sprintf( "LIBRARY.HULU.TV.%s" , series_id ) )
	next_episode_name := s.Get( fmt.Sprintf( "LIBRARY.HULU.TV.%s.%s" , series_id , next_episode ) )
	log.Debug( fmt.Sprintf( "HuluTVID( %s ) --> %s === " , series_id , next_episode , next_episode_name ) )
	uri := fmt.Sprintf( "https://www.hulu.com/watch/%s" , next_episode )
	log.Debug( uri )
	s.HuluOpenURI( uri )
	s.Set( "active_player_now_playing_id" , next_episode )
	s.Set( "active_player_now_playing_uri" , uri )
	return c.JSON( fiber.Map{
		"url": "/hulu/tv/:series_id" ,
		"series_id": series_id ,
		"next_episode_id": next_episode ,
		"next_episode_name": next_episode_name ,
		"result": true ,
	})
}

func ( s *Server ) HuluTVNext( c *fiber.Ctx ) ( error ) {
	log.Debug( "HuluTVNext()" )
	series_id := s.Get( "STATE.HULU.NOW_PLAYING.TV.SERIES_ID" )
	next_episode := circular_set.Next( s.DB , fmt.Sprintf( "LIBRARY.HULU.TV.%s" , series_id ) )
	next_episode_name := s.Get( fmt.Sprintf( "LIBRARY.HULU.TV.%s.%s" , series_id , next_episode ) )
	log.Debug( fmt.Sprintf( "HuluTVNext( %s ) --> %s === " , series_id , next_episode , next_episode_name ) )
	uri := fmt.Sprintf( "https://www.hulu.com/watch/%s" , next_episode )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
	s.Set( "active_player_now_playing_id" , next_episode )
	s.Set( "active_player_now_playing_uri" , uri )
	return c.JSON( fiber.Map{
		"url": "/hulu/tv/:id/next" ,
		"series_id": series_id ,
		"next_episode_id": next_episode ,
		"next_episode_name": next_episode_name ,
		"result": true ,
	})
}

func ( s *Server ) HuluTVPrevious( c *fiber.Ctx ) ( error ) {
	log.Debug( "HuluTVPrevious()" )
	series_id := s.Get( "STATE.HULU.NOW_PLAYING.TV.SERIES_ID" )
	next_episode := circular_set.Previous( s.DB , fmt.Sprintf( "LIBRARY.HULU.TV.%s" , series_id ) )
	next_episode_name := s.Get( fmt.Sprintf( "LIBRARY.HULU.TV.%s.%s" , series_id , next_episode ) )
	log.Debug( fmt.Sprintf( "HuluTVPrevious( %s ) --> %s === " , series_id , next_episode , next_episode_name ) )
	uri := fmt.Sprintf( "https://www.hulu.com/watch/%s" , next_episode )
	log.Debug( uri )
	s.ADB.OpenURI( uri )
	s.Set( "active_player_now_playing_id" , next_episode )
	s.Set( "active_player_now_playing_uri" , uri )
	return c.JSON( fiber.Map{
		"url": "/hulu/tv/:id/previous" ,
		"series_id": series_id ,
		"next_episode_id": next_episode ,
		"next_episode_name": next_episode_name ,
		"result": true ,
	})
}