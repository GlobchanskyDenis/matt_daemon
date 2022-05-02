package auth

import (
	"matt-daemon/pkg/dto"
	"errors"
)

var gConf *dto.AuthConfigList

func GetConfig() *dto.AuthConfigList {
	if gConf == nil {
		gConf = &dto.AuthConfigList{}
	}
	return gConf
}

func checkConfig() error {
	if gConf == nil {
		return errors.New("Auth package not configured")
	}
	return nil
}

func IsExist(login, password string) bool {
	if gConf == nil {
		return false
	}
	for _, authConf := range *gConf {
		if login == authConf.Login {
			if authConf.PasswordHash == hash(password) {
				return true
			}
		}
	}
	return false
}

func hash(src string) string {
	return src
}
