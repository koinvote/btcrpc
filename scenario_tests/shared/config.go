package shared

import "github.com/koinvote/btcrpc"

// 重新導出常用的類型，以便測試可以使用
type CreateRawTransactionInput = btcrpc.CreateRawTransactionInput
type CreateRawTransactionOutput = btcrpc.CreateRawTransactionOutput
type UTXO = btcrpc.UTXO

// TestConfig 包含所有測試的統一配置
type TestConfig struct {
	URL      string
	Username string
	Password string
}

// DefaultConfig 返回預設的測試配置
func DefaultConfig() *TestConfig {
	return &TestConfig{
		URL:      "http://bitcoin-core-regtest:18443",
		Username: "bitcoinrpc",
		Password: "test_password",
	}
}

// NewTestClient 創建使用預設配置的測試客戶端
func NewTestClient() *btcrpc.Client {
	config := DefaultConfig()
	return btcrpc.NewClient(config.URL, config.Username, config.Password)
}

// NewTestClientWithConfig 創建使用自定義配置的測試客戶端
func NewTestClientWithConfig(config *TestConfig) *btcrpc.Client {
	return btcrpc.NewClient(config.URL, config.Username, config.Password)
}
