package title

import (
	"errors"
	"strings"
	"unicode/utf8"

	"beacon/service/model"
)

const MaxLength = 128

var (
	ErrEmpty     = errors.New("告警消息标题不能为空")
	ErrLineBreak = errors.New("告警消息标题不能包含换行")
	ErrTooLong   = errors.New("告警消息标题不能超过 128 个字符")
)

// Normalize applies the shared title invariant used by persistence and every
// notification channel.
func Normalize(raw string) (string, error) {
	if strings.ContainsAny(raw, "\r\n") {
		return "", ErrLineBreak
	}
	title := strings.TrimSpace(raw)
	if title == "" {
		return "", ErrEmpty
	}
	if utf8.RuneCountInString(title) > MaxLength {
		return "", ErrTooLong
	}
	return title, nil
}

// OrDefault keeps notification delivery available if a legacy or damaged
// stored value violates the current invariant.
func OrDefault(raw string) string {
	title, err := Normalize(raw)
	if err != nil {
		return model.DefaultAlarmMessageTitle
	}
	return title
}
