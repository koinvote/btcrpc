# Bitcoin Core RPC 測試程式 Makefile

# 變數定義
GOCMD=go
GORUN=$(GOCMD) run
EXAMPLES_DIR=examples
SCENARIO_DIR=scenario_tests

.PHONY: help test-all basic wallet blockchain full scenario-all scenario-transfer scenario-utxo scenario-multisig scenario-fee scenario-monitor clean mod-check check-connection

# 預設目標
help:
	@echo "Bitcoin Core RPC 測試程式"
	@echo ""
	@echo "基本測試命令:"
	@echo "  make test-all      - 運行所有基本測試"
	@echo "  make basic         - 運行基本連線測試"
	@echo "  make wallet        - 運行錢包操作示例"
	@echo "  make blockchain    - 運行區塊鏈操作示例" 
	@echo "  make full          - 運行完整功能測試"
	@echo ""
	@echo "情境測試命令:"
	@echo "  make scenario-all      - 運行所有情境測試"
	@echo "  make scenario-transfer - 多錢包轉帳情境測試"
	@echo "  make scenario-utxo     - UTXO 管理情境測試"
	@echo "  make scenario-multisig - 多重簽名錢包情境測試"
	@echo "  make scenario-fee      - 手續費優化情境測試"
	@echo "  make scenario-monitor  - 區塊鏈監控情境測試"
	@echo ""
	@echo "工具命令:"
	@echo "  make clean         - 清理建置檔案"
	@echo "  make check-connection - 檢查 Bitcoin Core 連線"
	@echo ""
	@echo "確保 Bitcoin Core 在以下設定下運行:"
	@echo "  RPC URL: http://bitcoin-core:18443"
	@echo "  RPC User: bitcoinrpc"
	@echo "  RPC Password: test_password"

# 運行所有基本測試
test-all: basic wallet blockchain full

# 基本連線測試
basic:
	@echo "=== 運行基本連線測試 ==="
	$(GORUN) $(EXAMPLES_DIR)/basic_example.go

# 錢包操作示例
wallet:
	@echo "=== 運行錢包操作示例 ==="
	$(GORUN) $(EXAMPLES_DIR)/wallet_example.go

# 區塊鏈操作示例
blockchain:
	@echo "=== 運行區塊鏈操作示例 ==="
	$(GORUN) $(EXAMPLES_DIR)/blockchain_example.go

# 完整功能測試
full:
	@echo "=== 運行完整功能測試 ==="
	$(GORUN) $(EXAMPLES_DIR)/test_btcrpc.go

# 運行所有情境測試
scenario-all: scenario-transfer scenario-utxo scenario-multisig scenario-fee scenario-monitor
	@echo ""
	@echo "🎉 所有情境測試執行完成！"

# 多錢包轉帳情境測試
scenario-transfer:
	@echo ""
	@echo "🔸 執行多錢包轉帳情境測試..."
	@echo "=============================================="
	$(GORUN) $(SCENARIO_DIR)/transfer_test/multi_wallet_transfer_test.go

# UTXO 管理情境測試
scenario-utxo:
	@echo ""
	@echo "🔸 執行 UTXO 管理情境測試..."
	@echo "=============================================="
	$(GORUN) $(SCENARIO_DIR)/utxo_test/utxo_management_test.go

# 多重簽名錢包情境測試
scenario-multisig:
	@echo ""
	@echo "🔸 執行多重簽名錢包情境測試..."
	@echo "=============================================="
	$(GORUN) $(SCENARIO_DIR)/multisig_test/multisig_wallet_test.go

# 手續費優化情境測試
scenario-fee:
	@echo ""
	@echo "🔸 執行手續費優化情境測試..."
	@echo "=============================================="
	$(GORUN) $(SCENARIO_DIR)/fee_test/fee_optimization_test.go

# 區塊鏈監控情境測試
scenario-monitor:
	@echo ""
	@echo "🔸 執行區塊鏈監控情境測試..."
	@echo "=============================================="
	$(GORUN) $(SCENARIO_DIR)/monitor_test/blockchain_monitoring_test.go

# 清理
clean:
	$(GOCMD) clean
	rm -f *.exe *.out

# 檢查 Go 模組
mod-check:
	$(GOCMD) mod verify
	$(GOCMD) mod tidy

# 檢查連線
check-connection:
	@echo "檢查 Bitcoin Core 連線..."
	@curl -s --user bitcoinrpc:test_password --data-binary '{"jsonrpc":"1.0","id":"test","method":"getblockchaininfo","params":[]}' -H 'content-type: text/plain;' http://bitcoin-core:18443/ > /dev/null && echo "✓ Bitcoin Core 連線成功" || echo "❌ Bitcoin Core 連線失敗"