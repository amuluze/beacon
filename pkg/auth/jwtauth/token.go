// Package jwtauth
// Date: 2024/3/27 15:03
// Author: Amu
// Description:
package jwtauth

import "encoding/json"

type tokenInfo struct {
	Token     string `json:"token"`      // 访问令牌
	TokenType string `json:"token_type"` // 令牌类型
	ExpiresAt string `json:"expires_at"` // 过期时间
}

func (t *tokenInfo) GetAccessToken() string {
	return t.Token
}

func (t *tokenInfo) GetTokenType() string {
	return t.TokenType
}

func (t *tokenInfo) GetExpiresAt() string {
	return t.ExpiresAt
}

func (t *tokenInfo) EncodeToJSON() ([]byte, error) {
	return json.Marshal(t)
}
