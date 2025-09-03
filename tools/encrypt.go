/*
 *
 * encrypt.go
 * tools
 *
 * Created by lintao on 2023/11/9 16:33
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package tools

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha25(text string) string {
	// 创建一个新的SHA-256哈希器实例
	hasher := sha256.New()

	// 写入数据到哈希器
	hasher.Write([]byte(text))

	// 计算最终的散列值
	// Sum方法会返回一个新的slice，其中包含了哈希的结果
	hashed := hasher.Sum(nil)

	// 将散列值转换为16进制字符串
	hashedStr := hex.EncodeToString(hashed)
	return hashedStr
}
