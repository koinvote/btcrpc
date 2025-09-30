package btcrpc

import (
	"encoding/json"
	"fmt"
)

// === Raw Transaction Functions ===

// CreateRawTransaction creates a new raw transaction spending the given inputs and creating new outputs.
func (c *Client) CreateRawTransaction(inputs []CreateRawTransactionInput, outputs map[string]interface{}, locktime int64, replaceable bool) (string, error) {
	params := []interface{}{inputs, outputs}

	if locktime != 0 {
		params = append(params, locktime)
	}

	if replaceable {
		params = append(params, replaceable)
	}

	resp, err := c.call("createrawtransaction", params)
	if err != nil {
		return "", fmt.Errorf("createrawtransaction RPC call failed: %w", err)
	}

	if resp.Error != nil {
		return "", fmt.Errorf("createrawtransaction RPC error: %s", resp.Error.Message)
	}

	var result string
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal createrawtransaction result: %w", err)
	}

	return result, nil
}

// SignRawTransactionWithWallet signs inputs for raw transaction (serialized, hex-encoded).
func (c *Client) SignRawTransactionWithWallet(walletName, hexstring string, prevtxs []interface{}, sighashtype string) (*SignRawTransactionResponse, error) {
	params := []interface{}{hexstring}

	if len(prevtxs) > 0 {
		params = append(params, prevtxs)
	} else {
		params = append(params, nil)
	}

	if sighashtype != "" {
		params = append(params, sighashtype)
	}

	resp, err := c.callWithWallet("signrawtransactionwithwallet", params, walletName)
	if err != nil {
		return nil, fmt.Errorf("signrawtransactionwithwallet RPC call failed: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("signrawtransactionwithwallet RPC error: %s", resp.Error.Message)
	}

	var result SignRawTransactionResponse
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal signrawtransactionwithwallet result: %w", err)
	}

	return &result, nil
}

// SendRawTransaction submits raw transaction (serialized, hex-encoded) to local node and network.
func (c *Client) SendRawTransaction(hexstring string, maxfeerate float64) (string, error) {
	params := []interface{}{hexstring}

	if maxfeerate > 0 {
		params = append(params, maxfeerate)
	}

	resp, err := c.call("sendrawtransaction", params)
	if err != nil {
		return "", fmt.Errorf("sendrawtransaction RPC call failed: %w", err)
	}

	if resp.Error != nil {
		return "", fmt.Errorf("sendrawtransaction RPC error: %s", resp.Error.Message)
	}

	var result string
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal sendrawtransaction result: %w", err)
	}

	return result, nil
}

// === Private Key Management Functions ===

// DumpPrivKey reveals the private key corresponding to address.
func (c *Client) DumpPrivKey(walletName, address string) (string, error) {
	params := []interface{}{address}

	resp, err := c.callWithWallet("dumpprivkey", params, walletName)
	if err != nil {
		return "", fmt.Errorf("dumpprivkey RPC call failed: %w", err)
	}

	if resp.Error != nil {
		return "", fmt.Errorf("dumpprivkey RPC error: %s", resp.Error.Message)
	}

	var result string
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal dumpprivkey result: %w", err)
	}

	return result, nil
}

// ImportPrivKey adds a private key (as returned by dumpprivkey) to your wallet.
func (c *Client) ImportPrivKey(walletName, privkey, label string, rescan bool) error {
	params := []interface{}{privkey}

	if label != "" {
		params = append(params, label)
	} else {
		params = append(params, "")
	}

	params = append(params, rescan)

	resp, err := c.callWithWallet("importprivkey", params, walletName)
	if err != nil {
		return fmt.Errorf("importprivkey RPC call failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("importprivkey RPC error: %s", resp.Error.Message)
	}

	return nil
}

// === Multisignature Functions ===

// CreateMultisig creates a multi-signature address with n signature of m keys required.
func (c *Client) CreateMultisig(nrequired int, keys []string, addressType string) (*CreateMultisigResponse, error) {
	params := []interface{}{nrequired, keys}

	if addressType != "" {
		params = append(params, addressType)
	}

	resp, err := c.call("createmultisig", params)
	if err != nil {
		return nil, fmt.Errorf("createmultisig RPC call failed: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("createmultisig RPC error: %s", resp.Error.Message)
	}

	var result CreateMultisigResponse
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal createmultisig result: %w", err)
	}

	return &result, nil
}

// AddMultisigAddress adds a nrequired-to-sign multisignature address to the wallet.
func (c *Client) AddMultisigAddress(walletName string, nrequired int, keys []string, label, addressType string) (*AddMultisigAddressResponse, error) {
	params := []interface{}{nrequired, keys}

	if label != "" {
		params = append(params, label)
	} else {
		params = append(params, "")
	}

	if addressType != "" {
		params = append(params, addressType)
	}

	resp, err := c.callWithWallet("addmultisigaddress", params, walletName)
	if err != nil {
		return nil, fmt.Errorf("addmultisigaddress RPC call failed: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("addmultisigaddress RPC error: %s", resp.Error.Message)
	}

	var result AddMultisigAddressResponse
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal addmultisigaddress result: %w", err)
	}

	return &result, nil
}

// === Mempool Functions ===

// GetRawMempool returns all transaction ids in memory pool as array.
func (c *Client) GetRawMempool(verbose bool, mempoolSequence bool) (interface{}, error) {
	params := []interface{}{verbose}

	if mempoolSequence {
		params = append(params, mempoolSequence)
	}

	resp, err := c.call("getrawmempool", params)
	if err != nil {
		return nil, fmt.Errorf("getrawmempool RPC call failed: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("getrawmempool RPC error: %s", resp.Error.Message)
	}

	if verbose {
		var result GetRawMempoolVerboseResponse
		if err := json.Unmarshal(resp.Result, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal verbose getrawmempool result: %w", err)
		}
		return result, nil
	} else {
		var result GetRawMempoolResponse
		if err := json.Unmarshal(resp.Result, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal getrawmempool result: %w", err)
		}
		return result, nil
	}
}

// GetRawMempoolSimple returns all transaction ids in memory pool as a slice of strings.
func (c *Client) GetRawMempoolSimple() ([]string, error) {
	result, err := c.GetRawMempool(false, false)
	if err != nil {
		return nil, err
	}

	if txids, ok := result.(GetRawMempoolResponse); ok {
		return []string(txids), nil
	}

	return nil, fmt.Errorf("unexpected response type from getrawmempool")
}

// GetRawMempoolVerbose returns detailed information about all transactions in memory pool.
func (c *Client) GetRawMempoolVerbose() (map[string]GetRawMempoolEntry, error) {
	result, err := c.GetRawMempool(true, false)
	if err != nil {
		return nil, err
	}

	if verboseResult, ok := result.(GetRawMempoolVerboseResponse); ok {
		return map[string]GetRawMempoolEntry(verboseResult), nil
	}

	return nil, fmt.Errorf("unexpected response type from verbose getrawmempool")
}
