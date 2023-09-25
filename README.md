# Fire C2 Server

## Example

- http://localhost:5954/spotify/playlist-shuffle/3UMDmO2YJb8DgUjpSBu8y9?k=asdf

## TODO

- switching back to redis
	- https://github.com/0187773933/RedisCircularList
- Find some new way to get the config and db reference into logger package.
	- apparently there is init() in go? that gets called before main starts ...
- Audio playing locally through speakers
- HTML Control Panels
- TVs
	- Vizio
		- https://github.com/0187773933/VizioController
	- LG
		- https://github.com/48723247842/LGTVController
		- https://github.com/snabb/webostv
- VLC
	- Random M3U8 Music Tracks
- Spotify
	- Detect if media control overlay is already open
		- add just generic get status pre-call
			- don't press anything , get 1 screenshot , detect multiple things
				- now playing , which key index , shuffle status , media_controls open?
- Disney+
	- asdf
- YouTube
	- asdf
- Twitch
	- asdf