package internal

import (
	"log"
	"net"
)

// by default, use the standard logger, can be overwritten using kasa.SetLogger(l)
var Klogger kasalogger = log.Default()

// Any log interface that has Println and Printf will do
type kasalogger interface {
	Println(...interface{})
	Printf(string, ...interface{})
}
type Device struct {
	IP     string
	Port   int
	Parsed net.IP
	Debug  bool
}

// KasaDevice is the primary type, defined by kasa devices
type KasaDevice struct {
	System    GetSysinfo `json:"system"`
	Dimmer    Dimmer     `json:"smartlife.iot.dimmer"`
	NetIf     NetIf      `json:"netif"`
	Countdown Countdown  `json:"count_down"`
	Emeter    EmeterSub  `json:"emeter"`
	Schedule  Schedule   `json:"schedule"`
}

// GetSysinfo is defined by kasa devices
type GetSysinfo struct {
	Sysinfo Sysinfo `json:"get_sysinfo"`
}

// Sysinfo is defined by kasa devices
type Sysinfo struct {
	SWVersion      string   `json:"sw_ver"`
	HWVersion      string   `json:"hw_ver"`
	Model          string   `json:"model"`
	DeviceID       string   `json:"deviceId"`
	OEMID          string   `json:"oemId"`
	HWID           string   `json:"hwId"`
	RSSI           int8     `json:"rssi"`
	Longitude      int      `json:"longitude_i"`
	Latitude       int      `json:"latitude_i"`
	Alias          string   `json:"alias"`
	Status         string   `json:"status"`
	MIC            string   `json:"mic_type"`
	Feature        string   `json:"feature"`
	MAC            string   `json:"mac"`
	Updating       uint8    `json:"updating"`
	LEDOff         uint8    `json:"led_off"`
	RelayState     uint8    `json:"relay_state"`
	Brightness     uint8    `json:"brightness"`
	OnTime         int      `json:"on_time"`
	ActiveMode     string   `json:"active_mode"`
	DevName        string   `json:"dev_name"`
	Children       []Child  `json:"children"`
	NumChildren    uint8    `json:"child_num"`
	NTCState       int      `json:"ntc_state"`
	PreferredState []Preset `json:"preferred_state"`
	ErrCode        int8     `json:"error_code"`
}

// Dimmer is defined by kasa devices
type Dimmer struct {
	Parameters DimmerParameters `json:"get_dimmer_parameters"`
	ErrCode    int8             `json:"err_code"`
	ErrMsg     string           `json:"err_msg"`
}

// DimmerParameters is defined by kasa devices
type DimmerParameters struct {
	MinThreshold  uint16 `json:"minThreshold"`
	FadeOnTime    uint16 `json:"fadeOnTime"`
	FadeOffTime   uint16 `json:"fadeOffTime"`
	GentleOnTime  uint16 `json:"gentleOnTime"`
	GentleOffTime uint16 `json:"gentleOffTime"`
	RampRate      uint16 `json:"rampRate"`
	BulbType      uint8  `json:"bulb_type"`
	ErrCode       int8   `json:"err_code"`
	ErrMsg        string `json:"err_msg"`
}

// Child is defined by kasa devices
type Child struct {
	ID         string `json:"id"`
	RelayState uint8  `json:"state"`
	Alias      string `json:"alias"`
	OnTime     int    `json:"on_time"`
	// NextAction
}

// Preset is defined by kasa devices
type Preset struct {
	Index      uint8 `json:"index"`
	Brightness uint8 `json:"brightness"`
}

// NetIf is defined by kasa devices
// {"netif":{"get_stainfo":{"ssid":"IoT8417","key_type":3,"rssi":-61,"err_code":0}}}
type NetIf struct {
	StaInfo StaInfo `json:"get_stainfo"`
	ErrCode int8    `json:"err_code"`
	ErrMsg  string  `json:"err_msg"`
}

// StaInfo is defined by kasa devices
type StaInfo struct {
	SSID    string `json:"ssid"`
	KeyType int8   `json:"key_type"`
	RSSI    int8   `json:"rssi"`
	ErrCode int8   `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

// {"emeter":{"get_realtime":{"current_ma":1799,"voltage_mv":121882,"power_mw":174545,"total_wh":547,"err_code":0}}}
// {"emeter":{"get_daystat":{"day_list":[{"year":2021,"month":2,"day":6,"energy_wh":842},{"year":2021,"month":2,"day":7,"energy_wh":1142}],"err_code":0}}}

// EmeterSub is defined by kasa devices
type EmeterSub struct {
	Realtime EmeterRealtime `json:"get_realtime"`
	DayStat  EmeterDaystat  `json:"get_daystat"`
	ErrCode  int8           `json:"err_code"`
	ErrMsg   string         `json:"err_msg"`
}

// EmeterRealtime is defined by kasa devices
type EmeterRealtime struct {
	Slot      uint8  `json:"slot_id"`
	CurrentMA uint   `json:"current_ma"`
	VoltageMV uint   `json:"voltage_mv"`
	PowerMW   uint   `json:"power_mw"`
	TotalWH   uint   `json:"total_wh"`
	ErrCode   int8   `json:"err_code"`
	ErrMsg    string `json:"err_msg"`
}

// EmeterDaystat is defined by kasa devices
type EmeterDaystat struct {
	List    []EmeterDay `json:"day_list"`
	ErrCode int8        `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
}

// EmeterDay is defined by kasa devices
type EmeterDay struct {
	Year  uint16 `json:"year"`
	Month uint8  `json:"month"`
	Day   uint8  `json:"day"`
	WH    uint16 `json:"energy_wh"`
}

// Countdown is defined by kasa devices
type Countdown struct {
	GetRules GetRules `json:"get_rules"`
	DelRules DelRules `json:"delete_all_rules"`
	AddRule  AddRule  `json:"add_rule"`
}

// GetRules is defined by kasa devices
type GetRules struct {
	RuleList     []Rule `json:"rule_list"`
	ErrorCode    int8   `json:"err_code"`
	ErrorMessage string `json:"err_msg"`
}

// Rule is defined by kasa devices
type Rule struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Enable    uint8  `json:"enable"`
	Delay     uint16 `json:"delay"`
	Active    uint8  `json:"act"`
	Remaining uint16 `json:"remain"`
}

// DelRules is defined by kasa devices
type DelRules struct {
	ErrorCode    int8   `json:"err_code"`
	ErrorMessage string `json:"err_msg"`
}

// AddRule is defined by kasa devices
type AddRule struct {
	ID           string `json:"id"`
	ErrorCode    int8   `json:"err_code"`
	ErrorMessage string `json:"err_msg"`
}

// SetLogger allows applications to register their own logging interface
func SetLogger(l kasalogger) {
	Klogger = l
}

// MonthStat is defined bykasa devices
// {"schedule":{"get_monthstat":{"month_list":[{"year":2024,"month":1,"time":125},{"year":2024,"month":2,"time":1098}],"err_code":0}}}
type Schedule struct {
	GetMonthStat MonthStat `json:"get_monthstat"`
}

type MonthStat struct {
	MonthList []MonthItem `json:"month_list"`
}

type MonthItem struct {
	Year  int16 `json:"year"`
	Month int8  `json:"month"`
	Time  int64 `json:"time"`
}
