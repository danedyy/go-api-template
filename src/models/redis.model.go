package models

var RedisKeys = struct {
	DataAuthStateTokens     string
	AccessToken             string
	RefreshToken            string
	ConfirmEmail            string
	PasswordReset           string
	ConfirmBeneficiaryEmail string
	PasswordRetries         string
	ValidateBvn             string
}{
	DataAuthStateTokens:     "data:auth:state-tokens",
	AccessToken:             "auth:user:access:token",
	RefreshToken:            "auth:user:refresh:token",
	ConfirmEmail:            "auth:user:confirm:email:token",
	PasswordReset:           "auth:user:password-reset:email:token",
	ConfirmBeneficiaryEmail: "beneficiary:confirm:email",
	PasswordRetries:         "password:retries",
	ValidateBvn:             "bvn:validate",
}
