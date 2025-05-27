package utils

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
)

// 检查环境
func CheckEnvironment() bool {
	APP_ENV := os.Getenv("APP_ENV")
	if APP_ENV == "" {
		return false
	}
	if strings.Contains(APP_ENV, "development") {
		return true
	}
	return false
}

// Base62字符集
const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// ToBase62 将字节数组转换为Base62字符串
func ToBase62(data []byte) string {
	bi := new(big.Int).SetBytes(data)
	result := ""

	for bi.Sign() > 0 {
		mod := new(big.Int)
		bi.DivMod(bi, big.NewInt(62), mod)
		result = string(base62Chars[mod.Int64()]) + result
	}

	// 确保至少包含一个字符
	if result == "" {
		result = "0"
	}

	return result
}

// FromBase62 将Base62字符串转换回字节数组
func FromBase62(s string) ([]byte, error) {
	bi := new(big.Int)

	for _, c := range s {
		bi.Mul(bi, big.NewInt(62))
		pos := strings.IndexRune(base62Chars, c)
		if pos < 0 {
			return nil, fmt.Errorf("无效字符: %c", c)
		}
		bi.Add(bi, big.NewInt(int64(pos)))
	}

	return bi.Bytes(), nil
}

// EncryptUserIDCompact 使用更紧凑的方法加密用户ID (Base62版本)
func EncryptUserIDCompact(userID int, key []byte) (string, error) {

	// 简单加密方法：将userID与唯一的密钥混合
	// 使用SHA256生成派生密钥
	h := sha256.New()
	h.Write(key)

	// 将userID转换为字节并与派生密钥混合
	userIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(userIDBytes, uint32(userID))

	// 将userID与派生密钥的前4字节进行XOR操作
	derivedKey := h.Sum(nil)[:4]
	encrypted := make([]byte, 4)
	for i := 0; i < 4; i++ {
		encrypted[i] = userIDBytes[i] ^ derivedKey[i]
	}

	// 转换为Base62编码字符串，使结果更短
	return ToBase62(encrypted), nil
}

// DecryptUserIDCompact 解密用户ID (Base62版本)
func DecryptUserIDCompact(encrypted string, key []byte) (int, error) {
	// 使用SHA256生成相同的密钥派生密钥
	h := sha256.New()
	h.Write(key)
	derivedKey := h.Sum(nil)[:4]

	// Base62解码
	encryptedBytes, err := FromBase62(encrypted)
	if err != nil {
		log.Println("Base62解码失败:", err)
		return 0, fmt.Errorf("Base62解码失败: %w", err)
	}

	// 确保长度不超过4字节
	if len(encryptedBytes) > 4 {
		log.Println("解码后数据长度超过4字节:", len(encryptedBytes))
		return 0, fmt.Errorf("解码后数据长度超过4字节: %d", len(encryptedBytes))
	}

	// 如果长度不足4字节，则在前面填充0
	if len(encryptedBytes) < 4 {
		padding := make([]byte, 4-len(encryptedBytes))
		encryptedBytes = append(padding, encryptedBytes...)
	}

	// 与加密使用相同的XOR操作解密
	decrypted := make([]byte, 4)
	for i := 0; i < 4; i++ {
		decrypted[i] = encryptedBytes[i] ^ derivedKey[i]
	}

	// 转换回用户ID，使用uint32确保不会出现负数
	userID := int(binary.BigEndian.Uint32(decrypted))
	return userID, nil
}
