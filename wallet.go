package btcrpc

import (
	"encoding/json"
	"fmt"
)

// CreateWallet calls the createwallet RPC method
// walletName: name of the wallet to create
// disablePrivateKeys: disable the possibility of private keys (only watchonly/solvable addresses)
// blank: create a blank wallet (no seed, no keys)
// passphrase: encrypt the wallet with this passphrase
// avoidReuse: keep track of coin reuse for better privacy
func (c *Client) CreateWallet(walletName string, disablePrivateKeys, blank bool, passphrase string, avoidReuse bool) (*CreateWalletResponse, error) {
	// Prepare parameters
	params := []interface{}{walletName}

	// Add optional parameters if they differ from defaults
	if disablePrivateKeys || blank || passphrase != "" || avoidReuse {
		params = append(params, disablePrivateKeys)
	}
	if blank || passphrase != "" || avoidReuse {
		params = append(params, blank)
	}
	if passphrase != "" || avoidReuse {
		params = append(params, passphrase)
	}
	if avoidReuse {
		params = append(params, avoidReuse)
	}

	// Call the RPC method
	resp, err := c.call("createwallet", params)
	if err != nil {
		return nil, fmt.Errorf("failed to call createwallet: %v", err)
	}

	// Parse the result
	var result CreateWalletResponse
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create wallet response: %v", err)
	}

	return &result, nil
}

// LoadWallet calls the loadwallet RPC method
func (c *Client) LoadWallet(walletName string) error {
	// Call the RPC method
	_, err := c.call("loadwallet", []interface{}{walletName})
	if err != nil {
		return fmt.Errorf("failed to call loadwallet: %v", err)
	}
	return nil
}

// ListWallets calls the listwallets RPC method to get loaded wallets
func (c *Client) ListWallets() ([]string, error) {
	// Call the RPC method
	resp, err := c.call("listwallets", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to call listwallets: %v", err)
	}

	// Parse the result
	var wallets []string
	if err := json.Unmarshal(resp.Result, &wallets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallets list: %v", err)
	}

	return wallets, nil
}

// GetNewAddress calls the getnewaddress RPC method with specific wallet
// Note: This method requires a wallet to be loaded
// walletName: name of the wallet to use
// label: label for the address (optional)
// addressType: type of address to generate ("legacy", "p2sh-segwit", "bech32", "bech32m")
func (c *Client) GetNewAddress(walletName, label, addressType string) (string, error) {
	// Prepare parameters
	var params []interface{}

	// Add label if provided
	if label != "" {
		params = append(params, label)
	}

	// Add address type if provided
	if addressType != "" {
		// If we have addressType but no label, we need to add empty label first
		if label == "" {
			params = []interface{}{"", addressType}
		} else {
			params = append(params, addressType)
		}
	}

	// Call the RPC method with wallet endpoint
	resp, err := c.callWithWallet("getnewaddress", params, walletName)
	if err != nil {
		return "", fmt.Errorf("failed to call getnewaddress: %v", err)
	}

	// Parse the result (getnewaddress returns a string directly)
	var address string
	if err := json.Unmarshal(resp.Result, &address); err != nil {
		return "", fmt.Errorf("failed to unmarshal new address: %v", err)
	}

	return address, nil
}

// GetBalance calls the getbalance RPC method
// This method returns the wallet's available balance
// walletName: name of the wallet to check balance for
// minconf: minimum number of confirmations (optional, default 0)
// includeWatchonly: include watch-only addresses (optional, default false)
func (c *Client) GetBalance(walletName string, minconf *int, includeWatchonly *bool) (float64, error) {
	// Prepare parameters
	var params []interface{}

	// Add minconf if provided
	if minconf != nil {
		params = append(params, "*", *minconf) // "*" for all accounts (legacy parameter)
	}

	// Add include_watchonly if provided
	if includeWatchonly != nil {
		// If we have includeWatchonly but no minconf, use default minconf of 0
		if minconf == nil {
			params = []interface{}{"*", 0, *includeWatchonly}
		} else {
			params = append(params, *includeWatchonly)
		}
	}

	// Call the RPC method with wallet endpoint
	resp, err := c.callWithWallet("getbalance", params, walletName)
	if err != nil {
		return 0, fmt.Errorf("failed to call getbalance: %v", err)
	}

	// Parse the result (getbalance returns a number directly)
	var balance float64
	if err := json.Unmarshal(resp.Result, &balance); err != nil {
		return 0, fmt.Errorf("failed to unmarshal balance: %v", err)
	}

	return balance, nil
}

// SendToAddress calls the sendtoaddress RPC method
// walletName: name of the wallet to send from
// address: destination bitcoin address
// amount: amount to send (in BTC)
// comment: optional comment for the transaction
// commentTo: optional comment for the recipient
// subtractFeeFromAmount: if true, fee will be deducted from the amount being sent
// replaceable: allow transaction to be replaced by fee (BIP 125)
// confTarget: confirmation target in blocks for fee estimation
// estimateMode: fee estimation mode ("UNSET", "ECONOMICAL", "CONSERVATIVE")
func (c *Client) SendToAddress(walletName, address string, amount float64, comment, commentTo string, subtractFeeFromAmount, replaceable bool, confTarget int, estimateMode string) (string, error) {
	// Prepare parameters - address and amount are required
	params := []interface{}{address, amount}

	// Add optional parameters if they are provided
	if comment != "" || commentTo != "" || subtractFeeFromAmount || replaceable || confTarget > 0 || estimateMode != "" {
		params = append(params, comment)
	}
	if commentTo != "" || subtractFeeFromAmount || replaceable || confTarget > 0 || estimateMode != "" {
		params = append(params, commentTo)
	}
	if subtractFeeFromAmount || replaceable || confTarget > 0 || estimateMode != "" {
		params = append(params, subtractFeeFromAmount)
	}
	if replaceable || confTarget > 0 || estimateMode != "" {
		params = append(params, replaceable)
	}
	if confTarget > 0 || estimateMode != "" {
		params = append(params, confTarget)
	}
	if estimateMode != "" {
		params = append(params, estimateMode)
	}

	// Call the RPC method with wallet endpoint
	resp, err := c.callWithWallet("sendtoaddress", params, walletName)
	if err != nil {
		return "", fmt.Errorf("failed to call sendtoaddress: %v", err)
	}

	// Parse the result (sendtoaddress returns a transaction ID string)
	var txid string
	if err := json.Unmarshal(resp.Result, &txid); err != nil {
		return "", fmt.Errorf("failed to unmarshal transaction ID: %v", err)
	}

	return txid, nil
}

// ListTransactions calls the listtransactions RPC method
// walletName: name of the wallet to list transactions from
// label: optional label to filter transactions (use "*" for all)
// count: maximum number of transactions to return (default 10)
// skip: number of transactions to skip (for pagination)
// includeWatchonly: include watch-only transactions
func (c *Client) ListTransactions(walletName, label string, count, skip int, includeWatchonly bool) ([]Transaction, error) {
	// Prepare parameters
	var params []interface{}

	// Add label (use "*" for all transactions if not specified)
	if label == "" {
		label = "*"
	}
	params = append(params, label)

	// Add count if specified or use default
	if count == 0 {
		count = 10 // Bitcoin Core default
	}
	params = append(params, count)

	// Add skip if specified
	if skip > 0 {
		params = append(params, skip)
	}

	// Add includeWatchonly if needed
	if includeWatchonly {
		if skip == 0 {
			params = append(params, 0) // Add default skip
		}
		params = append(params, includeWatchonly)
	}

	// Call the RPC method with wallet endpoint
	resp, err := c.callWithWallet("listtransactions", params, walletName)
	if err != nil {
		return nil, fmt.Errorf("failed to call listtransactions: %v", err)
	}

	// Parse the result
	var transactions []Transaction
	if err := json.Unmarshal(resp.Result, &transactions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transactions list: %v", err)
	}

	return transactions, nil
}

// ValidateAddress calls the validateaddress RPC method
// address: bitcoin address to validate
func (c *Client) ValidateAddress(address string) (*ValidateAddressResponse, error) {
	// Prepare parameters
	params := []interface{}{address}

	// Call the RPC method (no wallet endpoint needed)
	resp, err := c.call("validateaddress", params)
	if err != nil {
		return nil, fmt.Errorf("failed to call validateaddress: %v", err)
	}

	// Parse the result
	var result ValidateAddressResponse
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal validate address response: %v", err)
	}

	return &result, nil
}

// SendToAddressSimple calls the sendtoaddress RPC method with minimal parameters for regtest
// This version is optimized for regtest environments where fee estimation might not work
func (c *Client) SendToAddressSimple(walletName, address string, amount float64) (string, error) {
	// Use only required parameters for regtest
	params := []interface{}{address, amount}

	// Call the RPC method with wallet endpoint
	resp, err := c.callWithWallet("sendtoaddress", params, walletName)
	if err != nil {
		return "", fmt.Errorf("failed to call sendtoaddress: %v", err)
	}

	// Parse the result (sendtoaddress returns a transaction ID string)
	var txid string
	if err := json.Unmarshal(resp.Result, &txid); err != nil {
		return "", fmt.Errorf("failed to unmarshal transaction ID: %v", err)
	}

	return txid, nil
}

// ListUnspent calls the listunspent RPC method
// walletName: name of the wallet to query
// minconf: minimum number of confirmations (default 1)
// maxconf: maximum number of confirmations (default 9999999)
// addresses: optional array of addresses to filter
// includeUnsafe: include outputs that are not safe to spend
// queryOptions: additional query options (amount filters, etc.)
func (c *Client) ListUnspent(walletName string, minconf, maxconf int, addresses []string, includeUnsafe bool, queryOptions map[string]interface{}) ([]UTXO, error) {
	// Prepare parameters
	var params []interface{}

	// Add minconf (default 1)
	if minconf == 0 {
		minconf = 1
	}
	params = append(params, minconf)

	// Add maxconf (default 9999999)
	if maxconf == 0 {
		maxconf = 9999999
	}
	params = append(params, maxconf)

	// Add addresses filter if provided
	if len(addresses) > 0 {
		params = append(params, addresses)
	} else {
		params = append(params, []string{}) // empty array
	}

	// Add includeUnsafe if specified
	if includeUnsafe {
		params = append(params, includeUnsafe)
	}

	// Add queryOptions if provided
	if len(queryOptions) > 0 {
		if !includeUnsafe {
			params = append(params, false) // default includeUnsafe
		}
		params = append(params, queryOptions)
	}

	// Call the RPC method with wallet endpoint
	resp, err := c.callWithWallet("listunspent", params, walletName)
	if err != nil {
		return nil, fmt.Errorf("failed to call listunspent: %v", err)
	}

	// Parse the result
	var utxos []UTXO
	if err := json.Unmarshal(resp.Result, &utxos); err != nil {
		return nil, fmt.Errorf("failed to unmarshal UTXO list: %v", err)
	}

	return utxos, nil
}

// GetTransaction calls the gettransaction RPC method
// walletName: name of the wallet containing the transaction
// txid: transaction ID to retrieve
// includeWatchonly: include watch-only addresses
// verbose: return detailed information
func (c *Client) GetTransaction(walletName, txid string, includeWatchonly, verbose bool) (*GetTransactionResponse, error) {
	// Prepare parameters
	params := []interface{}{txid}

	// Add includeWatchonly if specified
	if includeWatchonly {
		params = append(params, includeWatchonly)
	}

	// Add verbose if specified
	if verbose {
		if !includeWatchonly {
			params = append(params, false) // default includeWatchonly
		}
		params = append(params, verbose)
	}

	// Call the RPC method with wallet endpoint
	resp, err := c.callWithWallet("gettransaction", params, walletName)
	if err != nil {
		return nil, fmt.Errorf("failed to call gettransaction: %v", err)
	}

	// Parse the result
	var result GetTransactionResponse
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %v", err)
	}

	return &result, nil
}

// EstimateSmartFee calls the estimatesmartfee RPC method
// confTarget: confirmation target in blocks (between 1 - 1008)
// estimateMode: fee estimation mode ("UNSET", "ECONOMICAL", "CONSERVATIVE")
func (c *Client) EstimateSmartFee(confTarget int, estimateMode string) (*EstimateSmartFeeResponse, error) {
	// Validate confTarget
	if confTarget < 1 || confTarget > 1008 {
		return nil, fmt.Errorf("confirmation target must be between 1 and 1008, got: %d", confTarget)
	}

	// Prepare parameters
	params := []interface{}{confTarget}

	// Add estimate mode if provided
	if estimateMode != "" {
		params = append(params, estimateMode)
	}

	// Call the RPC method (no wallet endpoint needed)
	resp, err := c.call("estimatesmartfee", params)
	if err != nil {
		return nil, fmt.Errorf("failed to call estimatesmartfee: %v", err)
	}

	// Parse the result
	var result EstimateSmartFeeResponse
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fee estimate: %v", err)
	}

	return &result, nil
}

// GetWalletInfo calls the getwalletinfo RPC method
// walletName: name of the wallet to get info for
func (c *Client) GetWalletInfo(walletName string) (*GetWalletInfoResponse, error) {
	// Call the RPC method with wallet endpoint
	resp, err := c.callWithWallet("getwalletinfo", []interface{}{}, walletName)
	if err != nil {
		return nil, fmt.Errorf("failed to call getwalletinfo: %v", err)
	}

	// Parse the result
	var result GetWalletInfoResponse
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallet info: %v", err)
	}

	return &result, nil
}

// ListAddressGroupings calls the listaddressgroupings RPC method
// walletName: name of the wallet to list address groupings for
func (c *Client) ListAddressGroupings(walletName string) ([]AddressGrouping, error) {
	// Call the RPC method with wallet endpoint
	resp, err := c.callWithWallet("listaddressgroupings", []interface{}{}, walletName)
	if err != nil {
		return nil, fmt.Errorf("failed to call listaddressgroupings: %v", err)
	}

	// The response is a nested array structure: [[[address, amount, label], ...], ...]
	var rawResult [][][]interface{}
	if err := json.Unmarshal(resp.Result, &rawResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal address groupings raw: %v", err)
	}

	// Convert the raw result to our typed structure
	var groupings []AddressGrouping
	for _, rawGroup := range rawResult {
		var group AddressGrouping
		for _, rawAddr := range rawGroup {
			if len(rawAddr) >= 2 {
				addrInfo := AddressInfo{
					Address: rawAddr[0].(string),
					Amount:  rawAddr[1].(float64),
				}
				// Label is optional (third element)
				if len(rawAddr) >= 3 && rawAddr[2] != nil {
					addrInfo.Label = rawAddr[2].(string)
				}
				group = append(group, addrInfo)
			}
		}
		if len(group) > 0 {
			groupings = append(groupings, group)
		}
	}

	return groupings, nil
}
