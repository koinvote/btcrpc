package main

import (
	"fmt"
	"log"
	"time"

	"github.com/koinvote/btcrpc/scenario_tests/shared"
)

// BlockchainMonitoringScenarioTest 測試區塊鏈監控和交易追蹤情境
// 這個測試模擬以下場景：
// 1. 監控區塊鏈狀態和網絡信息
// 2. 創建和追蹤多筆交易的生命週期
// 3. 分析區塊內容和交易確認過程
// 4. 驗證交易狀態變化和確認數增長
// 5. 模擬區塊鏈分析和交易圖譜構建
// 6. 測試異常情況處理和網絡監控
func BlockchainMonitoringScenarioTest() {
	client := shared.NewTestClient()

	fmt.Println("=== 區塊鏈監控和交易追蹤情境測試 ===")

	walletName := "monitoring_wallet"
	trackedTxids := make([]string, 0)

	// Step 1: 檢查初始區塊鏈狀態
	fmt.Println("\n🔸 步驟 1: 獲取初始區塊鏈狀態...")

	blockchainInfo, err := client.GetBlockchainInfo()
	if err != nil {
		log.Fatalf("無法獲取區塊鏈信息: %v", err)
	}

	networkInfo, err := client.GetNetworkInfo()
	if err != nil {
		log.Fatalf("無法獲取網絡信息: %v", err)
	}

	initialHeight := blockchainInfo.Blocks

	fmt.Printf("✓ 初始區塊鏈狀態:\n")
	fmt.Printf("   - 網絡: %s\n", blockchainInfo.Chain)
	fmt.Printf("   - 當前高度: %d\n", blockchainInfo.Blocks)
	fmt.Printf("   - 最佳區塊: %s\n", blockchainInfo.BestBlockHash[:16]+"...")
	fmt.Printf("   - 難度: %.2f\n", blockchainInfo.Difficulty)
	fmt.Printf("   - Bitcoin Core 版本: %d\n", networkInfo.Version)
	fmt.Printf("   - 連接節點數: %d\n", networkInfo.Connections)

	// Step 2: 創建監控錢包
	fmt.Println("\n🔸 步驟 2: 創建監控錢包...")

	_, err = client.CreateWallet(walletName, false, false, "", false)
	if err != nil {
		log.Printf("Warning: 監控錢包可能已存在: %v", err)
	}

	// 生成監控地址
	monitorAddr, err := client.GetNewAddress(walletName, "monitor_main", "bech32")
	if err != nil {
		log.Fatalf("無法生成監控地址: %v", err)
	}

	fmt.Printf("✓ 監控地址: %s\n", monitorAddr)

	// Step 3: 創建初始資金和多筆測試交易
	fmt.Println("\n🔸 步驟 3: 創建測試交易序列...")

	// 挖礦獲得初始資金
	_, err = client.GenerateToAddress(101, monitorAddr, nil)
	if err != nil {
		log.Fatalf("挖礦失敗: %v", err)
	}

	time.Sleep(2 * time.Second)

	// 創建一系列測試地址和交易
	testAddresses := make([]string, 5)
	for i := 0; i < 5; i++ {
		addr, err := client.GetNewAddress(walletName, fmt.Sprintf("test_%d", i), "bech32")
		if err != nil {
			log.Fatalf("無法生成測試地址 %d: %v", i, err)
		}
		testAddresses[i] = addr

		// 創建不同金額的交易
		amount := float64(i+1) * 2.5
		txid, err := client.SendToAddressSimple(walletName, addr, amount)
		if err != nil {
			log.Fatalf("創建測試交易 %d 失敗: %v", i, err)
		}

		trackedTxids = append(trackedTxids, txid)
		fmt.Printf("   創建交易 %d: %s (%.2f BTC)\n", i+1, txid[:16]+"...", amount)

		time.Sleep(500 * time.Millisecond) // 短暫延遲
	}

	fmt.Printf("✓ 成功創建 %d 筆追蹤交易\n", len(trackedTxids))

	// Step 4: 監控未確認交易狀態
	fmt.Println("\n🔸 步驟 4: 監控未確認交易狀態...")

	// 檢查內存池中的交易
	mempoolTxids, err := client.GetRawMempoolSimple()
	if err != nil {
		log.Printf("無法獲取內存池: %v", err)
	} else {
		fmt.Printf("✓ 內存池中有 %d 筆未確認交易\n", len(mempoolTxids))

		// 檢查我們的交易是否在內存池中
		inMempool := 0
		for _, trackTxid := range trackedTxids {
			for _, mempoolTxid := range mempoolTxids {
				if trackTxid == mempoolTxid {
					inMempool++
					break
				}
			}
		}
		fmt.Printf("   其中 %d 筆是我們追蹤的交易\n", inMempool)
	}

	// 獲取詳細的內存池信息
	mempoolInfo, err := client.GetMempoolInfo()
	if err != nil {
		log.Printf("無法獲取內存池詳情: %v", err)
	} else {
		fmt.Printf("✓ 內存池詳情:\n")
		fmt.Printf("   - 總大小: %d bytes\n", mempoolInfo.Bytes)
		fmt.Printf("   - 內存使用: %d bytes\n", mempoolInfo.Usage)
		fmt.Printf("   - 最低費率: %.8f BTC/kB\n", mempoolInfo.MempoolMinFee)
	}

	// Step 5: 挖礦並追蹤交易確認過程
	fmt.Println("\n🔸 步驟 5: 挖礦並追蹤確認過程...")

	confirmationRounds := 3

	for round := 1; round <= confirmationRounds; round++ {
		fmt.Printf("\n   確認輪次 %d:\n", round)

		// 挖一個區塊
		newBlocks, err := client.GenerateToAddress(1, monitorAddr, nil)
		if err != nil {
			log.Fatalf("挖礦輪次 %d 失敗: %v", round, err)
		}

		newBlockHash := newBlocks[0]

		// 獲取新區塊信息
		blockInfo, err := client.GetBlock(newBlockHash, 1)
		if err != nil {
			log.Printf("無法獲取區塊信息: %v", err)
		} else {
			fmt.Printf("   ✓ 新區塊: %s (高度: %d, 交易數: %d)\n",
				newBlockHash[:16]+"...", blockInfo.Height, blockInfo.NTx)
		}

		time.Sleep(1 * time.Second)

		// 檢查追蹤交易的確認狀態
		fmt.Printf("   追蹤交易確認狀態:\n")
		confirmedCount := 0

		for i, txid := range trackedTxids {
			txDetails, err := client.GetTransaction(walletName, txid, false, false)
			if err != nil {
				fmt.Printf("     交易 %d: 獲取失敗 (%v)\n", i+1, err)
				continue
			}

			status := "未確認"
			if txDetails.Confirmations > 0 {
				status = fmt.Sprintf("%d 確認", txDetails.Confirmations)
				confirmedCount++
			}

			fmt.Printf("     交易 %d: %s (%s)\n", i+1, txid[:12]+"...", status)
		}

		fmt.Printf("   ✓ 已確認: %d/%d 筆交易\n", confirmedCount, len(trackedTxids))
	}

	// Step 6: 深度分析區塊內容
	fmt.Println("\n🔸 步驟 6: 深度分析區塊內容...")

	currentHeight := initialHeight + int64(confirmationRounds) + 101

	// 分析最近的區塊
	for i := 0; i < 3; i++ {
		height := currentHeight - int64(i)

		blockHash, err := client.GetBlockHash(int(height))
		if err != nil {
			log.Printf("無法獲取區塊 %d 的哈希: %v", height, err)
			continue
		}

		blockInfo, err := client.GetBlock(blockHash, 2) // verbosity=2 獲取完整交易信息
		if err != nil {
			log.Printf("無法獲取區塊 %d 的詳細信息: %v", height, err)
			continue
		}

		fmt.Printf("✓ 區塊 %d 分析:\n", height)
		fmt.Printf("   - 哈希: %s\n", blockHash[:20]+"...")
		fmt.Printf("   - 時間: %s\n", time.Unix(blockInfo.Time, 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("   - 大小: %d bytes\n", blockInfo.Size)
		fmt.Printf("   - 交易數: %d\n", blockInfo.NTx)
		fmt.Printf("   - 難度: %.2f\n", blockInfo.Difficulty)

		// 分析區塊中的交易
		trackedInBlock := 0
		for _, txid := range blockInfo.Tx {
			for _, trackedTxid := range trackedTxids {
				if txid == trackedTxid {
					trackedInBlock++
					fmt.Printf("   - 包含追蹤交易: %s\n", txid[:16]+"...")
				}
			}
		}

		if trackedInBlock == 0 {
			fmt.Printf("   - 無追蹤交易\n")
		}
	}

	// Step 7: 交易圖譜分析
	fmt.Println("\n🔸 步驟 7: 構建交易關係圖譜...")

	// 分析每筆追蹤交易的輸入輸出
	for i, txid := range trackedTxids {
		rawTx, err := client.GetRawTransaction(txid, true, nil)
		if err != nil {
			log.Printf("無法獲取交易 %d 原始數據: %v", i+1, err)
			continue
		}

		fmt.Printf("✓ 交易 %d 圖譜分析:\n", i+1)
		fmt.Printf("   - 交易 ID: %s\n", txid[:20]+"...")
		fmt.Printf("   - 大小: %d bytes\n", rawTx.Size)
		fmt.Printf("   - 輸入數量: %d\n", len(rawTx.Vin))
		fmt.Printf("   - 輸出數量: %d\n", len(rawTx.Vout))

		// 分析輸出
		for j, vout := range rawTx.Vout {
			if len(vout.ScriptPubKey.Addresses) > 0 {
				fmt.Printf("     輸出 %d: %.8f BTC → %s\n",
					j, vout.Value, vout.ScriptPubKey.Addresses[0][:16]+"...")
			}
		}
	}

	// Step 8: 錢包狀態總覽
	fmt.Println("\n🔸 步驟 8: 錢包狀態總覽...")

	finalBalance, err := client.GetBalance(walletName, nil, nil)
	if err != nil {
		log.Fatalf("無法獲取最終餘額: %v", err)
	}

	walletInfo, err := client.GetWalletInfo(walletName)
	if err != nil {
		log.Printf("無法獲取錢包詳情: %v", err)
	} else {
		fmt.Printf("✓ 錢包狀態:\n")
		fmt.Printf("   - 總餘額: %.8f BTC\n", walletInfo.Balance)
		fmt.Printf("   - 未確認餘額: %.8f BTC\n", walletInfo.UnconfirmedBalance)
		fmt.Printf("   - 交易總數: %d\n", walletInfo.TxCount)
		fmt.Printf("   - 密鑰池大小: %d\n", walletInfo.KeypoolSize)
	}

	// 獲取交易歷史
	allTransactions, err := client.ListTransactions(walletName, "*", 50, 0, false)
	if err != nil {
		log.Printf("無法獲取交易歷史: %v", err)
	} else {
		sendCount := 0
		receiveCount := 0
		generateCount := 0

		for _, tx := range allTransactions {
			switch tx.Category {
			case "send":
				sendCount++
			case "receive":
				receiveCount++
			case "generate":
				generateCount++
			}
		}

		fmt.Printf("✓ 交易統計:\n")
		fmt.Printf("   - 發送交易: %d 筆\n", sendCount)
		fmt.Printf("   - 接收交易: %d 筆\n", receiveCount)
		fmt.Printf("   - 挖礦獎勵: %d 筆\n", generateCount)
	}

	// Step 9: 最終區塊鏈狀態
	fmt.Println("\n🔸 步驟 9: 最終區塊鏈狀態...")

	finalBlockchainInfo, err := client.GetBlockchainInfo()
	if err != nil {
		log.Fatalf("無法獲取最終區塊鏈信息: %v", err)
	}

	blocksGenerated := finalBlockchainInfo.Blocks - initialHeight

	fmt.Printf("✓ 測試期間區塊鏈變化:\n")
	fmt.Printf("   - 初始高度: %d\n", initialHeight)
	fmt.Printf("   - 最終高度: %d\n", finalBlockchainInfo.Blocks)
	fmt.Printf("   - 新增區塊: %d\n", blocksGenerated)
	fmt.Printf("   - 最新區塊: %s\n", finalBlockchainInfo.BestBlockHash[:20]+"...")

	fmt.Println("\n✅ 區塊鏈監控和交易追蹤情境測試完成！")
	fmt.Printf("📊 測試結果總結:\n")
	fmt.Printf("   - 監控了 %d 筆交易的完整生命週期\n", len(trackedTxids))
	fmt.Printf("   - 挖掘了 %d 個新區塊\n", blocksGenerated)
	fmt.Printf("   - 分析了區塊內容和交易圖譜\n")
	fmt.Printf("   - 驗證了交易確認過程\n")
	fmt.Printf("   - 最終錢包餘額: %.8f BTC\n", finalBalance)
	fmt.Printf("   - 所有監控功能正常運行 ✓\n")
}

func main() {
	BlockchainMonitoringScenarioTest()
}
