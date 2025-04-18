package types

import (
	"time"
)

// TokenClaim data for token apps session
type TokenClaim struct {
	UserUuid            string    `json:"user_uuid"`
	Username            string    `json:"username"`
	SubscriptionUuid    string    `json:"subscription_uuid"`
	IsSubscription      bool      `json:"is_subscription"`
	SubscriptionStartAt time.Time `json:"subscription_start_at"`
	SubscriptionEndAt   time.Time `json:"subscription_end_at"`
	RoleId              string    `json:"role_id"`
	RoleKey             string    `json:"role_key"`
	StatusVerify        string    `json:"status_verify"`
	MaxLimitedAccess    time.Time `json:"max_limited_access"`
	DeviceId            string    `json:"device_id"`
	Email               string    `json:"email"`
	Channel             string    `json:"channel"`
	DeviceInfo          string    `json:"device_info"`
	AppVersion          string    `json:"app_version"`
	BusinessName        string    `json:"business_name"`
	ClientId            string    `json:"client_id"`
	ClientName          string    `json:"client_name"`
	ClientKey           string    `json:"client_key"`
	ClientSecret        string    `json:"client_secret"`
	PublicKey           string    `json:"public_key"`
	// identity user
	State                 string `json:"state"`
	City                  string `json:"city"`
	Nik                   string `json:"nik"`
	FullName              string `json:"fullname"`
	PlaceOfBirth          string `json:"place_of_birth"`
	DateOfBirth           string `json:"date_of_birth"`
	Gender                string `json:"gender"`
	Address               string `json:"address"`
	RtRw                  string `json:"rt_rw"`
	AdministrativeVillage string `json:"administrative_village"`
	District              string `json:"district"`
	Religion              string `json:"religion"`
	MaritalStatus         string `json:"marital_status"`
	Occupation            string `json:"occupation"`
	Nationality           string `json:"nationality"`
	BloodType             string `json:"blood_type"`
}

func (tc TokenClaim) GetUserUuid() string               { return tc.UserUuid }
func (tc TokenClaim) GetUsername() string               { return tc.Username }
func (tc TokenClaim) GetSubscriptionUuid() string       { return tc.SubscriptionUuid }
func (tc TokenClaim) GetIsSubscription() bool           { return tc.IsSubscription }
func (tc TokenClaim) GetRoleId() string                 { return tc.RoleId }
func (tc TokenClaim) GetRoleKey() string                { return tc.RoleKey }
func (tc TokenClaim) GetStatus() string                 { return tc.StatusVerify }
func (tc TokenClaim) GetMaxLimitedAccess() time.Time    { return tc.MaxLimitedAccess }
func (tc TokenClaim) GetDeviceId() string               { return tc.DeviceId }
func (tc TokenClaim) GetEmail() string                  { return tc.Email }
func (tc TokenClaim) GetChannel() string                { return tc.Channel }
func (tc TokenClaim) GetDeviceInfo() string             { return tc.DeviceInfo }
func (tc TokenClaim) GetAppVersion() string             { return tc.AppVersion }
func (tc TokenClaim) GetBusinessName() string           { return tc.BusinessName }
func (tc TokenClaim) GetClientId() string               { return tc.ClientId }
func (tc TokenClaim) GetClientName() string             { return tc.ClientName }
func (tc TokenClaim) GetClientKey() string              { return tc.ClientKey }
func (tc TokenClaim) GetClientSecret() string           { return tc.ClientSecret }
func (tc TokenClaim) GetPublicKey() string              { return tc.PublicKey }
func (tc TokenClaim) GetState() string                  { return tc.State }
func (tc TokenClaim) GetCity() string                   { return tc.City }
func (tc TokenClaim) GetNik() string                    { return tc.Nik }
func (tc TokenClaim) GetFullname() string               { return tc.FullName }
func (tc TokenClaim) GetPlaceOfBirth() string           { return tc.PlaceOfBirth }
func (tc TokenClaim) GetDateOfBirth() string            { return tc.DateOfBirth }
func (tc TokenClaim) GetGender() string                 { return tc.Gender }
func (tc TokenClaim) GetAddress() string                { return tc.Address }
func (tc TokenClaim) GetRtRw() string                   { return tc.RtRw }
func (tc TokenClaim) GetAdministrativeVillage() string  { return tc.AdministrativeVillage }
func (tc TokenClaim) GetDistrict() string               { return tc.District }
func (tc TokenClaim) GetReligion() string               { return tc.Religion }
func (tc TokenClaim) GetMaritalStatus() string          { return tc.MaritalStatus }
func (tc TokenClaim) GetOccupation() string             { return tc.Occupation }
func (tc TokenClaim) GetNationality() string            { return tc.Nationality }
func (tc TokenClaim) GetBloodType() string              { return tc.BloodType }
func (tc TokenClaim) GetSubscriptionStartAt() time.Time { return tc.SubscriptionStartAt }
func (tc TokenClaim) GetSubscriptionEndAt() time.Time   { return tc.SubscriptionEndAt }
