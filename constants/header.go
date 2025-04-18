package constants

const (
	// Authorization header
	Authorization string = "Authorization"
	// ContentType header
	ContentType string = "Content-Type"
	// ApiKey header
	ApiKey string = "API_KEY"
	// ApplicationApiKey header
	AppliacationApiKey string = "x-api-key"
	// ApplicationVersion header
	ApplicationVersion string = "x-app-version"
	// ApplicationChannel header
	ApplicationChannel string = "x-channel"
	// ApplicationDevice header
	ApplicationDevice string = "x-device-id"
	// ApplicationTimestamp header
	ApplicationTimestamp string = "x-timestamp"
	// ApplicationSignature header
	ApplicationSignature string = "x-signature"
	// ApplicationSkipSignature header
	ApplicationSkipSignature string = "x-skip-signature"
	// ApplicationOriginalIp header
	ApplicationOriginalIp string = "x-original-forwarded-for"
	// ApplicationLatitude header
	ApplicationLatitude string = "x-Latitude-id"
	// ApplicationLatitude header
	ApplicationLongitude string = "x-Longitude-id"
	// ApplicationLatitude header
	ApplicationMockLocation string = "x-Mock-id"
	// ApplicationTimezone header
	ApplicationTimezone string = "x-timezone"
	// ApplicationTradeUuid for ref id transaction to provider
	ApplicationTradeUuid string = "x-trade-uuid"
	// ApplicationClientKey header
	ApplicationClientKey string = "x-client-key"
	// ApplicationClientSecret header
	ApplicationClientSecret string = "x-client-secret"

	// Information profile
	ApplicationUserId             string = "x-user-uuid"
	ApplicationUsername           string = "x-username"
	ApplicationIdentityType       string = "x-identity-type"
	ApplicationIdentityCountry    string = "x-identity-country"
	ApplicationIdentityProvince   string = "x-identity-province"
	ApplicationIdentityCity       string = "x-identity-city"
	ApplicationIdentityPostalCode string = "x-identity-postal-code"
	ApplicationIdentityNumber     string = "x-identity-number"
	ApplicationIdentityName       string = "x-identity-name"
	ApplicationIdentityPoB        string = "x-identity-pob"
	ApplicationIdentityDoB        string = "x-identity-dob"
	ApplicationIdentityAddress    string = "x-identity-address"

	// Bearer header value authorization
	Bearer string = "Bearer"
	// Basic header value authorization
	Basic string = "Basic"
	// ApplicationJson header value of content-type
	ApplicationJson string = "application/json"
	// ApplicationUrlEncoded header value of content-type
	ApplicationUrlEncoded string = "application/x-www-form-urlencoded"
)
