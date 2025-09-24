package jwt

import sdk "github.com/jack0829/enmooy-account-sdk/jwt"

var (
	Auth *sdk.JWT
)

func InitAuth(key string, salt string) {
	Auth = sdk.New([]byte(key), salt)
}
