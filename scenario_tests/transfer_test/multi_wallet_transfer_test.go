package main

import (
	"fmt"
	"log"
	"time"

	"github.com/koinvote/btcrpc"
)

// MultiWalletTransferScenarioTest 測試多錢包間轉帳情境
// 這個測試模擬以下場景：
// 1. 創建兩個錢包 (Alice 和 Bob)
// 2. 為每個錢包生成地址
// 3. 挖礦給 Alice 錢包獲取初始資金
// 4. Alice 轉帳給 Bob
// 5. 驗證轉帳後的餘額變化
// 6. Bob 再轉帳回給 Alice
// 7. 最終驗證所有餘額
func MultiWalletTransferScenarioTest() {
	// Bitcoin Core RPC 客戶端配置
	client := btcrpc.NewClient(
		"http://bitcoin-core:18443", // regtest 模式
		"bitcoinrpc",
		"test_password",
	)

	fmt.Println("=== 多錢包轉帳情境測試 ===")

	// Step 1: 創建兩個錢包
	fmt.Println("\n🔸 步驟 1: 創建錢包...")

	aliceWallet := "alice_wallet"
	bobWallet := "bob_wallet"

	// 創建 Alice 的錢包
	_, err := client.CreateWallet(aliceWallet, false, false, "", false)
	if err != nil {
		log.Printf("Warning: Alice 錢包可能已存在: %v", err)
	} else {
		fmt.Printf("✓ 成功創建 %s\n", aliceWallet)
	}

	// 創建 Bob 的錢包
	_, err = client.CreateWallet(bobWallet, false, false, "", false)
	if err != nil {
		log.Printf("Warning: Bob 錢包可能已存在: %v", err)
	} else {
		fmt.Printf("✓ 成功創建 %s\n", bobWallet)
	}

	// Step 2: 為每個錢包生成地址
	fmt.Println("\n🔸 步驟 2: 生成錢包地址...")

	aliceAddr, err := client.GetNewAddress(aliceWallet, "alice_main", "bech32")
	if err != nil {
		log.Fatalf("無法為 Alice 生成地址: %v", err)
	}
	fmt.Printf("✓ Alice 地址: %s\n", aliceAddr)

	bobAddr, err := client.GetNewAddress(bobWallet, "bob_main", "bech32")
	if err != nil {
		log.Fatalf("無法為 Bob 生成地址: %v", err)
	}
	fmt.Printf("✓ Bob 地址: %s\n", bobAddr)

	// Step 3: 挖礦給 Alice 獲取初始資金
	fmt.Println("\n🔸 步驟 3: 挖礦獲取初始資金...")

	// 挖 101 個區塊給 Alice (需要超過 100 個區塊才能花費挖礦獎勵)
	blockHashes, err := client.GenerateToAddress(101, aliceAddr, nil)
	if err != nil {
		log.Fatalf("挖礦失敗: %v", err)
	}
	fmt.Printf("✓ 成功挖掘 %d 個區塊\n", len(blockHashes))

	// 等待區塊確認
	time.Sleep(2 * time.Second)

	// 檢查 Alice 的初始餘額
	aliceBalance, err := client.GetBalance(aliceWallet, nil, nil)
	if err != nil {
		log.Fatalf("無法獲取 Alice 餘額: %v", err)
	}
	fmt.Printf("✓ Alice 初始餘額: %.8f BTC\n", aliceBalance)

	if aliceBalance == 0 {
		log.Fatalf("Alice 餘額為 0，無法進行轉帳測試")
	}

	// Step 4: Alice 轉帳給 Bob
	fmt.Println("\n🔸 步驟 4: Alice 轉帳給 Bob...")

	transferAmount := 10.0 // 轉帳 10 BTC

	txid1, err := client.SendToAddressSimple(aliceWallet, bobAddr, transferAmount)
	if err != nil {
		log.Fatalf("Alice 轉帳給 Bob 失敗: %v", err)
	}
	fmt.Printf("✓ 轉帳交易 ID: %s\n", txid1)

	// 挖一個區塊確認交易
	_, err = client.GenerateToAddress(1, aliceAddr, nil)
	if err != nil {
		log.Fatalf("確認交易失敗: %v", err)
	}

	// 等待交易確認
	time.Sleep(2 * time.Second)

	// Step 5: 驗證轉帳後的餘額
	fmt.Println("\n🔸 步驟 5: 驗證第一次轉帳後的餘額...")

	aliceBalanceAfter1, err := client.GetBalance(aliceWallet, nil, nil)
	if err != nil {
		log.Fatalf("無法獲取 Alice 轉帳後餘額: %v", err)
	}

	bobBalanceAfter1, err := client.GetBalance(bobWallet, nil, nil)
	if err != nil {
		log.Fatalf("無法獲取 Bob 轉帳後餘額: %v", err)
	}

	fmt.Printf("✓ Alice 轉帳後餘額: %.8f BTC\n", aliceBalanceAfter1)
	fmt.Printf("✓ Bob 轉帳後餘額: %.8f BTC\n", bobBalanceAfter1)

	// 驗證 Bob 確實收到了轉帳
	if bobBalanceAfter1 != transferAmount {
		log.Fatalf("Bob 餘額不正確，期望: %.8f，實際: %.8f", transferAmount, bobBalanceAfter1)
	}

	// Step 6: Bob 轉帳回給 Alice
	fmt.Println("\n🔸 步驟 6: Bob 轉帳回給 Alice...")

	returnAmount := 5.0 // 轉回 5 BTC

	txid2, err := client.SendToAddressSimple(bobWallet, aliceAddr, returnAmount)
	if err != nil {
		log.Fatalf("Bob 轉帳給 Alice 失敗: %v", err)
	}
	fmt.Printf("✓ 回轉交易 ID: %s\n", txid2)

	// 挖一個區塊確認交易
	_, err = client.GenerateToAddress(1, aliceAddr, nil)
	if err != nil {
		log.Fatalf("確認回轉交易失敗: %v", err)
	}

	// 等待交易確認
	time.Sleep(2 * time.Second)

	// Step 7: 最終驗證所有餘額
	fmt.Println("\n🔸 步驟 7: 最終餘額驗證...")

	aliceFinalBalance, err := client.GetBalance(aliceWallet, nil, nil)
	if err != nil {
		log.Fatalf("無法獲取 Alice 最終餘額: %v", err)
	}

	bobFinalBalance, err := client.GetBalance(bobWallet, nil, nil)
	if err != nil {
		log.Fatalf("無法獲取 Bob 最終餘額: %v", err)
	}

	fmt.Printf("✓ Alice 最終餘額: %.8f BTC\n", aliceFinalBalance)
	fmt.Printf("✓ Bob 最終餘額: %.8f BTC\n", bobFinalBalance)

	// 驗證餘額邏輯
	expectedBobBalance := transferAmount - returnAmount // 10 - 5 = 5
	if bobFinalBalance != expectedBobBalance {
		log.Fatalf("Bob 最終餘額不正確，期望: %.8f，實際: %.8f", expectedBobBalance, bobFinalBalance)
	}

	// 檢查交易記錄
	fmt.Println("\n🔸 步驟 8: 驗證交易記錄...")

	aliceTransactions, err := client.ListTransactions(aliceWallet, "*", 10, 0, false)
	if err != nil {
		log.Fatalf("無法獲取 Alice 交易記錄: %v", err)
	}

	bobTransactions, err := client.ListTransactions(bobWallet, "*", 10, 0, false)
	if err != nil {
		log.Fatalf("無法獲取 Bob 交易記錄: %v", err)
	}

	fmt.Printf("✓ Alice 有 %d 筆交易記錄\n", len(aliceTransactions))
	fmt.Printf("✓ Bob 有 %d 筆交易記錄\n", len(bobTransactions))

	// 顯示最近的交易
	if len(aliceTransactions) > 0 {
		lastTx := aliceTransactions[0]
		fmt.Printf("   Alice 最近交易: %s, 金額: %.8f, 類型: %s\n",
			lastTx.TxID[:16]+"...", lastTx.Amount, lastTx.Category)
	}

	if len(bobTransactions) > 0 {
		lastTx := bobTransactions[0]
		fmt.Printf("   Bob 最近交易: %s, 金額: %.8f, 類型: %s\n",
			lastTx.TxID[:16]+"...", lastTx.Amount, lastTx.Category)
	}

	fmt.Println("\n✅ 多錢包轉帳情境測試完成！")
	fmt.Printf("📊 測試結果總結:\n")
	fmt.Printf("   - Alice 從 %.8f BTC → %.8f BTC\n", aliceBalance, aliceFinalBalance)
	fmt.Printf("   - Bob 從 0 BTC → %.8f BTC\n", bobFinalBalance)
	fmt.Printf("   - 總計執行了 2 筆轉帳交易\n")
	fmt.Printf("   - 所有餘額驗證通過 ✓\n")
}

func main() {
	MultiWalletTransferScenarioTest()
}
