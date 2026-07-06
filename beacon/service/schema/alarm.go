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

type AlarmThresholdQueryReply struct {
	Data []AlarmThreshold `json:"data"`
}

type AlarmThresholdUpdateArgs struct {
	ID        uint   `json:"id" validate:"required"`
	Type      string `json:"type" validate:"lte=64"`
	Duration  int    `json:"duration" validate:"gte=0"`
	Threshold int    `json:"threshold" validate:"gte=0"`
}
