package tv

import (
	"fmt"
	logger "github.com/0187773933/FireC2Server/v1/logger"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	types "github.com/0187773933/FireC2Server/v1/types"
	lg_tv "github.com/0187773933/LGTVController/v1/controller"
	tg_tv_types "github.com/0187773933/LGTVController/v1/types"
	vizio_tv "github.com/0187773933/VizioController/controller"
)

var log = logger.GetLogger()

type TV struct {
	WakeOnLan bool `yaml:"tv_wake_on_lan"`
	Brand string `yaml:"tv_brand"`
	IP string `yaml:"tv_ip"`
	MAC string `yaml:"tv_mac"`
	DefaultVolume int `yaml:"tv_default_volume"`
	DefaultInput int `yaml:"tv_default_input"`
	WebSocketPort string `yaml:"websocket_port"`
	LGClientKey string `yaml:"lg_client_key"`
	VizioAuthToken string `yaml:"vizio_auth_token"`
	TimeoutSeconds int `yaml:"timeout_seconds"`
	LG *lg_tv.Controller `yaml:"-"`
}

// TODO: Samsung TV
// http://wiki.samygo.tv/index.php?title=Main_Page
// https://github.com/ninjasphere/go-samsung-tv
// https://github.com/YuukanOO/rtv/blob/master/main.go

func New( config *types.ConfigFile ) ( result *TV ) {
	result = &TV{
		Brand: config.TVBrand ,
		IP: config.TVIP ,
		MAC: config.TVMAC ,
		WebSocketPort: config.TVWebSocketPort ,
		LGClientKey: config.TVLGClientKey ,
		VizioAuthToken: config.TVVizioAuthToken ,
		TimeoutSeconds: config.TVTimeoutSeconds ,
	}
	// Initial Setup
	switch config.TVBrand{
		case "lg":
			lg_tv_config := &tg_tv_types.ConfigFile{
				TVIP: config.TVIP ,
				TVMAC: config.TVMAC ,
				WebSocketPort: config.TVWebSocketPort ,
				ClientKey: config.TVLGClientKey ,
				TimeoutSeconds: config.TVTimeoutSeconds ,
			}
			result.LG = lg_tv.New( lg_tv_config )
			break;
		case "samsung":
			log.Debug( "samsung === todo" )
			break;
		case "vizio":
			log.Debug( "github.com/0187773933/VizioController/controller v0 doesn't require setup" )
			break;
	}
	return
}

func ( tv *TV ) Prepare() {
	status := tv.Status()
	if status.Power == false {
		if tv.WakeOnLan == true { utils.WakeOnLan( tv.MAC ) }
		tv.PowerOn()
	}
	tv.MuteOff()
	if status.Volume != tv.DefaultVolume { tv.SetVolume( tv.DefaultVolume ) }
	if status.Input != tv.DefaultInput { tv.SetInput( tv.DefaultInput ) }
}

type Status struct {
	Volume int `json:"volume"`
	Input int `json:"input"`
	Power bool `json:"power"`
}
func ( tv *TV ) Status() ( result Status ) {
	log.Debug( "TV.Status()" )
	result.Volume = tv.GetVolume()
	log.Debug( "Volume === " , result.Volume )
	result.Input = tv.GetInput()
	log.Debug( "Input === " , result.Input )
	result.Power = tv.GetPowerStatus()
	log.Debug( "Power === " , result.Power )
	// result.Mute = tv.GetPowerStatus()
	return;
}

func ( tv *TV ) PowerOn() {
	switch tv.Brand{
		case "lg":
			tv.LG.API( "power_on" )
			break;
		case "samsung":
			break;
		case "vizio":
			vizio_tv.PowerOn( tv.IP , tv.VizioAuthToken )
	}
}

func ( tv *TV ) PowerOff() {
	switch tv.Brand{
		case "lg":
			tv.LG.API( "power_off" )
			break;
		case "samsung":
			break;
		case "vizio":
			vizio_tv.PowerOff( tv.IP , tv.VizioAuthToken )
	}
}

func ( tv *TV ) GetPowerStatus() ( result bool ) {
	result = false
	switch tv.Brand{
		case "lg":
			read_result := tv.LG.API( "get_volume" )
			if read_result == "error reading message" {
				result = true
			} else if read_result == "timeout while reading message" {
				result = true
			}
			break;
		case "samsung":
			break;
		case "vizio":
			//vizio_tv.PowerOff( tv.IP , tv.VizioAuthToken )
			log.Debug( "vizio === GetPowerStatus() === to do" )
	}
	return;
}

func ( tv *TV ) GetInput() ( result int ) {
	log.Debug( "TV.GetInput()" )
	switch tv.Brand{
		case "lg":
			result_string := tv.LG.API( "get_inputs" )
			log.Debug( "lg === to do , unknown what get_inputs list is" , result_string )
			// result = utils.StringToInt( result_string )
			// log.Println( "LG-3" )
			result = 1
			break;
		case "samsung":
			break;
		case "vizio":
			v_result := vizio_tv.GetCurrentInput( tv.IP , tv.VizioAuthToken )
			log.Debug( "vizio === to do , probably have to split strings" , v_result , v_result.Name )
			result = 1
	}
	log.Debug( "done" )
	return;
}

func ( tv *TV ) SetInput( hdmi_input int ) {
	switch tv.Brand{
		case "lg":
			tv.LG.API( "set_input" , tg_tv_types.Payload{
				"inputId": fmt.Sprintf( "HDMI-%d" , hdmi_input ) ,
			})
			break;
		case "samsung":
			break;
		case "vizio":
			vizio_tv.SetInput( tv.IP , tv.VizioAuthToken , fmt.Sprintf( "HDMI-%d" , hdmi_input ) )
	}
}

func ( tv *TV ) MuteOn() {
	switch tv.Brand{
		case "lg":
			tv.LG.API( "set_mute" , tg_tv_types.Payload{
				"mute": true ,
			})
			break;
		case "samsung":
			break;
		case "vizio":
			vizio_tv.MuteOn( tv.IP , tv.VizioAuthToken )
	}
}

func ( tv *TV ) MuteOff() {
	switch tv.Brand{
		case "lg":
			tv.LG.API( "set_mute" , tg_tv_types.Payload{
				"mute": false ,
			})
			break;
		case "samsung":
			break;
		case "vizio":
			vizio_tv.MuteOff( tv.IP , tv.VizioAuthToken )
	}
}

// func ( tv *TV ) GetMute() {
// 	switch tv.Brand{
// 		case "lg":
// 			audio_status := tv.LG.API( "get_audio_status" )
// 			log.Println( "lg === to do" , audio_status )
// 			break;
// 		case "samsung":
// 			break;
// 		case "vizio":
// 			vizio_tv.MuteOn( tv.IP , tv.VizioAuthToken )
// 	}
// }


func ( tv *TV ) GetVolume() ( result int ) {
	result = -1
	switch tv.Brand{
		case "lg":
			result_string := tv.LG.API( "get_volume" )
			result = utils.StringToInt( result_string )
			break;
		case "samsung":
			break;
		case "vizio":
			result = vizio_tv.GetVolume( tv.IP , tv.VizioAuthToken )
	}
	return;
}

func ( tv *TV ) SetVolume( volume_level int ) {
	switch tv.Brand{
		case "lg":
			tv.LG.API( "set_input" , tg_tv_types.Payload{
				"volume": volume_level ,
			})
			break;
		case "samsung":
			break;
		case "vizio":
			// current_volume := tv.GetVolume()
			// vizio_tv.VolumeUp( tv.IP , tv.VizioAuthToken )
			// vizio_tv.VolumeDown( tv.IP , tv.VizioAuthToken )
			log.Debug( "TODO" )
	}
}
