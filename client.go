package btcrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client represents a Bitcoin Core RPC client
type Client struct {
	url      string
	username string
	password string
	client   *http.Client
}

// NewClient creates a new Bitcoin Core RPC client
func NewClient(url, username, password string) *Client {
	return &Client{
		url:      url,
		username: username,
		password: password,
		client:   &http.Client{},
	}
}

// call performs a JSON-RPC call to Bitcoin Core
func (c *Client) call(method string, params []interface{}) (*RPCResponse, error) {
	return c.callWithWallet(method, params, "")
}

// callWithWallet performs a JSON-RPC call to Bitcoin Core with specific wallet
func (c *Client) callWithWallet(method string, params []interface{}, walletName string) (*RPCResponse, error) {
	// Create RPC request
	rpcReq := RPCRequest{
		Method:  method,
		Params:  params,
		ID:      1,
		JsonRPC: "1.0",
	}

	// Serialize request to JSON
	reqBody, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create URL with wallet endpoint if needed
	url := c.url
	if walletName != "" {
		url = c.url + "/wallet/" + walletName
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.username, c.password)

	// Perform HTTP request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s, body: %s", resp.Status, string(body))
	}

	// Parse JSON-RPC response
	var rpcResp RPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// Check for RPC error
	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	return &rpcResp, nil
}
