package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

func CreateSign(uid []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(uid)
	sign := h.Sum(nil)
	result := uid
	result = append(result, sign...)
	return hex.EncodeToString(result)
}

func DecodeSign(sign string) (uint32, error) {
	data, err := hex.DecodeString(sign)
	if err != nil {
		return 0, err
	}
	ID := binary.BigEndian.Uint32(data[:4])
	return ID, nil
}

func CompareSign(sign string, secret string) (uint32, error) {
	decodedID, err := DecodeSign(sign)
	fmt.Printf("decoded is %d \n", decodedID)
	if err != nil {
		return 0, nil
	}

	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, decodedID)
	tmpSign := CreateSign(size, secret)
	if sign != tmpSign {
		return 0, errors.New("invalid token")
	}

	return decodedID, nil
}
