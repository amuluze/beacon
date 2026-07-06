// Package jwtauth
// Date: 2024/3/27 15:03
// Author: Amu
// Description:
package jwtauth

import "encoding/json"

type tokenInfo struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (t *tokenInfo) GetAccessToken() string {
	return t.AccessToken
}

func (t *tokenInfo) GetRefreshToken() string {
	return t.RefreshToken
}

func (t *tokenInfo) EncodeToJSON() ([]byte, error) {
	return json.Marshal(t)
}
