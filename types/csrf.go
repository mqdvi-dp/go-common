package types

type CsrfToken struct {
	UserUuid string `json:"user_uuid"`
	DeviceId string `json:"device_id"`
}
