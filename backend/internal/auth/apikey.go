package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashAPIKey 与管理密钥校验使用同一摘要算法
func HashAPIKey(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
