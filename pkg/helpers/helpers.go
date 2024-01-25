package helpers

import (
	"coffee-shop/pkg/custom_errors"
	"encoding/base64"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GetUTCTimeStamp() time.Time {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	return now
}

func GetRandomTxnPin() int {
	return rangeIn(100000, 999999)
}

func rangeIn(low, hi int) int {
	rand.Seed(time.Now().UnixNano())
	return low + rand.Intn(hi-low)
}

func DecodeCursor(encodedCursor string) (res time.Time, err error) {
	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return
	}
	res, err = time.Parse(time.RFC3339Nano, string(byt))
	if err != nil {
		return
	}
	return
}

func EncodeCursor(t time.Time) string {
	key := t.Format(time.RFC3339Nano)
	return base64.StdEncoding.EncodeToString([]byte(key))
}

func HashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePasswords(hashedPwd string, plainPwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		return false, custom_errors.UIError.New(nil, false, "Invalid password!")
	}
	return true, err
}

func Ptr[T any](x T) *T {
	return &x
}
