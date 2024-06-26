package models

type DecodedStateToken struct {
	UserID string
	Code   string
}
type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}
