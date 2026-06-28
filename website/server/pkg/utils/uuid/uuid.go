// Package uuid
// Date: 2024/3/27 16:50
// Author: Amu
// Description:
package uuid

import (
	"fmt"
	"unicode"

	"github.com/google/uuid"
)

// UUID Define alias
type UUID = uuid.UUID

// NewUUID Create uuid
func NewUUID() (UUID, error) {
	return uuid.NewRandom()
}

// MustUUID Create uuid(Throw panic if something goes wrong)
func MustUUID() UUID {
	v, err := NewUUID()
	if err != nil {
		panic(err)
	}
	return v
}

// MustString Create uuid
func MustString() string {
	return MustUUID().String()
}

// MustParseUUIString convert uuid str to uuid
func MustParseUUIString(uuidStr string) UUID {
	return uuid.MustParse(uuidStr)
}

// MustID 必须返回以字符为起始的ID
func MustID() string {
	uid := MustString()
	for !unicode.IsLetter(rune(uid[0])) {
		uid = MustString()
	}
	return uid
}

// EncodeUUID converts a uuid byte array to UUID standard string form.
func EncodeUUID(src [16]byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", src[0:4], src[4:6], src[6:8], src[8:10], src[10:16])
}
