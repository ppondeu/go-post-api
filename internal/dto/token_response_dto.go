package dto

type TokenResponseDto struct {
	AccessToken string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
