# Fire C2 Server

[![Go Reference](https://pkg.go.dev/badge/github.com/0187773933/FireC2Server.svg)](https://pkg.go.dev/github.com/0187773933/FireC2Server)

## Example

- http://localhost:5954/spotify/playlist-shuffle/3UMDmO2YJb8DgUjpSBu8y9?k=asdf

## TODO

- all "fresh" / "wake-up" commands need to first check if we are on the amazon firecube profile selection screen ...
	`adb shell pm list users`
- switching back to redis
	- https://github.com/0187773933/RedisCircularList
- Find some new way to get the config and db reference into logger package.
	- apparently there is init() in go? that gets called before main starts ...
- Audio playing locally through speakers here instead of on streamdeck
- HTML Control Panels
- TVs
	- Vizio
		- https://github.com/0187773933/VizioController
			- done , just need to implement here now
	- LG
		- https://github.com/48723247842/LGTVController
		- https://github.com/snabb/webostv
		- probably have to write an app , or use hdmi cec
	- HISENSE
		- https://github.com/Krazy998/mqtt-hisensetv
		- https://github.com/yosssi/gmq
	- Samsung
		- https://github.com/McKael/samtv
		- https://github.com/mgoff/go-samsung-exlink
		- https://github.com/rainu/samsung-remote-mqtt
		- https://github.com/rainu/samsung-remote
	- HDMICEC
		- https://kwikwai.com/products-2/?show-controller=true
		- https://github.com/torvalds/linux/blob/master/drivers/media/cec/platform/cec-gpio/cec-gpio.c
		- https://github.com/torvalds/linux/blob/master/drivers/media/cec/core/cec-pin.c
		- https://github.com/Kapp4/agocontrol/tree/master/devices
		- https://github.com/MichaelPote/arduino-hdmi-irw
		- https://github.com/s-moch/CEC
		- https://github.com/floe/CEC
- VLC
	- Random M3U8 Music Tracks
- Spotify
	- Detect if media control overlay is already open
		- add just generic get status pre-call
			- don't press anything , get 1 screenshot , detect multiple things
				- now playing , which key index , shuffle status , media_controls open?
- YouTube
	- adhoc playlist support
	- fix s.YouTubeIsVideoIdAvailable( video_id )
		- stream could just be ended , but video_id still valid
			- delete these too
- Disney
	- detect spining circle , stall-out
	- force app reload with same uuid
- Twitch
	- take screenshots of stream ui
	- have to detect where we are in the quality selection menu
	- "weird" quality selection menu trap
- start storing better state details to improve status
- fix sleep ?
	- `adb shell settings get secure sleep_timeout`
	- `adb shell settings put secure sleep_timeout 2147483647`
- PushOver Notifications
- Automatic Staging on HDMI of Twitch Live Streams

## ADB First Time Connection

1. Enter new Docker Container and Run :

	```
	adb connect $FireCubeIP:5555
	failed to authenticate to 192.168.4.193:5555
	```

2. Accept the connection on the Fire Cube TV
3. Confirm Connection :
	`adb devices`
4.
