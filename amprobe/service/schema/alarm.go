// Package schema
// Date:   2024/10/14 17:27
// Author: Amu
// Description:
package schema

type AlarmThreshold struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	Duration  int    `json:"duration"`
	Threshold int    `json:"threshold"`
}

type AlarmThresholdUpdateArgs struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	Duration  int    `json:"duration"`
	Threshold int    `json:"threshold"`
}
