package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sublink/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AccessKey struct {
	ID            int        `gorm:"primaryKey" json:"id"`
	UserID        int        `gorm:"not null;index" json:"userID"` // 关联到用户的外键
	Username      string     `gorm:"not null;index" json:"username"`
	AccessKeyHash string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"-"` // API Key 哈希值，不返回给前端
	CreatedAt     time.Time  `gorm:"" json:"createdAt"`
	ExpiredAt     *time.Time `gorm:"index" json:"expiredAt"`               // 过期时间（可选）
	Description   string     `gorm:"type:varchar(255)" json:"description"` // 备注
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

// Delete 删除 AccessKey (物理删除)
func (accessKey *AccessKey) Delete() error {
	// 先从数据库获取完整的 AccessKey 信息（包括 Username）
	var fullAccessKey AccessKey
	err := DB.First(&fullAccessKey, accessKey.ID).Error
	if err != nil {
		return fmt.Errorf("获取 AccessKey 信息失败: %w", err)
	}

	// 删除数据库记录
	err = DB.Unscoped().Delete(accessKey).Error
	if err != nil {
		return err
	}

	return nil
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

	apiKey := fmt.Sprintf("subE_%s_%s", encryptedID, randomHex)

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

// CleanupExpiredAccessKeys 清理过期的 AccessKey
func CleanupExpiredAccessKeys() error {
	log.Println("开始清理过期的 AccessKey")

	// 查找所有过期的 AccessKey
	var expiredKeys []AccessKey
	err := DB.Where("expired_at IS NOT NULL AND expired_at < ?", time.Now()).Find(&expiredKeys).Error
	if err != nil {
		log.Printf("查询过期 AccessKey 失败: %v", err)
		return fmt.Errorf("查询过期 AccessKey 失败: %w", err)
	}

	if len(expiredKeys) == 0 {
		log.Println("没有发现过期的 AccessKey")
		return nil
	}

	log.Printf("发现 %d 个过期的 AccessKey，准备清理", len(expiredKeys))

	// 批量删除过期的 AccessKey
	for _, key := range expiredKeys {
		err := key.Delete()
		if err != nil {
			log.Printf("删除过期 AccessKey 失败，ID: %d, 错误: %v", key.ID, err)
			continue
		}
		log.Printf("成功删除过期 AccessKey，ID: %d, Username: %s", key.ID, key.Username)
	}

	log.Printf("过期 AccessKey 清理完成，共处理 %d 个", len(expiredKeys))
	return nil
}

// StartAccessKeyCleanupScheduler 启动 AccessKey 清理定时任务
func StartAccessKeyCleanupScheduler() {
	log.Println("启动 AccessKey 清理定时任务")
	// 每小时执行一次清理
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		defer ticker.Stop()
		for range ticker.C {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("AccessKey 清理任务异常: %v", r)
					}
				}()
				CleanupExpiredAccessKeys()
			}()
		}
	}()

	// 启动时立即执行一次清理
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("AccessKey 初始清理任务异常: %v", r)
			}
		}()
		CleanupExpiredAccessKeys()
	}()
}
