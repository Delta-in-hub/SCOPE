package utils

import (
	"bytes"        // 用于比较 byte 切片
	"crypto/rand"  // 用于生成随机密钥
	"encoding/hex" // 用于密钥的 hex 编码/解码
	"testing"
)

// 辅助函数：生成一个随机的 32 字节 AES 密钥并返回其十六进制字符串
// 在测试设置失败时调用 t.Fatalf
func generateValidHexKey(t *testing.T) string {
	t.Helper() // 标记为测试辅助函数
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}
	return hex.EncodeToString(key)
}

func TestEncryptDecrypt(t *testing.T) {
	// 生成一些密钥用于测试
	validKey1 := generateValidHexKey(t)
	validKey2 := generateValidHexKey(t)
	// 确保两个有效密钥不同，用于测试“错误密钥”场景
	if validKey1 == validKey2 {
		t.Fatal("Generated identical keys, cannot run wrong key test reliably")
	}

	tests := []struct {
		name           string
		plaintext      string
		keyInHex       string
		wantEncryptErr bool // 是否期望加密阶段出错

		// --- 以下字段仅在 wantEncryptErr 为 false 时相关 ---
		decryptKeyHex  string              // 用于解密的密钥，如果为空，则使用 keyInHex
		tamperFunc     func([]byte) []byte // 在解密前篡改密文的函数
		wantDecryptErr bool                // 是否期望解密阶段出错
	}{
		// --- 成功场景 ---
		{
			name:           "valid encryption and decryption",
			plaintext:      "hello world, this is a standard test message.",
			keyInHex:       validKey1,
			wantEncryptErr: false,
			decryptKeyHex:  "",  // 使用加密密钥解密
			tamperFunc:     nil, // 不篡改
			wantDecryptErr: false,
		},
		{
			name:           "empty plaintext",
			plaintext:      "",
			keyInHex:       validKey1,
			wantEncryptErr: false,
			decryptKeyHex:  "",
			tamperFunc:     nil,
			wantDecryptErr: false,
		},
		{
			name:           "different valid key",
			plaintext:      "test with another key",
			keyInHex:       validKey2, // 使用第二个有效密钥
			wantEncryptErr: false,
			decryptKeyHex:  "",
			tamperFunc:     nil,
			wantDecryptErr: false,
		},

		// --- 加密失败场景 ---
		{
			name:           "encrypt with invalid key size (short)",
			plaintext:      "test message",
			keyInHex:       "0123456789abcdef", // 16位，太短
			wantEncryptErr: true,
			// 解密部分不适用
		},
		{
			name:           "encrypt with invalid key size (long)",
			plaintext:      "test message",
			keyInHex:       validKey1 + "00", // 33 字节，太长
			wantEncryptErr: true,
			// 解密部分不适用
		},
		{
			name:           "encrypt with invalid hex key characters",
			plaintext:      "test message",
			keyInHex:       "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcXYZ", // 包含无效字符
			wantEncryptErr: true,
			// 解密部分不适用
		},

		// --- 解密失败场景 (加密应成功) ---
		{
			name:           "decrypt with wrong key",
			plaintext:      "secret data",
			keyInHex:       validKey1, // 用 key1 加密
			wantEncryptErr: false,
			decryptKeyHex:  validKey2, // 尝试用 key2 解密
			tamperFunc:     nil,
			wantDecryptErr: true, // 期望解密失败
		},
		{
			name:           "decrypt tampered nonce",
			plaintext:      "sensitive info",
			keyInHex:       validKey1,
			wantEncryptErr: false,
			decryptKeyHex:  "", // 使用正确的密钥
			tamperFunc: func(ct []byte) []byte { // 篡改 Nonce (通常是前 12 字节)
				if len(ct) > 0 {
					ct[0] = ct[0] ^ 0xff // 翻转 Nonce 的第一个字节的位
				}
				return ct
			},
			wantDecryptErr: true, // 期望解密失败 (认证失败)
		},
		{
			name:           "decrypt tampered ciphertext data",
			plaintext:      "more secrets",
			keyInHex:       validKey1,
			wantEncryptErr: false,
			decryptKeyHex:  "", // 使用正确的密钥
			tamperFunc: func(ct []byte) []byte { // 篡改实际密文数据部分
				nonceSize := 12 // GCM 标准 Nonce 大小
				if len(ct) > nonceSize {
					ct[nonceSize] = ct[nonceSize] ^ 0xff // 翻转密文数据的第一个字节
				}
				return ct
			},
			wantDecryptErr: true, // 期望解密失败 (认证失败)
		},
		{
			name:           "decrypt truncated ciphertext (missing data)",
			plaintext:      "cannot be too short",
			keyInHex:       validKey1,
			wantEncryptErr: false,
			decryptKeyHex:  "",
			tamperFunc: func(ct []byte) []byte {
				nonceSize := 12
				if len(ct) > nonceSize+1 { // 确保至少有 Nonce 和 1 字节数据
					return ct[:len(ct)-1] // 移除最后一个字节
				}
				// 如果太短无法截断数据，就返回一个肯定会失败的长度
				if len(ct) > nonceSize {
					return ct[:nonceSize] // 只返回 Nonce
				}
				return ct[:1] // 返回更短的，肯定不够
			},
			wantDecryptErr: true, // 期望解密失败
		},
		{
			name:           "decrypt truncated ciphertext (missing nonce)",
			plaintext:      "needs nonce",
			keyInHex:       validKey1,
			wantEncryptErr: false,
			decryptKeyHex:  "",
			tamperFunc: func(ct []byte) []byte {
				nonceSize := 12
				if len(ct) >= nonceSize {
					return ct[nonceSize:] // 只返回数据部分
				}
				return []byte{} // 如果连 Nonce 都不够长，返回空
			},
			wantDecryptErr: true, // Decrypt 函数内部应检查长度，或 GCM Open 会失败
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// --- 测试加密 ---
			plaintextBytes := []byte(tt.plaintext)
			ciphertext, errEncrypt := Encrypt(plaintextBytes, tt.keyInHex)

			// 检查加密错误是否符合预期
			if tt.wantEncryptErr {
				if errEncrypt == nil {
					t.Errorf("Encrypt() error = nil, wantErr true")
				}
				// 如果期望加密出错，则此测试用例结束
				return
			}
			// 如果不期望加密出错，但出错了
			if !tt.wantEncryptErr && errEncrypt != nil {
				t.Fatalf("Encrypt() unexpected error = %v", errEncrypt) // 使用 Fatalf，因为后续解密无意义
			}

			// --- 如果加密成功，继续测试解密 ---

			// 确定用于解密的密钥
			decKey := tt.keyInHex
			if tt.decryptKeyHex != "" {
				decKey = tt.decryptKeyHex
			}

			// 如果定义了篡改函数，则篡改密文
			if tt.tamperFunc != nil {
				// 注意：直接修改 ciphertext 会影响原始切片，这在这里是期望的
				// 如果不想修改原始的，需要先复制：tamperedCt := append([]byte(nil), ciphertext...)
				ciphertext = tt.tamperFunc(ciphertext)
			}

			// --- 测试解密 ---
			decryptedBytes, errDecrypt := Decrypt(ciphertext, decKey)

			// 检查解密错误是否符合预期
			if tt.wantDecryptErr {
				if errDecrypt == nil {
					t.Errorf("Decrypt() error = nil, wantDecryptErr true")
				}
			} else {
				// 如果不期望解密出错，但出错了
				if errDecrypt != nil {
					t.Errorf("Decrypt() unexpected error = %v", errDecrypt)
				} else {
					// 如果解密成功，比较解密后的内容和原始明文
					if !bytes.Equal(decryptedBytes, plaintextBytes) {
						// 使用 %q 可以清晰显示字符串内容，包括可能的空格等
						t.Errorf("Decrypt() got = %q, want %q", string(decryptedBytes), tt.plaintext)
					}
				}
			}
		})
	}
}
