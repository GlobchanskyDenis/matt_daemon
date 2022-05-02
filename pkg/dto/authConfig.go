package dto

import ()

type AuthConfigList []AuthConfig

type AuthConfig struct {
	Login        string `conf:"Login"`
	PasswordHash string `conf:"PasswordHash"`
}
