# Bitcoin Core RPC 工具函式進階情境測試

## 📋 測試概述

本文件說明針對 btcrpc 工具函式庫所設計的進階情境測試套件。這些測試專注於驗證 RPC 功能在真實業務場景中的表現，而非簡單的 API 調用測試。透過模擬複雜的多步驟操作，確保工具函式在生產環境中的穩定性和正確性。

## 🎯 測試設計理念

### 情境導向方法
這些測試採用情境導向的設計思維，不同於傳統的單元測試：

- **真實業務流程**：每個測試都模擬實際使用中的完整業務流程
- **多函式協作**：驗證多個 RPC 函式間的協同作業效果
- **狀態連續性**：測試跨多個操作步驟的狀態變化和數據一致性

### 測試深度與廣度
- **端到端驗證**：從錢包創建到交易確認的完整操作鏈
- **資金流追蹤**：詳細驗證每筆資金變動的正確性
- **邊界條件**：測試各種異常情況和邊界條件的處理
- **性能評估**：分析不同操作策略對效率和成本的影響

### 測試目標

這些測試旨在驗證 Bitcoin Core RPC API 在複雜業務場景下的穩定性和正確性，包括：
- 多錢包資金管理
- UTXO 選擇和優化
- 多重簽名安全機制
- 手續費估算和優化
- 區塊鏈監控和交易追蹤

## 🏗️ 測試結構

### 涵蓋的工具函式模組
這些情境測試驗證了以下 btcrpc 工具函式的功能：

**錢包管理函式** (`wallet.go`)
- `CreateWallet`, `LoadWallet`, `ListWallets`
- `GetNewAddress`, `GetBalance`
- `SendToAddress`, `SendToAddressSimple`
- `ListTransactions`, `GetTransaction`
- `ListUnspent`, `ValidateAddress`
- `GetWalletInfo`, `ListAddressGroupings`

**區塊鏈操作函式** (`blockchain.go`)
- `GetBlockchainInfo`, `GetNetworkInfo`
- `GenerateToAddress`, `GetBlock`, `GetBlockHash`
- `GetRawTransaction`, `GetMempoolInfo`

**進階功能函式** (`advanced.go`)
- `CreateRawTransaction`, `SignRawTransactionWithWallet`
- `SendRawTransaction`, `DumpPrivKey`, `ImportPrivKey`
- `CreateMultisig`, `AddMultisigAddress`
- `GetRawMempool`, `EstimateSmartFee`

### 測試檔案組織
```
scenario_tests/
├── README.md                       # 詳細測試說明
├── transfer_test/                  # 多錢包轉帳情境
├── utxo_test/                      # UTXO 管理情境
├── multisig_test/                  # 多重簽名情境
├── fee_test/                       # 手續費優化情境
└── monitor_test/                   # 區塊鏈監控情境
```

## 🧪 情境測試詳述

### 1. 多錢包轉帳情境測試 (`multi_wallet_transfer_test.go`)
**測試目標：** 驗證錢包管理和轉帳函式的協同運作

**涉及函式：**
- `CreateWallet()` - 創建多個獨立錢包
- `GetNewAddress()` - 生成接收地址
- `GenerateToAddress()` - 挖礦獲取測試資金
- `GetBalance()` - 查詢錢包餘額
- `SendToAddressSimple()` - 執行轉帳操作
- `ListTransactions()` - 驗證交易記錄

**情境流程：** Alice 和 Bob 錢包創建 → 資金初始化 → 雙向轉帳 → 餘額驗證 → 交易歷史確認

### 2. UTXO 管理情境測試 (`utxo_management_test.go`)
**測試目標：** 驗證 UTXO 相關函式在複雜場景下的表現

**涉及函式：**
- `ListUnspent()` - 查詢和分析 UTXO
- `SendToAddressSimple()` - 創建不同大小的交易
- `GetTransaction()` - 分析交易詳情和手續費
- `ListAddressGroupings()` - 檢查地址分組
- `GetBalance()` - 驗證餘額變化

**情境流程：** 多樣化 UTXO 創建 → UTXO 分布分析 → 不同規模轉帳測試 → 手續費效率評估

### 3. 多重簽名情境測試 (`multisig_wallet_test.go`)
**測試目標：** 驗證多重簽名相關函式的安全性和可靠性

**涉及函式：**
- `CreateMultisig()` - 創建多重簽名地址
- `AddMultisigAddress()` - 將多重簽名地址加入錢包
- `ValidateAddress()` - 驗證地址有效性
- `CreateRawTransaction()` - 創建原始交易
- `SignRawTransactionWithWallet()` - 多方簽名
- `SendRawTransaction()` - 廣播已簽名交易

**情境流程：** 3方錢包設置 → 2-of-3 多重簽名創建 → 資金充值 → 多方簽名轉帳 → 安全性驗證

### 4. 手續費優化情境測試 (`fee_optimization_test.go`)
**測試目標：** 驗證手續費相關函式的準確性和優化策略

**涉及函式：**
- `EstimateSmartFee()` - 智能手續費估算
- `SendToAddress()` - 設置不同手續費參數的轉帳
- `GetTransaction()` - 分析實際手續費
- `GetRawTransaction()` - 計算交易大小和費率
- `GetMempoolInfo()` - 檢查內存池狀態
- `ListUnspent()` - UTXO 合併建議分析

**情境流程：** 手續費估算測試 → 不同費率交易執行 → 實際費用分析 → RBF 功能驗證 → 優化建議

### 5. 區塊鏈監控情境測試 (`blockchain_monitoring_test.go`)
**測試目標：** 驗證區塊鏈監控和數據分析函式的完整性

**涉及函式：**
- `GetBlockchainInfo()` - 獲取區塊鏈狀態
- `GetNetworkInfo()` - 獲取網絡信息
- `GetBlock()` - 分析區塊內容
- `GetBlockHash()` - 查詢特定高度區塊
- `GetRawTransaction()` - 構建交易圖譜
- `GetMempoolInfo()` - 監控內存池狀態

**情境流程：** 初始狀態監控 → 交易創建和追蹤 → 區塊確認過程 → 交易圖譜分析 → 狀態變化驗證

## 🚀 執行測試

### 前置條件

1. **Bitcoin Core 設定：**
   ```bash
   # regtest 模式設定
   RPC URL: http://bitcoin-core:18443
   RPC User: bitcoinrpc
   RPC Password: test_password
   Network: regtest
   ```

2. **Go 環境：**
   ```bash
   Go 1.24.7 或更高版本
   ```

### 運行方式

#### 檢查連線
```bash
make check-connection
```

#### 運行單個情境測試
```bash
# 多錢包轉帳測試
make scenario-transfer

# UTXO 管理測試  
make scenario-utxo

# 多重簽名測試
make scenario-multisig

# 手續費優化測試
make scenario-fee

# 區塊鏈監控測試
make scenario-monitor
```

#### 運行所有情境測試
```bash
make scenario-all
```

#### 查看幫助
```bash
make help
```

### 測試輸出

每個測試都會提供詳細的執行日誌，包括：
- 📋 測試步驟說明
- ✅ 成功操作確認
- 📊 數據分析結果
- 🔍 驗證結果摘要

範例輸出：
```
=== 多錢包轉帳情境測試 ===

🔸 步驟 1: 創建錢包...
✓ 成功創建 alice_wallet
✓ 成功創建 bob_wallet

🔸 步驟 2: 生成錢包地址...
✓ Alice 地址: bcrt1q...
✓ Bob 地址: bcrt1q...

...

✅ 多錢包轉帳情境測試完成！
📊 測試結果總結:
   - Alice 從 5000.00000000 BTC → 4990.00000000 BTC
   - Bob 從 0 BTC → 5.00000000 BTC
   - 總計執行了 2 筆轉帳交易
   - 所有餘額驗證通過 ✓
```

## 🔧 故障排除

### 常見問題

1. **連線失敗**
   ```
   ❌ 無法連接到 Bitcoin Core
   ```
   **解決方案：** 確認 Bitcoin Core 正在運行且 RPC 設定正確

2. **錢包已存在**
   ```
   Warning: 錢包可能已存在
   ```
   **說明：** 這是正常警告，測試會繼續進行

3. **餘額不足**
   ```
   ❌ Alice 餘額為 0，無法進行轉帳測試
   ```
   **解決方案：** 確保 regtest 網絡正常，重新執行測試

### 除錯技巧

1. **檢查區塊鏈狀態：**
   ```bash
   curl -s --user bitcoinrpc:test_password \
        --data-binary '{"jsonrpc":"1.0","id":"test","method":"getblockchaininfo","params":[]}' \
        -H 'content-type: text/plain;' \
        http://bitcoin-core:18443/
   ```

2. **查看錢包列表：**
   ```bash
   curl -s --user bitcoinrpc:test_password \
        --data-binary '{"jsonrpc":"1.0","id":"test","method":"listwallets","params":[]}' \
        -H 'content-type: text/plain;' \
        http://bitcoin-core:18443/
   ```

## � 測試特色

### 業務邏輯驗證
- ✅ **資金守恆定律**：確保所有資金變動都有明確的來源和去向
- ✅ **狀態一致性**：驗證複雜操作後各個組件狀態的一致性  
- ✅ **異常恢復**：測試系統在異常情況下的恢復能力

### 性能和效率
- ⚡ **UTXO 優化**：分析和優化交易輸入選擇策略
- � **手續費效率**：測試不同手續費策略的成本效益
- 📈 **擴展性**：驗證系統在高負載下的表現

### 安全性保證
- 🔒 **多重簽名**：驗證企業級安全機制的可靠性
- 🔐 **私鑰管理**：測試密鑰管理和權限控制
- 🛡️ **交易完整性**：確保交易數據的不可篡改性

## 📈 測試價值與意義

### 對 btcrpc 工具函式庫的驗證
- **函式可靠性**：確保各個 RPC 函式在複雜場景下正確運行
- **函式協作性**：驗證多個函式組合使用時的穩定性
- **邊界條件處理**：測試函式在異常情況下的錯誤處理能力
- **性能表現**：評估函式在不同負載下的執行效率

### 對使用者的開發支持
- **最佳實踐範例**：提供真實場景的函式使用示範
- **整合測試模板**：為基於 btcrpc 的應用開發提供測試參考
- **問題預防**：提前發現和解決潛在的整合問題
- **文檔補充**：透過實際測試展示函式的進階用法

### 對專案品質的提升
- **回歸測試**：確保新版本不會破壞現有功能
- **功能完整性**：驗證所有主要 RPC 功能都被正確實現
- **錯誤處理**：測試各種異常情況的處理機制
- **性能基準**：建立函式效能的基準參考點

## 🔍 測試執行建議

### 測試環境要求
- **Bitcoin Core**: regtest 模式運行
- **RPC 設定**: 確保 RPC 服務正確配置
- **Go 環境**: 1.24.7 或更高版本
- **網絡隔離**: 獨立的測試環境，避免影響其他系統

### 執行策略
- **漸進式測試**: 建議先執行單個測試熟悉流程
- **完整驗證**: 定期執行所有測試確保整體功能正常
- **日誌分析**: 詳細檢查測試輸出以了解函式行為
- **性能監控**: 關注測試執行時間和資源使用情況

### 擴展建議
- **自定義情境**: 基於實際需求添加新的測試情境
- **參數調整**: 根據不同環境調整測試參數
- **錯誤注入**: 增加故意製造錯誤的測試案例
- **壓力測試**: 在高負載情況下驗證函式穩定性

## 🤝 貢獻指南

歡迎為這個測試套件貢獻新的情境測試！

### 新增測試的建議結構：

1. **明確的業務場景**
2. **多步驟驗證流程**
3. **詳細的執行日誌**
4. **完整的結果驗證**
5. **異常情況處理**

### 程式碼風格：

- 使用中文注釋說明業務邏輯
- 提供詳細的執行步驟輸出
- 包含完整的錯誤處理
- 驗證所有關鍵數據點

---

**注意：** 這些測試專為 regtest 環境設計，請勿在主網或測試網使用真實資金運行！