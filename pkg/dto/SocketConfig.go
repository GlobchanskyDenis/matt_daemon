package dto

import (

)

type SocketConfig struct {
	IsTcpSocket bool   `conf:"IsTcpSocket"`
	Ip          string `conf:"Ip"`
	Port        uint   `conf:"Port"`
}