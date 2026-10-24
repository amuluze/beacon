// Package jwtauth
// Date: 2024/3/27 15:03
// Author: Amu
// Description:
package jwtauth

type tokenInfo struct {
	accessToken  string
	refreshToken string
}

func (t *tokenInfo) GetAccessToken() string {
	return t.accessToken
}

func (t *tokenInfo) GetRefreshToken() string {
	return t.refreshToken
}
