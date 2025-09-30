package btcrpc

import (
	"encoding/json"
	"fmt"
)

// GetBlockchainInfo calls the getblockchaininfo RPC method
func (c *Client) GetBlockchainInfo() (*BlockchainInfo, error) {
	// Call the RPC method
	resp, err := c.call("getblockchaininfo", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to call getblockchaininfo: %v", err)
	}

	// Parse the result
	var info BlockchainInfo
	if err := json.Unmarshal(resp.Result, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal blockchain info: %v", err)
	}

	return &info, nil
}

// GetNetworkInfo calls the getnetworkinfo RPC method
func (c *Client) GetNetworkInfo() (*NetworkInfo, error) {
	// Call the RPC method
	resp, err := c.call("getnetworkinfo", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to call getnetworkinfo: %v", err)
	}

	// Parse the result
	var info NetworkInfo
	if err := json.Unmarshal(resp.Result, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal network info: %v", err)
	}

	return &info, nil
}

// GenerateToAddress calls the generatetoaddress RPC method
// This method is used in regtest mode to generate blocks to a specific address
// nblocks: number of blocks to generate
// address: address to receive the block rewards
// maxtries: maximum number of iterations to try (optional, default 1000000)
func (c *Client) GenerateToAddress(nblocks int, address string, maxtries *int) (GenerateToAddressResponse, error) {
	// Prepare parameters
	params := []interface{}{nblocks, address}

	// Add maxtries if provided
	if maxtries != nil {
		params = append(params, *maxtries)
	}

	// Call the RPC method
	resp, err := c.call("generatetoaddress", params)
	if err != nil {
		return nil, fmt.Errorf("failed to call generatetoaddress: %v", err)
	}

	// Parse the result (generatetoaddress returns an array of block hashes)
	var blockHashes GenerateToAddressResponse
	if err := json.Unmarshal(resp.Result, &blockHashes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal generate to address response: %v", err)
	}

	return blockHashes, nil
}

// GetBlock calls the getblock RPC method
// blockhash: hash of the block to retrieve
// verbosity: 0=raw hex, 1=json object, 2=json object with transaction data
func (c *Client) GetBlock(blockhash string, verbosity int) (*GetBlockResponse, error) {
	// Prepare parameters
	params := []interface{}{blockhash}

	// Add verbosity if specified (default is 1)
	if verbosity != 1 {
		params = append(params, verbosity)
	}

	// Call the RPC method
	resp, err := c.call("getblock", params)
	if err != nil {
		return nil, fmt.Errorf("failed to call getblock: %v", err)
	}

	// For verbosity 0, return raw hex (not implemented in this struct)
	if verbosity == 0 {
		return nil, fmt.Errorf("verbosity 0 (raw hex) not supported, use verbosity 1 or 2")
	}

	// Parse the result
	var block GetBlockResponse
	if err := json.Unmarshal(resp.Result, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block info: %v", err)
	}

	return &block, nil
}

// GetBlockHash calls the getblockhash RPC method
// height: block height to get hash for
func (c *Client) GetBlockHash(height int) (string, error) {
	// Prepare parameters
	params := []interface{}{height}

	// Call the RPC method
	resp, err := c.call("getblockhash", params)
	if err != nil {
		return "", fmt.Errorf("failed to call getblockhash: %v", err)
	}

	// Parse the result
	var blockHash GetBlockHashResponse
	if err := json.Unmarshal(resp.Result, &blockHash); err != nil {
		return "", fmt.Errorf("failed to unmarshal block hash: %v", err)
	}

	return string(blockHash), nil
}

// GetRawTransaction calls the getrawtransaction RPC method
// txid: transaction ID to retrieve
// verbose: if false, return hex string; if true, return JSON object
// blockhash: optional block hash to look in (for performance)
func (c *Client) GetRawTransaction(txid string, verbose bool, blockhash *string) (*GetRawTransactionResponse, error) {
	// Prepare parameters
	params := []interface{}{txid, verbose}

	// Add blockhash if provided
	if blockhash != nil {
		params = append(params, *blockhash)
	}

	// Call the RPC method
	resp, err := c.call("getrawtransaction", params)
	if err != nil {
		return nil, fmt.Errorf("failed to call getrawtransaction: %v", err)
	}

	// For non-verbose mode, return hex string in Hex field
	if !verbose {
		var hexString string
		if err := json.Unmarshal(resp.Result, &hexString); err != nil {
			return nil, fmt.Errorf("failed to unmarshal raw transaction hex: %v", err)
		}
		return &GetRawTransactionResponse{Hex: hexString}, nil
	}

	// For verbose mode, parse the full JSON object
	var tx GetRawTransactionResponse
	if err := json.Unmarshal(resp.Result, &tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal raw transaction: %v", err)
	}

	return &tx, nil
}

// GetMempoolInfo calls the getmempoolinfo RPC method
func (c *Client) GetMempoolInfo() (*GetMempoolInfoResponse, error) {
	// Call the RPC method
	resp, err := c.call("getmempoolinfo", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to call getmempoolinfo: %v", err)
	}

	// Parse the result
	var info GetMempoolInfoResponse
	if err := json.Unmarshal(resp.Result, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal mempool info: %v", err)
	}

	return &info, nil
}
