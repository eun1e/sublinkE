package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"os"
	"sublink/utils"
	"time"
)

type AccessKey struct {
	ID            int    `gorm:"primaryKey"`
	UserID        int    `gorm:"not null;index"` // 关联到用户的外键
	Username      string `gorm:"not null;index"`
	AccessKeyHash string `gorm:"type:varchar(255);not null;uniqueIndex"` // API Key 哈希值
	CreatedAt     time.Time
	ExpiredAt     *time.Time     `gorm:"index"`             // 过期时间（可选）
	Description   string         `gorm:"type:varchar(255)"` // 备注
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// Generate 保存 AccessKey
func (accessKey *AccessKey) Generate() error {
	return DB.Create(accessKey).Error
}

// FindValidAccessKeys 查找未过期的 AccessKey
func FindValidAccessKeys(userID int) ([]AccessKey, error) {
	var accessKeys []AccessKey
	err := DB.Where("user_id = ?", userID).
		Where("expired_at IS NULL OR expired_at > ?", time.Now()).
		Find(&accessKeys).Error
	return accessKeys, err
}

// Delete 删除 AccessKey (软删除)
func (accessKey *AccessKey) Delete() error {
	return DB.Delete(accessKey).Error
}

// GenerateAPIKey 生成一个新的 API Key,单用户系统直接全随机不编码用户信息
func (accessKey *AccessKey) GenerateAPIKey() (string, error) {
	encryptionKey := os.Getenv("API_ENCRYPTION_KEY")
	if encryptionKey == "" {
		return "", fmt.Errorf("未设置API_ENCRYPTION_KEY环境变量")
	}
	encryptedID, err := utils.EncryptUserIDCompact(accessKey.UserID, []byte(encryptionKey))
	if err != nil {
		log.Println("加密用户ID失败:", err)
		return "", fmt.Errorf("加密用户ID失败: %w", err)
	}
	randomBytes := make([]byte, 18)
	_, err = rand.Read(randomBytes)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("生成随机数据失败: %w", err)
	}

	randomHex := hex.EncodeToString(randomBytes)

	apiKey := fmt.Sprintf("subX_%s_%s", encryptedID, randomHex)

	hashedKey, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("哈希API密钥失败: %w", err)
	}
	accessKey.AccessKeyHash = string(hashedKey)

	return apiKey, nil
}

// VerifyKey 验证提供的 API Key 是否与存储的哈希匹配
func (accessKey *AccessKey) VerifyKey(providedKey string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(accessKey.AccessKeyHash), []byte(providedKey))
	return err == nil
}
