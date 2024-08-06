package sessions

import (
	"log"

	"aidanwoods.dev/go-paseto"
	"github.com/risersh/rest/conf"
)

var PrivateKey paseto.V4AsymmetricSecretKey
var PublicKey paseto.V4AsymmetricPublicKey

func InitKeys() {
	var err error
	log.Printf("private key: %s", conf.Config.BaseConfig.Sessions.PrivateKey)
	PrivateKey, err = paseto.NewV4AsymmetricSecretKeyFromHex(conf.Config.Sessions.PrivateKey)
	if err != nil {
		panic(err)
	}
	PublicKey, err = paseto.NewV4AsymmetricPublicKeyFromHex(conf.Config.Sessions.PublicKey)
	if err != nil {
		panic(err)
	}
}
