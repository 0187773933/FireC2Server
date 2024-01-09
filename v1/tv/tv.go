package tv

import (
	"fmt"
	logger "github.com/0187773933/FireC2Server/v1/logger"
	utils "github.com/0187773933/FireC2Server/v1/utils"
	types "github.com/0187773933/FireC2Server/v1/types"
	lg_tv "github.com/0187773933/LGTVController/v1/controller"
	tg_tv_types "github.com/0187773933/LGTVController/v1/types"
	vizio_tv "github.com/0187773933/VizioController/v1/controller"
	hdmi_cec "github.com/0187773933/HDMICEC/v1/controller"
	try "github.com/manucorporat/try"
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
	VIZIO *vizio_tv.Controller `yaml:"-"`
	HDMICEC hdmi_cec.Controller `yaml:"-"`
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
		DefaultVolume: config.TVDefaultVolume ,
		DefaultInput: config.TVDefaultInput ,
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
			result.VIZIO = vizio_tv.New( config.TVIP , config.TVVizioAuthToken )
			break;
		case "HDMICEC":
			result.HDMICEC = hdmi_cec.New()
			break;
	}
	return
}

func ( tv *TV ) Prepare() {
	status := tv.Status()
	if status.Power == false {
		log.Debug( "TV Power Was Off , Turning On" )
		if tv.WakeOnLan == true { utils.WakeOnLan( tv.MAC ) }
		switch tv.Brand{
			case "lg":
				tv.PowerOn()
				break;
			case "samsung":
				tv.PowerOn()
				break;
			case "vizio":
				tv.PowerOn()
			case "HDMICEC":
				try.This( func() {
					tv.HDMICEC.PowerOn()
				}).Catch(func(e try.E) {
					log.Debug( "failed to send hdmi cec power-on command" )
				})
				try.This( func() {
					tv.HDMICEC.SelectHDMI1()
				}).Catch(func(e try.E) {
					log.Debug( "failed to send hdmi cec command to force hdmi-1 awake" )
				})
				break;
		}
	}
	log.Debug( "Turning Mute OFF" )
	tv.MuteOff()
	if status.Volume != tv.DefaultVolume {
		log.Debug( fmt.Sprintf( "Current Volume === %d , Target Volume === %d" , status.Volume , tv.DefaultVolume ) )
		tv.SetVolume( tv.DefaultVolume )
	}
	if status.Input != tv.DefaultInput {
		log.Debug( fmt.Sprintf( "Current Input === %d , Target Input === %d" , status.Input , tv.DefaultInput ) )
		tv.SetInput( tv.DefaultInput )
	}
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
			tv.VIZIO.PowerOn()
		case "HDMICEC":
			tv.HDMICEC.SelectHDMI1()
			break;
	}
}

func ( tv *TV ) PowerOff() {
	switch tv.Brand{
		case "lg":
			tv.LG.API( "power_off" )
			break;
		case "samsung":
			log.Debug( "samsung === PowerOff() === to do" )
			break;
		case "vizio":
			tv.VIZIO.PowerOff()
		case "HDMICEC":
			tv.HDMICEC.PowerOff()
			break;
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
			log.Debug( "samsung === GetPowerStatus() === to do" )
			break;
		case "vizio":
			//vizio_tv.PowerOff( tv.IP , tv.VizioAuthToken )
			tv.VIZIO.PowerGetState()
		case "HDMICEC":
			log.Debug( "HDMICEC === GetPowerStatus() === to do" )
			break;
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
			log.Debug( "samsung === GetInput() === to do" )
			break;
		case "vizio":
			current_input := tv.VIZIO.InputGetCurrent()
			// log.Debug( fmt.Sprintf( "current_input === %s" , current_input.Name ) )
			switch current_input.Name {
				case "hdmi1":
					result = 1
					break;
				case "hdmi2":
					result = 2
					break;
				case "hdmi3":
					result = 3
					break;
				case "hdmi4":
					result = 4
					break;
				default:
					result = -1
					break
			}
			break;
		case "HDMICEC":
			log.Debug( "HDMICEC === GetInput() === to do" )
			break;
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
			log.Debug( "samsung === SetInput() === to do" )
			break;
		case "vizio":
			target_input := fmt.Sprintf( "HDMI-%d" , hdmi_input )
			fmt.Println( "setting vizio to :" , target_input )
			tv.VIZIO.InputSet( target_input )
		case "HDMICEC":
			log.Debug( "HDMICEC === SetInput() === to do" )
			break;
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
			log.Debug( "samsung === MuteOn() === to do" )
			break;
		case "vizio":
			tv.VIZIO.MuteOn()
		case "HDMICEC":
			log.Debug( "HDMICEC === MuteOn() === to do" )
			break;
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
			log.Debug( "samsung === MuteOff() === to do" )
			break;
		case "vizio":
			tv.VIZIO.MuteOff()
		case "HDMICEC":
			log.Debug( "HDMICEC === MuteOff() === to do" )
			break;
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
			log.Debug( "samsung === GetVolume() === to do" )
			break;
		case "vizio":
			result = tv.VIZIO.VolumeGet()
		case "HDMICEC":
			log.Debug( "HDMICEC === GetVolume() === to do" )
			break;
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
			log.Debug( "samsung === GetVolume() === to do" )
			break;
		case "vizio":
			tv.VIZIO.VolumeSet( volume_level )
		case "HDMICEC":
			log.Debug( "HDMICEC === SetVolume() === to do" )
			break;
	}
}
