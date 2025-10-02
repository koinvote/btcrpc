package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/koinvote/btcrpc/scenario_tests/shared"
)

// UTXOManagementScenarioTest 測試 UTXO 管理情境
// 這個測試模擬以下場景：
// 1. 創建錢包並獲取多個 UTXO
// 2. 分析 UTXO 分布和狀態
// 3. 執行小額和大額轉帳來測試 UTXO 選擇
// 4. 驗證 UTXO 的變化和交易費用計算
// 5. 測試地址重用和隱私保護
func UTXOManagementScenarioTest() {
	client := shared.NewTestClient()

	fmt.Println("=== UTXO 管理情境測試 ===")

	walletName := "utxo_test_wallet"
	recipientWallet := "utxo_recipient_wallet"

	// Step 1: 創建測試錢包
	fmt.Println("\n🔸 步驟 1: 創建測試錢包...")

	_, err := client.CreateWallet(walletName, false, false, "", false)
	if err != nil {
		log.Printf("Warning: 主錢包可能已存在: %v", err)
	}

	_, err = client.CreateWallet(recipientWallet, false, false, "", false)
	if err != nil {
		log.Printf("Warning: 接收錢包可能已存在: %v", err)
	}

	// Step 2: 生成多個地址並獲得不同金額的 UTXO
	fmt.Println("\n🔸 步驟 2: 創建多個地址並分別挖礦...")

	addresses := make([]string, 5)
	for i := 0; i < 5; i++ {
		addr, err := client.GetNewAddress(walletName, fmt.Sprintf("addr_%d", i), "bech32")
		if err != nil {
			log.Fatalf("無法生成地址 %d: %v", i, err)
		}
		addresses[i] = addr
		fmt.Printf("   地址 %d: %s\n", i+1, addr)
	}

	// 為不同地址挖礦不同數量的區塊來創建不同大小的 UTXO
	mineAmounts := []int{10, 20, 5, 15, 25} // 挖礦區塊數

	for i, addr := range addresses {
		fmt.Printf("   為地址 %d 挖掘 %d 個區塊...\n", i+1, mineAmounts[i])
		_, err := client.GenerateToAddress(mineAmounts[i], addr, nil)
		if err != nil {
			log.Fatalf("為地址 %d 挖礦失敗: %v", i+1, err)
		}
	}

	// 再挖掘一些區塊來確保所有獎勵都可以使用
	fmt.Println("   挖掘確認區塊...")
	_, err = client.GenerateToAddress(10, addresses[0], nil)
	if err != nil {
		log.Fatalf("挖掘確認區塊失敗: %v", err)
	}

	time.Sleep(3 * time.Second)

	// Step 3: 分析 UTXO 分布
	fmt.Println("\n🔸 步驟 3: 分析 UTXO 分布...")

	utxos, err := client.ListUnspent(walletName, 1, 9999999, nil, false, nil)
	if err != nil {
		log.Fatalf("無法列出 UTXO: %v", err)
	}

	fmt.Printf("✓ 找到 %d 個 UTXO:\n", len(utxos))

	// 按金額排序 UTXO
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Amount > utxos[j].Amount
	})

	totalAmount := 0.0
	for i, utxo := range utxos {
		totalAmount += utxo.Amount
		fmt.Printf("   UTXO %d: %.8f BTC (確認: %d, 地址: %s...)\n",
			i+1, utxo.Amount, utxo.Confirmations, utxo.Address[:10])
	}

	fmt.Printf("✓ 總 UTXO 金額: %.8f BTC\n", totalAmount)

	// Step 4: 測試小額轉帳 (應該選擇最小的合適 UTXO)
	fmt.Println("\n🔸 步驟 4: 測試小額轉帳...")

	recipientAddr, err := client.GetNewAddress(recipientWallet, "recipient_1", "bech32")
	if err != nil {
		log.Fatalf("無法生成接收地址: %v", err)
	}

	smallAmount := 1.0
	fmt.Printf("   轉帳金額: %.8f BTC 到 %s\n", smallAmount, recipientAddr[:20]+"...")

	// 記錄轉帳前的 UTXO
	utxosBefore := len(utxos)

	txid1, err := client.SendToAddressSimple(walletName, recipientAddr, smallAmount)
	if err != nil {
		log.Fatalf("小額轉帳失敗: %v", err)
	}

	fmt.Printf("✓ 小額轉帳交易 ID: %s\n", txid1[:16]+"...")

	// 確認交易
	_, err = client.GenerateToAddress(1, addresses[0], nil)
	if err != nil {
		log.Fatalf("確認小額轉帳失敗: %v", err)
	}

	time.Sleep(2 * time.Second)

	// 檢查轉帳後的 UTXO 變化
	utxosAfterSmall, err := client.ListUnspent(walletName, 1, 9999999, nil, false, nil)
	if err != nil {
		log.Fatalf("無法列出轉帳後 UTXO: %v", err)
	}

	fmt.Printf("✓ 轉帳前 UTXO 數量: %d, 轉帳後: %d\n", utxosBefore, len(utxosAfterSmall))

	// Step 5: 測試大額轉帳 (可能需要組合多個 UTXO)
	fmt.Println("\n🔸 步驟 5: 測試大額轉帳...")

	largeAmount := totalAmount * 0.8 // 使用 80% 的餘額
	recipientAddr2, err := client.GetNewAddress(recipientWallet, "recipient_2", "bech32")
	if err != nil {
		log.Fatalf("無法生成第二個接收地址: %v", err)
	}

	fmt.Printf("   轉帳金額: %.8f BTC 到 %s\n", largeAmount, recipientAddr2[:20]+"...")

	// 檢查當前餘額
	currentBalance, err := client.GetBalance(walletName, nil, nil)
	if err != nil {
		log.Fatalf("無法獲取當前餘額: %v", err)
	}

	fmt.Printf("   當前可用餘額: %.8f BTC\n", currentBalance)

	if largeAmount > currentBalance {
		largeAmount = currentBalance * 0.9 // 調整為 90% 以留出手續費空間
		fmt.Printf("   調整轉帳金額為: %.8f BTC\n", largeAmount)
	}

	txid2, err := client.SendToAddressSimple(walletName, recipientAddr2, largeAmount)
	if err != nil {
		log.Fatalf("大額轉帳失敗: %v", err)
	}

	fmt.Printf("✓ 大額轉帳交易 ID: %s\n", txid2[:16]+"...")

	// 確認交易
	_, err = client.GenerateToAddress(1, addresses[0], nil)
	if err != nil {
		log.Fatalf("確認大額轉帳失敗: %v", err)
	}

	time.Sleep(2 * time.Second)

	// Step 6: 分析轉帳後的 UTXO 狀態
	fmt.Println("\n🔸 步驟 6: 分析最終 UTXO 狀態...")

	finalUTXOs, err := client.ListUnspent(walletName, 1, 9999999, nil, false, nil)
	if err != nil {
		log.Fatalf("無法列出最終 UTXO: %v", err)
	}

	fmt.Printf("✓ 最終 UTXO 數量: %d\n", len(finalUTXOs))

	finalBalance := 0.0
	for i, utxo := range finalUTXOs {
		finalBalance += utxo.Amount
		fmt.Printf("   UTXO %d: %.8f BTC (確認: %d)\n",
			i+1, utxo.Amount, utxo.Confirmations)
	}

	fmt.Printf("✓ 最終餘額: %.8f BTC\n", finalBalance)

	// Step 7: 驗證接收方餘額
	fmt.Println("\n🔸 步驟 7: 驗證接收方餘額...")

	recipientBalance, err := client.GetBalance(recipientWallet, nil, nil)
	if err != nil {
		log.Fatalf("無法獲取接收方餘額: %v", err)
	}

	fmt.Printf("✓ 接收方總餘額: %.8f BTC\n", recipientBalance)

	expectedReceived := smallAmount + largeAmount
	if recipientBalance != expectedReceived {
		log.Fatalf("接收方餘額不正確，期望: %.8f，實際: %.8f", expectedReceived, recipientBalance)
	}

	// Step 8: 檢查交易詳情和費用
	fmt.Println("\n🔸 步驟 8: 分析交易詳情...")

	// 獲取第一筆交易詳情
	tx1Details, err := client.GetTransaction(walletName, txid1, false, true)
	if err != nil {
		log.Fatalf("無法獲取交易1詳情: %v", err)
	}

	// 獲取第二筆交易詳情
	tx2Details, err := client.GetTransaction(walletName, txid2, false, true)
	if err != nil {
		log.Fatalf("無法獲取交易2詳情: %v", err)
	}

	fmt.Printf("✓ 小額轉帳手續費: %.8f BTC\n", -tx1Details.Fee)
	fmt.Printf("✓ 大額轉帳手續費: %.8f BTC\n", -tx2Details.Fee)

	// Step 9: 驗證地址重用情況
	fmt.Println("\n🔸 步驟 9: 檢查地址使用情況...")

	addressGroupings, err := client.ListAddressGroupings(walletName)
	if err != nil {
		log.Fatalf("無法獲取地址分組: %v", err)
	}

	fmt.Printf("✓ 地址分組數量: %d\n", len(addressGroupings))
	for i, group := range addressGroupings {
		fmt.Printf("   分組 %d: %d 個地址\n", i+1, len(group))
		for j, addr := range group {
			fmt.Printf("     地址 %d: %s (%.8f BTC)\n",
				j+1, addr.Address[:20]+"...", addr.Amount)
		}
	}

	fmt.Println("\n✅ UTXO 管理情境測試完成！")
	fmt.Printf("📊 測試結果總結:\n")
	fmt.Printf("   - 初始 UTXO 數量: %d\n", utxosBefore)
	fmt.Printf("   - 最終 UTXO 數量: %d\n", len(finalUTXOs))
	fmt.Printf("   - 執行了 2 筆轉帳交易\n")
	fmt.Printf("   - 小額轉帳: %.8f BTC (手續費: %.8f BTC)\n", smallAmount, -tx1Details.Fee)
	fmt.Printf("   - 大額轉帳: %.8f BTC (手續費: %.8f BTC)\n", largeAmount, -tx2Details.Fee)
	fmt.Printf("   - 總手續費: %.8f BTC\n", (-tx1Details.Fee)+(-tx2Details.Fee))
	fmt.Printf("   - 餘額驗證通過 ✓\n")
}

func main() {
	UTXOManagementScenarioTest()
}
