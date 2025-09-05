/*
 *
 * encrypt.go
 * utils
 *
 * Created by lintao on 2023/11/9 16:33
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// =========================
// 配置 & 全局（线程安全）
// =========================

type Config struct {
	// BcryptCost 建议 10–14；默认 12。根据你的机器压测调整。
	BcryptCost int

	// EnablePrehash: 将明文先做 SHA-256，再喂给 bcrypt，避免 bcrypt 72字节截断问题。
	EnablePrehash bool

	// Pepper：服务端额外机密（可选）。建议用 Base64 编码放在环境变量里再注入。
	// 注意：Pepper 与盐不同；盐由 bcrypt 内部自动处理，Pepper 是全局密钥。
	Pepper []byte
}

const (
	DefaultBcryptCost    = 12
	EnvPepperBase64      = "PASS_PEPPER_B64" // 环境变量名（可自定义）
	maxUAStoreLen        = 256               // 如需记录UA时可参考
	defaultEnablePrehash = true
)

var (
	mu  sync.RWMutex
	cfg = Config{
		BcryptCost:    DefaultBcryptCost,
		EnablePrehash: defaultEnablePrehash,
		Pepper:        nil, // 可通过 SetPepperFromEnv() 或 SetConfig() 注入
	}
)

// SetConfig 全量设置配置（线程安全）。
func SetConfig(c Config) {
	mu.Lock()
	defer mu.Unlock()
	cfg = c
}

// SetBcryptCost 仅调整 cost。
func SetBcryptCost(cost int) {
	mu.Lock()
	defer mu.Unlock()
	cfg.BcryptCost = cost
}

// SetPepper 设置 Pepper。
func SetPepper(p []byte) {
	mu.Lock()
	defer mu.Unlock()
	cfg.Pepper = append([]byte(nil), p...)
}

// SetPepperFromEnv 从环境变量 PASS_PEPPER_B64 读取 Pepper（Base64 编码）。
func SetPepperFromEnv() error {
	if v := strings.TrimSpace(os.Getenv(EnvPepperBase64)); v != "" {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return err
		}
		SetPepper(b)
	}
	return nil
}

// GetConfig 拷贝一份当前配置。
func GetConfig() Config {
	mu.RLock()
	defer mu.RUnlock()
	return cfg
}

// =========================
// 核心API
// =========================

// Hash 生成密码哈希（BCrypt）。默认先做 SHA-256 预哈希，再执行 bcrypt。
// 返回形如 `$2b$12$...` 的哈希字符串。
func Hash(plaintext string) (string, error) {
	c := GetConfig()

	var material []byte
	if c.EnablePrehash {
		material = prehash(plaintext, c.Pepper)
	} else {
		// 若不用预哈希，注意 bcrypt 仅使用前72字节！
		material = []byte(plaintext)
	}

	h, err := bcrypt.GenerateFromPassword(material, c.BcryptCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

// Verify 校验明文是否匹配存量哈希。
// 返回：ok（是否匹配）、needRehash（是否建议用当前策略重算哈希）、err。
func Verify(storedHash, plaintext string) (ok bool, needRehash bool, err error) {
	if storedHash == "" {
		return false, false, errors.New("empty stored hash")
	}
	c := GetConfig()

	// 根据当前策略生成“校验输入”
	material := []byte(plaintext)
	if c.EnablePrehash {
		material = prehash(plaintext, c.Pepper)
	}

	// 比较
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), material); err != nil {
		// 不匹配或其他错误
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, false, nil
		}
		return false, false, err
	}

	// 计算是否需要重哈希（例如 cost 提升，或你未来关闭/开启了预哈希策略等）
	need, _ := NeedsRehash(storedHash, c)
	return true, need, nil
}

// NeedsRehash 判断一个存量哈希是否需要按当前策略重算。
// 当前策略仅依据 Bcrypt cost 判断（因为是否预哈希无法从 bcrypt 字符串反推）。
func NeedsRehash(storedHash string, currentCfg Config) (bool, error) {
	cost, err := bcrypt.Cost([]byte(storedHash))
	if err != nil {
		return false, err
	}
	// 规则：存量cost < 目标cost 则建议重算
	return cost < currentCfg.BcryptCost, nil
}

// MustRehashIfNeeded 在验证通过后调用；若 needRehash 为 true，立即重算并返回新哈希。
// 典型用法：登录成功时检测并平滑升级 cost。
func MustRehashIfNeeded(storedHash, plaintext string) (newHash string, changed bool, err error) {
	ok, need, err := Verify(storedHash, plaintext)
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, errors.New("password mismatch")
	}
	if !need {
		return storedHash, false, nil
	}
	h, err := Hash(plaintext)
	if err != nil {
		return "", false, err
	}
	return h, true, nil
}

// =========================
// 工具函数
// =========================

// 预哈希：SHA-256(password || 0x00 || pepper)
// 这样即便pepper为空也有确定性；并与明文简单拼接拉开域（避免歧义串联）。
func prehash(password string, pepper []byte) []byte {
	h := sha256.New()
	h.Write([]byte(password))
	h.Write([]byte{0})
	if len(pepper) > 0 {
		h.Write(pepper)
	}
	sum := h.Sum(nil)
	return sum
}

// TokenSafeNow：给外部调用者做时间戳（审计时可用）。
func TokenSafeNow() int64 { return time.Now().Unix() }
