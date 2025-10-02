package main

import (
	"fmt"
	"log"

	"github.com/koinvote/btcrpc/scenario_tests/shared"
)

// TestConfigValidation 驗證統一配置管理的正確性
func TestConfigValidation() {
	fmt.Println("=== 測試統一配置管理 ===")

	// 測試預設配置
	fmt.Println("\n🔸 測試預設配置...")
	
	config := shared.DefaultConfig()
	fmt.Printf("✓ URL: %s\n", config.URL)
	fmt.Printf("✓ Username: %s\n", config.Username)
	fmt.Printf("✓ Password: %s\n", config.Password)

	// 測試客戶端創建
	fmt.Println("\n🔸 測試客戶端創建...")
	
	client := shared.NewTestClient()
	if client == nil {
		log.Fatal("❌ 無法創建測試客戶端")
	}
	fmt.Println("✓ 成功創建測試客戶端")

	// 測試自定義配置
	fmt.Println("\n🔸 測試自定義配置...")
	
	customConfig := &shared.TestConfig{
		URL:      "http://custom-bitcoin:18443",
		Username: "custom_user",
		Password: "custom_pass",
	}
	
	customClient := shared.NewTestClientWithConfig(customConfig)
	if customClient == nil {
		log.Fatal("❌ 無法創建自定義配置客戶端")
	}
	fmt.Println("✓ 成功創建自定義配置客戶端")

	fmt.Println("\n✅ 統一配置管理測試通過")
}

func main() {
	TestConfigValidation()
}