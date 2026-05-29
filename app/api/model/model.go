package model

import "time"

type LimiterInfo struct {
	UserID   int64
	DeviceID string

	LastTime   int64
	LastTokens int64
	ExpireTime time.Duration
}

type TokenInfo struct {
	RefreshTokenID string
	UserID         int64
	DeviceID       string
	ExpireTime     time.Duration
}
