package btcrpc

import "encoding/json"

// RPCRequest 代表一個 JSON-RPC 請求 / represents a JSON-RPC request
type RPCRequest struct {
	Method  string        `json:"method"`  // 要調用的 RPC 方法名稱 / RPC method name to call
	Params  []interface{} `json:"params"`  // 方法參數陣列 / Array of method parameters
	ID      int           `json:"id"`      // 請求識別碼 / Request identifier
	JsonRPC string        `json:"jsonrpc"` // JSON-RPC 協議版本 / JSON-RPC protocol version
}

// RPCResponse 代表一個 JSON-RPC 回應 / represents a JSON-RPC response
type RPCResponse struct {
	Result json.RawMessage `json:"result"` // 方法執行結果的原始 JSON 數據 / Raw JSON data of method execution result
	Error  *RPCError       `json:"error"`  // 錯誤信息（如果有的話）/ Error information (if any)
	ID     int             `json:"id"`     // 對應請求的識別碼 / Identifier corresponding to the request
}

// RPCError 代表一個 JSON-RPC 錯誤 / represents a JSON-RPC error
type RPCError struct {
	Code    int    `json:"code"`    // 錯誤代碼 / Error code
	Message string `json:"message"` // 錯誤描述信息 / Error description message
}

// BlockchainInfo 代表 getblockchaininfo 的回應數據 / represents the response from getblockchaininfo
type BlockchainInfo struct {
	Chain                string  `json:"chain"`                // 區塊鏈網絡名稱（main、test、regtest）/ Blockchain network name (main, test, regtest)
	Blocks               int64   `json:"blocks"`               // 本地最佳區塊鏈的區塊數量 / Number of blocks in the local best block chain
	Headers              int64   `json:"headers"`              // 已驗證的區塊頭數量 / Number of validated headers
	BestBlockHash        string  `json:"bestblockhash"`        // 最佳區塊的哈希值 / Hash of the best block
	Difficulty           float64 `json:"difficulty"`           // 當前挖礦難度 / Current mining difficulty
	MedianTime           int64   `json:"mediantime"`           // 最佳區塊的中位時間 / Median time of the best block
	VerificationProgress float64 `json:"verificationprogress"` // 區塊鏈驗證進度（0-1）/ Blockchain verification progress (0-1)
	InitialBlockDownload bool    `json:"initialblockdownload"` // 是否正在進行初始區塊下載 / Whether initial block download is in progress
	ChainWork            string  `json:"chainwork"`            // 區塊鏈總工作量的十六進制表示 / Total chainwork in hex
	SizeOnDisk           int64   `json:"size_on_disk"`         // 區塊鏈在磁盤上的大小（字節）/ Size of blockchain on disk in bytes
	Pruned               bool    `json:"pruned"`               // 是否啟用區塊修剪 / Whether block pruning is enabled
}

// NetworkInfo 代表 getnetworkinfo 的回應數據 / represents the response from getnetworkinfo
type NetworkInfo struct {
	Version         int            `json:"version"`         // Bitcoin Core 版本號 / Bitcoin Core version number
	Subversion      string         `json:"subversion"`      // Bitcoin Core 子版本字符串 / Bitcoin Core subversion string
	ProtocolVersion int            `json:"protocolversion"` // P2P 協議版本 / P2P protocol version
	LocalServices   string         `json:"localservices"`   // 本地節點提供的服務（十六進制）/ Services provided by local node (hex)
	LocalRelay      bool           `json:"localrelay"`      // 是否中繼交易 / Whether transaction relay is enabled
	TimeOffset      int            `json:"timeoffset"`      // 時間偏移量（秒）/ Time offset in seconds
	Connections     int            `json:"connections"`     // 對等節點連接數 / Number of peer connections
	NetworkActive   bool           `json:"networkactive"`   // 網絡是否啟用 / Whether networking is enabled
	Networks        []Network      `json:"networks"`        // 支持的網絡列表 / List of supported networks
	RelayFee        float64        `json:"relayfee"`        // 最低中繼費用（BTC/kB）/ Minimum relay fee (BTC/kB)
	IncrementalFee  float64        `json:"incrementalfee"`  // 增量費用（BTC/kB）/ Incremental fee (BTC/kB)
	LocalAddresses  []LocalAddress `json:"localaddresses"`  // 本地地址列表 / List of local addresses
	Warnings        string         `json:"warnings"`        // 警告信息 / Warning messages
}

// Network 代表網絡信息 / represents network information
type Network struct {
	Name                      string `json:"name"`                        // 網絡名稱（ipv4、ipv6、onion、i2p）/ Network name (ipv4, ipv6, onion, i2p)
	Limited                   bool   `json:"limited"`                     // 網絡是否受限 / Whether the network is limited
	Reachable                 bool   `json:"reachable"`                   // 網絡是否可達 / Whether the network is reachable
	Proxy                     string `json:"proxy"`                       // 代理服務器地址 / Proxy server address
	ProxyRandomizeCredentials bool   `json:"proxy_randomize_credentials"` // 是否隨機化代理憑證 / Whether proxy credentials are randomized
}

// LocalAddress 代表本地地址信息 / represents local address information
type LocalAddress struct {
	Address string `json:"address"` // 本地 IP 地址 / Local IP address
	Port    int    `json:"port"`    // 端口號 / Port number
	Score   int    `json:"score"`   // 地址評分 / Address score
}

// CreateWalletResponse 代表 createwallet 的回應數據 / represents the response from createwallet
type CreateWalletResponse struct {
	Name    string `json:"name"`    // 創建的錢包名稱 / Name of the created wallet
	Warning string `json:"warning"` // 創建錢包時的警告信息 / Warning message during wallet creation
}

// GenerateToAddressResponse 代表 generatetoaddress 的回應，返回生成的區塊哈希陣列 / represents the response from generatetoaddress, returns an array of block hashes that were generated
type GenerateToAddressResponse []string

// BalanceResponse 代表 getbalance 的回應數據 / represents the response from getbalance
type BalanceResponse struct {
	Balance float64 `json:"balance"` // 錢包餘額（BTC）/ Wallet balance (BTC)
}

// SendToAddressResponse 代表 sendtoaddress 的回應，返回交易哈希（txid）/ represents the response from sendtoaddress, returns transaction hash (txid)
type SendToAddressResponse string

// Transaction 代表 listtransactions 返回的交易信息 / represents a transaction from listtransactions
type Transaction struct {
	Account           string   `json:"account"`                      // 帳戶名稱（已棄用）/ Account name (deprecated)
	Address           string   `json:"address"`                      // 交易相關的地址 / Address involved in the transaction
	Category          string   `json:"category"`                     // 交易類別：send（發送）、receive（接收）、generate（挖礦）、immature（未成熟）/ Transaction category: send, receive, generate, immature
	Amount            float64  `json:"amount"`                       // 交易金額（BTC）/ Transaction amount (BTC)
	Label             string   `json:"label"`                        // 地址標籤 / Address label
	Vout              int      `json:"vout"`                         // 輸出索引 / Output index
	Fee               float64  `json:"fee,omitempty"`                // 交易手續費（BTC）/ Transaction fee (BTC)
	Confirmations     int      `json:"confirmations"`                // 確認次數 / Number of confirmations
	BlockHash         string   `json:"blockhash,omitempty"`          // 包含此交易的區塊哈希 / Hash of block containing this transaction
	BlockIndex        int      `json:"blockindex,omitempty"`         // 交易在區塊中的索引 / Transaction index in the block
	BlockTime         int64    `json:"blocktime,omitempty"`          // 區塊時間戳 / Block timestamp
	TxID              string   `json:"txid"`                         // 交易 ID / Transaction ID
	WalletConflicts   []string `json:"walletconflicts"`              // 與此交易衝突的交易列表 / List of conflicting transactions
	Time              int64    `json:"time"`                         // 交易時間戳 / Transaction timestamp
	TimeReceived      int64    `json:"timereceived"`                 // 接收到交易的時間戳 / Time when transaction was received
	BIP125Replaceable string   `json:"bip125-replaceable,omitempty"` // 是否支持 BIP125 替換 / Whether BIP125 replacement is enabled
	Comment           string   `json:"comment,omitempty"`            // 交易註釋 / Transaction comment
	To                string   `json:"to,omitempty"`                 // 收款方註釋 / Recipient comment
}

// ListTransactionsResponse 代表 listtransactions 的回應，返回交易列表 / represents the response from listtransactions
type ListTransactionsResponse []Transaction

// ValidateAddressResponse 代表 validateaddress 的回應數據 / represents the response from validateaddress
type ValidateAddressResponse struct {
	IsValid        bool   `json:"isvalid"`                   // 地址是否有效 / Whether the address is valid
	Address        string `json:"address,omitempty"`         // 驗證的地址 / The validated address
	ScriptPubKey   string `json:"scriptPubKey,omitempty"`    // 地址對應的腳本公鑰（十六進制）/ Script public key for the address (hex)
	IsMine         bool   `json:"ismine,omitempty"`          // 地址是否屬於本錢包 / Whether the address belongs to this wallet
	IsScript       bool   `json:"isscript,omitempty"`        // 是否為腳本地址 / Whether it's a script address
	IsWitness      bool   `json:"iswitness,omitempty"`       // 是否為隔離見證地址 / Whether it's a witness address
	WitnessVersion int    `json:"witness_version,omitempty"` // 隔離見證版本 / Witness version
	WitnessProgram string `json:"witness_program,omitempty"` // 隔離見證程序（十六進制）/ Witness program (hex)
}

// UTXO 代表未花費交易輸出 / represents an unspent transaction output
type UTXO struct {
	TxID          string  `json:"txid"`                    // 交易 ID / Transaction ID
	Vout          int     `json:"vout"`                    // 輸出索引 / Output index
	Address       string  `json:"address"`                 // 輸出地址 / Output address
	Label         string  `json:"label"`                   // 地址標籤 / Address label
	ScriptPubKey  string  `json:"scriptPubKey"`            // 腳本公鑰（十六進制）/ Script public key (hex)
	Amount        float64 `json:"amount"`                  // 輸出金額（BTC）/ Output amount (BTC)
	Confirmations int     `json:"confirmations"`           // 確認次數 / Number of confirmations
	RedeemScript  string  `json:"redeemScript,omitempty"`  // 贖回腳本（十六進制）/ Redeem script (hex)
	WitnessScript string  `json:"witnessScript,omitempty"` // 隔離見證腳本（十六進制）/ Witness script (hex)
	Spendable     bool    `json:"spendable"`               // 是否可花費 / Whether it's spendable
	Solvable      bool    `json:"solvable"`                // 是否可解 / Whether it's solvable
	Safe          bool    `json:"safe"`                    // 是否安全（未受到雙花攻擊）/ Whether it's safe (not subject to double-spend)
}

// GetTransactionResponse 代表 gettransaction 的回應數據 / represents the response from gettransaction
type GetTransactionResponse struct {
	Amount            float64             `json:"amount"`                       // 交易淨金額（BTC）/ Net transaction amount (BTC)
	Fee               float64             `json:"fee,omitempty"`                // 交易手續費（BTC）/ Transaction fee (BTC)
	Confirmations     int                 `json:"confirmations"`                // 確認次數 / Number of confirmations
	BlockHash         string              `json:"blockhash,omitempty"`          // 包含此交易的區塊哈希 / Hash of block containing this transaction
	BlockIndex        int                 `json:"blockindex,omitempty"`         // 交易在區塊中的索引 / Transaction index in the block
	BlockTime         int64               `json:"blocktime,omitempty"`          // 區塊時間戳 / Block timestamp
	TxID              string              `json:"txid"`                         // 交易 ID / Transaction ID
	WalletConflicts   []string            `json:"walletconflicts"`              // 與此交易衝突的交易列表 / List of conflicting transactions
	Time              int64               `json:"time"`                         // 交易時間戳 / Transaction timestamp
	TimeReceived      int64               `json:"timereceived"`                 // 接收到交易的時間戳 / Time when transaction was received
	BIP125Replaceable string              `json:"bip125-replaceable,omitempty"` // 是否支持 BIP125 替換 / Whether BIP125 replacement is enabled
	Details           []TransactionDetail `json:"details"`                      // 交易詳細信息列表 / List of transaction details
	Hex               string              `json:"hex"`                          // 交易的十六進制表示 / Transaction in hex format
}

// TransactionDetail 代表 gettransaction 中的交易詳細信息 / represents transaction details within gettransaction
type TransactionDetail struct {
	Account   string  `json:"account,omitempty"`   // 帳戶名稱（已棄用）/ Account name (deprecated)
	Address   string  `json:"address,omitempty"`   // 交易相關的地址 / Address involved in the transaction
	Category  string  `json:"category"`            // 交易類別：send（發送）、receive（接收）/ Transaction category: send, receive
	Amount    float64 `json:"amount"`              // 交易金額（BTC）/ Transaction amount (BTC)
	Label     string  `json:"label,omitempty"`     // 地址標籤 / Address label
	Vout      int     `json:"vout"`                // 輸出索引 / Output index
	Fee       float64 `json:"fee,omitempty"`       // 交易手續費（BTC）/ Transaction fee (BTC)
	Abandoned bool    `json:"abandoned,omitempty"` // 交易是否被放棄 / Whether the transaction is abandoned
}

// EstimateSmartFeeResponse 代表 estimatesmartfee 的回應數據 / represents the response from estimatesmartfee
type EstimateSmartFeeResponse struct {
	FeeRate float64  `json:"feerate,omitempty"` // 估算的手續費率（BTC/kB）/ Estimated fee rate (BTC/kB)
	Errors  []string `json:"errors,omitempty"`  // 估算過程中的錯誤信息 / Errors during estimation
	Blocks  int      `json:"blocks"`            // 估算基於的區塊數 / Number of blocks used for estimation
}

// GetBlockResponse 代表 getblock 的回應數據 / represents the response from getblock
type GetBlockResponse struct {
	Hash              string   `json:"hash"`                        // 區塊哈希 / Block hash
	Confirmations     int      `json:"confirmations"`               // 確認次數 / Number of confirmations
	Height            int      `json:"height"`                      // 區塊高度 / Block height
	Version           int      `json:"version"`                     // 區塊版本 / Block version
	VersionHex        string   `json:"versionHex"`                  // 區塊版本（十六進制）/ Block version (hex)
	MerkleRoot        string   `json:"merkleroot"`                  // 默克爾樹根 / Merkle tree root
	Time              int64    `json:"time"`                        // 區塊時間戳 / Block timestamp
	MedianTime        int64    `json:"mediantime"`                  // 中位時間 / Median time
	Nonce             uint32   `json:"nonce"`                       // 隨機數 / Nonce
	Bits              string   `json:"bits"`                        // 難度目標（十六進制）/ Difficulty target (hex)
	Difficulty        float64  `json:"difficulty"`                  // 挖礦難度 / Mining difficulty
	ChainWork         string   `json:"chainwork"`                   // 區塊鏈總工作量（十六進制）/ Total chainwork (hex)
	NTx               int      `json:"nTx"`                         // 區塊中的交易數量 / Number of transactions in block
	PreviousBlockHash string   `json:"previousblockhash,omitempty"` // 前一個區塊的哈希 / Hash of previous block
	NextBlockHash     string   `json:"nextblockhash,omitempty"`     // 下一個區塊的哈希 / Hash of next block
	StrippedSize      int      `json:"strippedsize"`                // 去除見證數據的區塊大小 / Block size without witness data
	Size              int      `json:"size"`                        // 區塊大小（字節）/ Block size (bytes)
	Weight            int      `json:"weight"`                      // 區塊權重 / Block weight
	Tx                []string `json:"tx"`                          // 區塊中的交易 ID 列表 / List of transaction IDs in block
}

// GetWalletInfoResponse 代表 getwalletinfo 的回應數據 / represents the response from getwalletinfo
type GetWalletInfoResponse struct {
	WalletName            string      `json:"walletname"`                        // 錢包名稱 / Wallet name
	WalletVersion         int         `json:"walletversion"`                     // 錢包版本 / Wallet version
	Format                string      `json:"format"`                            // 錢包格式 / Wallet format
	Balance               float64     `json:"balance"`                           // 錢包餘額（BTC）/ Wallet balance (BTC)
	UnconfirmedBalance    float64     `json:"unconfirmed_balance"`               // 未確認餘額（BTC）/ Unconfirmed balance (BTC)
	ImmatureBalance       float64     `json:"immature_balance"`                  // 未成熟餘額（BTC）/ Immature balance (BTC)
	TxCount               int         `json:"txcount"`                           // 交易總數 / Total number of transactions
	KeypoolOldest         int64       `json:"keypoololdest"`                     // 密鑰池中最舊密鑰的時間戳 / Timestamp of oldest key in keypool
	KeypoolSize           int         `json:"keypoolsize"`                       // 密鑰池大小 / Size of keypool
	KeypoolSizeHdInternal int         `json:"keypoolsize_hd_internal,omitempty"` // HD 內部密鑰池大小 / Size of HD internal keypool
	UnlockedUntil         int64       `json:"unlocked_until,omitempty"`          // 錢包解鎖截止時間 / Time until wallet is unlocked
	PayTxFee              float64     `json:"paytxfee"`                          // 支付交易手續費（BTC/kB）/ Pay transaction fee (BTC/kB)
	HdSeedId              string      `json:"hdseedid,omitempty"`                // HD 種子 ID / HD seed ID
	PrivateKeysEnabled    bool        `json:"private_keys_enabled"`              // 是否啟用私鑰 / Whether private keys are enabled
	AvoidReuse            bool        `json:"avoid_reuse"`                       // 是否避免地址重用 / Whether address reuse is avoided
	Scanning              interface{} `json:"scanning"`                          // 掃描狀態 / Scanning status
	Descriptors           bool        `json:"descriptors"`                       // 是否使用描述符 / Whether descriptors are used
}

// AddressGrouping 代表一組地址及其餘額 / represents a group of addresses with their balances
type AddressGrouping []AddressInfo

// AddressInfo 代表地址分組中的地址信息 / represents address information in groupings
type AddressInfo struct {
	Address string  `json:"address"`         // 比特幣地址 / Bitcoin address
	Amount  float64 `json:"amount"`          // 地址餘額（BTC）/ Address balance (BTC)
	Label   string  `json:"label,omitempty"` // 地址標籤 / Address label
}

// ListAddressGroupingsResponse 代表 listaddressgroupings 的回應，返回地址分組列表 / represents the response from listaddressgroupings
type ListAddressGroupingsResponse []AddressGrouping

// GetBlockHashResponse 代表 getblockhash 的回應，返回區塊哈希 / represents the response from getblockhash
type GetBlockHashResponse string

// GetRawTransactionResponse 代表 getrawtransaction 的回應數據 / represents the response from getrawtransaction
type GetRawTransactionResponse struct {
	Hex           string               `json:"hex"`                     // 交易的十六進制表示 / Transaction in hex format
	TxID          string               `json:"txid,omitempty"`          // 交易 ID / Transaction ID
	Hash          string               `json:"hash,omitempty"`          // 交易哈希 / Transaction hash
	Size          int                  `json:"size,omitempty"`          // 交易大小（字節）/ Transaction size (bytes)
	VSize         int                  `json:"vsize,omitempty"`         // 虛擬交易大小 / Virtual transaction size
	Weight        int                  `json:"weight,omitempty"`        // 交易權重 / Transaction weight
	Version       int                  `json:"version,omitempty"`       // 交易版本 / Transaction version
	LockTime      int64                `json:"locktime,omitempty"`      // 鎖定時間 / Lock time
	Vin           []RawTransactionVin  `json:"vin,omitempty"`           // 交易輸入列表 / List of transaction inputs
	Vout          []RawTransactionVout `json:"vout,omitempty"`          // 交易輸出列表 / List of transaction outputs
	BlockHash     string               `json:"blockhash,omitempty"`     // 包含此交易的區塊哈希 / Hash of block containing this transaction
	Confirmations int                  `json:"confirmations,omitempty"` // 確認次數 / Number of confirmations
	BlockTime     int64                `json:"blocktime,omitempty"`     // 區塊時間戳 / Block timestamp
	Time          int64                `json:"time,omitempty"`          // 交易時間戳 / Transaction timestamp
}

// RawTransactionVin 代表原始交易中的交易輸入 / represents a transaction input in raw transaction
type RawTransactionVin struct {
	TxID        string                  `json:"txid,omitempty"`        // 輸入來源交易的 ID / ID of source transaction
	Vout        int                     `json:"vout,omitempty"`        // 輸入來源交易的輸出索引 / Output index of source transaction
	ScriptSig   RawTransactionScriptSig `json:"scriptSig,omitempty"`   // 腳本簽名 / Script signature
	TxInWitness []string                `json:"txinwitness,omitempty"` // 隔離見證數據 / Witness data
	Sequence    int64                   `json:"sequence,omitempty"`    // 序列號 / Sequence number
	Coinbase    string                  `json:"coinbase,omitempty"`    // 基礎幣交易數據（僅對挖礦交易）/ Coinbase data (mining transactions only)
}

// RawTransactionScriptSig 代表腳本簽名 / represents script signature
type RawTransactionScriptSig struct {
	Asm string `json:"asm"` // 腳本的彙編表示 / Assembly representation of script
	Hex string `json:"hex"` // 腳本的十六進制表示 / Hex representation of script
}

// RawTransactionVout 代表原始交易中的交易輸出 / represents a transaction output in raw transaction
type RawTransactionVout struct {
	Value        float64                    `json:"value"`        // 輸出金額（BTC）/ Output amount (BTC)
	N            int                        `json:"n"`            // 輸出索引 / Output index
	ScriptPubKey RawTransactionScriptPubKey `json:"scriptPubKey"` // 腳本公鑰 / Script public key
}

// RawTransactionScriptPubKey 代表腳本公鑰 / represents script public key
type RawTransactionScriptPubKey struct {
	Asm       string   `json:"asm"`                 // 腳本的彙編表示 / Assembly representation of script
	Hex       string   `json:"hex"`                 // 腳本的十六進制表示 / Hex representation of script
	ReqSigs   int      `json:"reqSigs,omitempty"`   // 所需簽名數量 / Required number of signatures
	Type      string   `json:"type"`                // 腳本類型（pubkey、pubkeyhash、scripthash 等）/ Script type (pubkey, pubkeyhash, scripthash, etc.)
	Addresses []string `json:"addresses,omitempty"` // 相關地址列表 / List of associated addresses
	Address   string   `json:"address,omitempty"`   // 相關地址（單個）/ Associated address (single)
}

// GetMempoolInfoResponse 代表 getmempoolinfo 的回應數據 / represents the response from getmempoolinfo
type GetMempoolInfoResponse struct {
	Loaded           bool    `json:"loaded"`           // 內存池是否已載入 / Whether mempool is loaded
	Size             int     `json:"size"`             // 內存池中的交易數量 / Number of transactions in mempool
	Bytes            int64   `json:"bytes"`            // 內存池交易的總大小（字節）/ Total size of mempool transactions (bytes)
	Usage            int64   `json:"usage"`            // 內存池的實際內存使用量（字節）/ Actual memory usage of mempool (bytes)
	MaxMempool       int64   `json:"maxmempool"`       // 內存池的最大大小（字節）/ Maximum size of mempool (bytes)
	MempoolMinFee    float64 `json:"mempoolminfee"`    // 內存池最低手續費率（BTC/kB）/ Minimum fee rate for mempool (BTC/kB)
	MinRelayTxFee    float64 `json:"minrelaytxfee"`    // 最低中繼手續費率（BTC/kB）/ Minimum relay fee rate (BTC/kB)
	UnbroadcastCount int     `json:"unbroadcastcount"` // 未廣播交易數量 / Number of unbroadcast transactions
}

// === Phase 5: Expert Level Functions ===

// CreateRawTransactionInput 代表創建原始交易的輸入 / represents an input for creating raw transaction
type CreateRawTransactionInput struct {
	TxID     string `json:"txid"`               // 輸入來源交易的 ID / ID of source transaction for input
	Vout     int    `json:"vout"`               // 輸入來源交易的輸出索引 / Output index of source transaction
	Sequence int64  `json:"sequence,omitempty"` // 序列號（可選）/ Sequence number (optional)
}

// CreateRawTransactionOutput 代表創建原始交易的輸出 / represents an output for creating raw transaction
type CreateRawTransactionOutput struct {
	Address string  `json:"address,omitempty"` // 目標地址 / Target address
	Amount  float64 `json:"amount,omitempty"`  // 輸出金額（BTC）/ Output amount (BTC)
	Data    string  `json:"data,omitempty"`    // 任意數據（十六進制）/ Arbitrary data (hex)
}

// SignRawTransactionResponse 代表 signrawtransactionwithwallet 的回應數據 / represents the response from signrawtransactionwithwallet
type SignRawTransactionResponse struct {
	Hex      string                    `json:"hex"`              // 簽名後的交易十六進制 / Signed transaction hex
	Complete bool                      `json:"complete"`         // 交易是否完全簽名 / Whether transaction is completely signed
	Errors   []SignRawTransactionError `json:"errors,omitempty"` // 簽名過程中的錯誤列表 / List of errors during signing
}

// SignRawTransactionError 代表簽名原始交易時的錯誤 / represents an error in signing raw transaction
type SignRawTransactionError struct {
	TxID      string `json:"txid"`      // 發生錯誤的交易 ID / Transaction ID where error occurred
	Vout      int    `json:"vout"`      // 發生錯誤的輸出索引 / Output index where error occurred
	ScriptSig string `json:"scriptSig"` // 腳本簽名 / Script signature
	Sequence  int64  `json:"sequence"`  // 序列號 / Sequence number
	Error     string `json:"error"`     // 錯誤描述 / Error description
}

// DumpPrivKeyResponse 代表 dumpprivkey 的回應數據 / represents the response from dumpprivkey
type DumpPrivKeyResponse struct {
	PrivateKey string `json:"privatekey"` // 導出的私鑰（WIF 格式）/ Exported private key (WIF format)
}

// ImportPrivKeyResponse 代表 importprivkey 的回應數據 / represents the response from importprivkey
type ImportPrivKeyResponse struct {
	Success bool `json:"success"` // 導入是否成功 / Whether import was successful
}

// CreateMultisigResponse 代表 createmultisig 的回應數據 / represents the response from createmultisig
type CreateMultisigResponse struct {
	Address      string `json:"address"`      // 多重簽名地址 / Multisig address
	RedeemScript string `json:"redeemScript"` // 贖回腳本（十六進制）/ Redeem script (hex)
}

// AddMultisigAddressResponse 代表 addmultisigaddress 的回應數據 / represents the response from addmultisigaddress
type AddMultisigAddressResponse struct {
	Address      string `json:"address"`      // 添加到錢包的多重簽名地址 / Multisig address added to wallet
	RedeemScript string `json:"redeemScript"` // 贖回腳本（十六進制）/ Redeem script (hex)
}

// GetRawMempoolResponse 代表 getrawmempool 的回應，返回內存池中的交易 ID 列表 / represents the response from getrawmempool
type GetRawMempoolResponse []string

// GetRawMempoolVerboseResponse 代表 getrawmempool 的詳細回應數據 / represents the verbose response from getrawmempool
type GetRawMempoolVerboseResponse map[string]GetRawMempoolEntry

// GetRawMempoolEntry 代表詳細內存池回應中的單個條目 / represents a single entry in verbose mempool response
type GetRawMempoolEntry struct {
	Vsize             int      `json:"vsize"`              // 虛擬交易大小 / Virtual transaction size
	Weight            int      `json:"weight"`             // 交易權重 / Transaction weight
	Fee               float64  `json:"fee"`                // 交易手續費（BTC）/ Transaction fee (BTC)
	ModifiedFee       float64  `json:"modifiedfee"`        // 修改後的手續費（BTC）/ Modified fee (BTC)
	Time              int64    `json:"time"`               // 交易進入內存池的時間 / Time when transaction entered mempool
	Height            int      `json:"height"`             // 交易進入內存池時的區塊高度 / Block height when transaction entered mempool
	DescendantCount   int      `json:"descendantcount"`    // 後代交易數量 / Number of descendant transactions
	DescendantSize    int      `json:"descendantsize"`     // 後代交易總大小 / Total size of descendant transactions
	DescendantFees    float64  `json:"descendantfees"`     // 後代交易總手續費（BTC）/ Total fees of descendant transactions (BTC)
	AncestorCount     int      `json:"ancestorcount"`      // 祖先交易數量 / Number of ancestor transactions
	AncestorSize      int      `json:"ancestorsize"`       // 祖先交易總大小 / Total size of ancestor transactions
	AncestorFees      float64  `json:"ancestorfees"`       // 祖先交易總手續費（BTC）/ Total fees of ancestor transactions (BTC)
	WTxID             string   `json:"wtxid"`              // 見證交易 ID / Witness transaction ID
	FeeRate           float64  `json:"feerate"`            // 手續費率（BTC/kB）/ Fee rate (BTC/kB)
	Depends           []string `json:"depends"`            // 依賴的交易 ID 列表 / List of dependent transaction IDs
	SpentBy           []string `json:"spentby"`            // 花費此交易輸出的交易 ID 列表 / List of transaction IDs spending this transaction's outputs
	BIP125Replaceable bool     `json:"bip125-replaceable"` // 是否支持 BIP125 替換 / Whether BIP125 replacement is enabled
	Unbroadcast       bool     `json:"unbroadcast"`        // 是否為未廣播交易 / Whether transaction is unbroadcast
}
