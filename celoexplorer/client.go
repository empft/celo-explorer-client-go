// Implements a client for celo explorer api through rpc http.
// Not all returned data will be converted to simpler, programmer-friendly data structure because they are not used.
// This code has not been tested. Use it at your own risk.
package celoexplorer

import (
	"encoding/hex"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.com/stevealexrs/celo-explorer-client-go/request"
)

const (
	BaseUrl        string = "https://explorer.celo.org/api"
	TestnetBaseUrl string = "https://alfajores-blockscout.celo-testnet.org/api"

	// Celo contract address
	CeloGold string = "471ece3750da237f93b8e339c536989b8978a438"
	// CeloUSD contract address
	CeloUSD  string = "765de816845861e75a25fca122bb6898b8b1282a"
	// CeloEUR contract address
	CeloEUR  string = "d8763cba276a3738e6de85b4b3bf5fded6d6ca73"
	
	// Alfajores Testnet Celo contract address
	TestnetCeloGold string = "F194afDf50B03e69Bd7D057c1Aa9e10c9954E4C9"
	// Alfajores Testnet CeloUSD contract address
	TestnetCeloUSD string = "874069Fa1Eb16D44d622F2e0Ca25eeA172369bC1"
	// Alfajores Testnet CeloEUR contract address
	TestnetCeloEUR string = "10c892A6EC43a53E45D0B916B4b7D383B1b78C0F"
)
type Client struct {
	req  *request.RequestClient
}

func New(url string) *Client {
	tr := &http.Transport{
		MaxIdleConns:    100,
		IdleConnTimeout: 30 * time.Second,
	}
	httpClient := &http.Client{Transport: tr}
	reqClient := request.NewRequestClientWithHttp(url, httpClient)
	return &Client{
		req: reqClient,
	}
}

// errors are ignored
func toBigInt(value string, base int) *big.Int {
	n := new(big.Int)
	n, _ = n.SetString(value, base)
	return n
}

func trim0x(s string) string {
	return strings.TrimPrefix(s, "0x")
}

// errors are ignored
func hexToByte(s string) []byte {
	src := []byte(s)
	dst := make([]byte, hex.DecodedLen(len(src)))
	n, _ := hex.Decode(dst, src)
	return dst[:n]
}

// Mimics Ethereum JSON RPC's eth_getBalance.
// Returns the wei balance (1 Celo = 10^18 wei) for an address as of the provided block (defaults to latest).
func (c *Client) EthGetBalance(address string, block *big.Int) (*big.Int, error) {
	bal, err := c.req.EthGetBalance(address, block)
	if err != nil {
		return nil, err
	}

	b := toBigInt(trim0x(bal), 16)
	return b, nil
}

// Get balance for address.
func (c *Client) Balance(address string) (*big.Int, error) {
	bal, err := c.req.Balance(address)
	if err != nil {
		return nil, err
	}
	
	b := toBigInt(string(bal), 10)
	return b, nil
}

type FetchedBalance struct {
	Address string
	Balance *big.Int
	Stale   bool
}

// Get balance for multiple addresses.
// If the balance hasn't been updated in a long time, we will double check with the node to fetch the absolute latest balance. This will not be reflected in the current request, but once it is updated, subsequent requests will show the updated balance. You can know that this is taking place via the `stale` attribute, which is set to `true` if a new balance is being fetched.
func (c *Client) BalanceMulti(address []string) ([]FetchedBalance, error) {
	bal, err := c.req.BalanceMulti(address)
	if err != nil {
		return nil, err
	}

	result := make([]FetchedBalance, len(bal))
	for i, v := range bal {
		result[i].Address = trim0x(v.Account)
		result[i].Balance = toBigInt(v.Balance, 10)
		result[i].Stale = v.Stale
	}
	return result, nil
}

type Transaction struct {
	BlockHash           string
	BlockNumber         *big.Int
	Confirmations       *big.Int
	Contractaddress     string
	CumulativeGasUsed   int
	Feecurrency         string
	From                string
	Gas                 int
	GasPrice            *big.Int
	GasUsed             int
	GatewayFee          int
	GatewayFeeRecipient string
	Hash                string
	Input               []byte
	IsError             bool
	Nonce               int
	Timestamp           time.Time
	To                  string
	TransactionIndex    int
	TxReceiptStatus     bool
	Value               *big.Int
}

// Get transactions sent by an address. Up to a maximum of 10,000 transactions.
func (c *Client) TxList(address string, sort *request.SortDirection, block *request.BlockRange, page *request.PageRange, filter *request.FilterDirection, timeRange *request.TimeRange) ([]Transaction, error) {
	txList, err := c.req.TxList(address, sort, block, page, filter, timeRange)
	if err != nil {
		return nil, err
	}

	transactions := make([]Transaction, len(txList))
	for i, v := range txList {
		transactions[i].BlockHash = trim0x(v.Blockhash)
		transactions[i].BlockNumber = toBigInt(v.Blocknumber, 10)
		transactions[i].Confirmations = toBigInt(v.Confirmations, 10)
		transactions[i].Contractaddress = trim0x(v.Contractaddress)
		transactions[i].CumulativeGasUsed, _ = strconv.Atoi(v.Cumulativegasused)

		transactions[i].Feecurrency = v.Feecurrency
		transactions[i].From = trim0x(v.From)
		transactions[i].Gas, _ = strconv.Atoi(v.Gas)

		transactions[i].GasPrice = toBigInt(v.Gasprice, 10)
		transactions[i].GasUsed, _ = strconv.Atoi(v.Gasused)

		transactions[i].GatewayFee, _ = strconv.Atoi(v.Gatewayfee)

		transactions[i].GatewayFeeRecipient = trim0x(v.Gatewayfeerecipient)
		transactions[i].Hash = trim0x(v.Hash)
		transactions[i].Input = hexToByte(trim0x(v.Input))

		if v.Iserror == "0" {
			transactions[i].IsError = false
		} else {
			transactions[i].IsError = true
		}

		transactions[i].Nonce, _ = strconv.Atoi(v.Nonce)

		iTime, _ := strconv.ParseInt(v.Timestamp, 10, 64)
		transactions[i].Timestamp = time.Unix(iTime, 0)

		transactions[i].To = trim0x(v.To)
		transactions[i].TransactionIndex, _ = strconv.Atoi(v.Transactionindex)

		if v.TxreceiptStatus == "1" {
			transactions[i].TxReceiptStatus = true
		} else {
			transactions[i].TxReceiptStatus = false
		}

		transactions[i].Value = toBigInt(v.Value, 10)
	}

	return transactions, nil
}

type TokenTransfer struct {
	Value *big.Int
	BlockHash string
	BlockNumber *big.Int
	Confirmations *big.Int
	ContractAddress string
	CumulativeGasUsed int
	From string
	To string
	Gas int
	Gasprice *big.Int
	Gasused int
	Hash string
	Input []byte
	LogIndex int
	Nonce int
	Timestamp time.Time
	TokenDecimal int
	TokenName string
	TokenSymbol string
	TransactionIndex int
}

// Get token transfer events to and from an address.
func (c *Client) TokenTx(address string, contractAddress *string, sort *request.SortDirection, block *request.BlockRange, page *request.PageRange) ([]TokenTransfer, error) {
	tokensList, err := c.req.TokenTx(address, contractAddress, sort, block, page)
	if err != nil {
		return nil, err
	}

	tokens := make([]TokenTransfer, len(tokensList))
	for i, v := range tokensList {
		tokens[i].Value = toBigInt(v.Value, 10)
		tokens[i].BlockHash = trim0x(v.Blockhash)
		tokens[i].BlockNumber = toBigInt(v.Blocknumber, 10)
		tokens[i].Confirmations = toBigInt(v.Confirmations, 10)
		tokens[i].ContractAddress = trim0x(v.Contractaddress)
		tokens[i].CumulativeGasUsed, _ = strconv.Atoi(v.Cumulativegasused)

		tokens[i].From = trim0x(v.From)
		tokens[i].To = trim0x(v.To)

		tokens[i].Gas, _ = strconv.Atoi(v.Gas)
		tokens[i].Gasprice = toBigInt(v.Gasprice, 10)
		tokens[i].Gasused, _ = strconv.Atoi(v.Gasused)
	
		tokens[i].Hash = trim0x(v.Hash)
		tokens[i].Input = hexToByte(trim0x(v.Input))
		tokens[i].LogIndex, _ = strconv.Atoi(v.Logindex) 

		tokens[i].Nonce, _ = strconv.Atoi(v.Nonce) 

		iTime, _ := strconv.ParseInt(v.Timestamp, 10, 64)
		tokens[i].Timestamp = time.Unix(iTime, 0)

		tokens[i].TokenDecimal, _ = strconv.Atoi(v.Tokendecimal) 
		tokens[i].TokenName = v.Tokenname
		tokens[i].TokenSymbol = v.Tokensymbol

		tokens[i].TransactionIndex, _ = strconv.Atoi(v.Transactionindex) 
	}
	return tokens, nil
}

// Get token account balance for token contract address.
func (c *Client) TokenBalance(contractAddress, address string) (*big.Int, error) {
	bal, err := c.req.TokenBalance(contractAddress, address)
	if err != nil {
		return nil, err
	}

	return toBigInt(string(bal), 10), nil
}

type Token struct {
	Balance *big.Int
	ContractAddress string
	Decimals int
	Name string
	Symbol string
	Type string
}

// Get list of tokens owned by address.
func (c *Client) TokenList(address string) ([]Token, error) {
	tokenList, err := c.req.TokenList(address)
	if err != nil {
		return nil, err
	}

	tokens := make([]Token, len(tokenList))
	for i, v := range tokenList {
		tokens[i].Balance = toBigInt(v.Balance, 10)
		tokens[i].ContractAddress = trim0x(v.Contractaddress)
		tokens[i].Decimals, _ = strconv.Atoi(v.Decimals)
		tokens[i].Name = v.Name
		tokens[i].Symbol = v.Symbol
		tokens[i].Type = v.Type
	}
	return tokens, nil
}

type EventLog struct {
	Address string
	BlockNumber *big.Int
	Data string
	// address of fee currency
	FeeCurrency string
	GasPrice *big.Int
	GasUsed int
	GatewayFee *big.Int
	GatewayfeeRecipient string
	LogIndex int
	Timestamp time.Time
	Topics []string
	TransactionHash string
	TransactionIndex int
}

// WARNING: This function may not work correctly since I am not sure whether the returned data is in hex or decimal form.
// Get event logs for an address and/or topics. Up to a maximum of 1,000 event logs.
func (c *Client) GetLogs(block request.BlockRangeAdv, contractAddress string, topics request.Topics) ([]EventLog, error) {
	logList, err := c.req.GetLogs(block, contractAddress, topics)
	if err != nil {
		return nil ,err
	}

	logs := make([]EventLog, len(logList))
	for i, v := range logList {
		logs[i].Address = trim0x(v.Address)
		logs[i].BlockNumber = toBigInt(trim0x(v.Blocknumber), 16)
		logs[i].Data = trim0x(v.Data)
		logs[i].FeeCurrency = trim0x(v.Feecurrency)
		logs[i].GasPrice = toBigInt(trim0x(v.Gasprice), 16)

		gasUsed, _ := strconv.ParseInt(trim0x(v.Gasused), 16, 64)
		logs[i].GasUsed = int(gasUsed)

		logs[i].GatewayFee = toBigInt(v.Gatewayfee, 10)
		logs[i].GatewayfeeRecipient = trim0x(v.Gatewayfeerecipient)
		lIndex, _  := strconv.ParseInt(trim0x(v.Logindex), 16, 64)
		logs[i].LogIndex = int(lIndex)

		iTime, _ := strconv.ParseInt(trim0x(v.Timestamp), 16, 64)
		logs[i].Timestamp = time.Unix(iTime, 0)

		array := make([]string, len(v.Topics))
		for i, v := range v.Topics {
			array[i] = trim0x(v)
		}
		logs[i].Topics = array

		logs[i].TransactionHash = trim0x(v.Transactionhash)

		tIndex, _ := strconv.ParseInt(trim0x(v.Transactionindex), 16, 64)
		logs[i].TransactionIndex = int(tIndex)
	}

	return logs, nil
}

type TokenInfo struct {
	Catalogued bool
	ContractAddress string
	Decimals int
	Name string
	Symbol string
	TotalSupply *big.Int
	Type string
}

// Get ERC-20 or ERC-721 token by contract address.
func (c *Client) GetToken(contractAddress string) (TokenInfo, error) {
	info, err := c.req.GetToken(contractAddress)
	if err != nil {
		return TokenInfo{}, err
	}

	deci, _ := strconv.Atoi(info.Decimals)
	return TokenInfo{
		Catalogued:      info.Cataloged,
		ContractAddress: trim0x(info.Contractaddress),
		Decimals:        deci,
		Name:            info.Name,
		Symbol:          info.Symbol,
		TotalSupply:     toBigInt(info.Totalsupply, 10),
		Type:            info.Type,
	}, nil
}

type TxLog struct {
	Address string
	Data 	[]byte
	Index	int
	Topics	[]string
}

type TransactionWithLogs struct {
	BlockNumber         *big.Int
	Confirmations       *big.Int
	Feecurrency         string
	From                string
	GasLimit            *big.Int
	GasPrice            *big.Int
	GasUsed             int
	GatewayFee          *big.Int
	GatewayFeeRecipient string
	Hash                string
	Input               []byte
	Logs				[]TxLog
	RevertReason		string
	Success				bool
	Timestamp           time.Time
	To                  string
	Value               *big.Int
}

// Get transaction info.
func (c *Client) GetTxInfo(txHash string) (TransactionWithLogs, error) {
	txInfo, err := c.req.GetTxInfo(txHash, nil)
	if err != nil {
		return TransactionWithLogs{}, err
	}

	gasUsed, _ := strconv.Atoi(txInfo.Gasused)

	iTime, _ := strconv.ParseInt(trim0x(txInfo.Timestamp), 10, 64)
	timestamp := time.Unix(iTime, 0)


	logs := make([]TxLog, len(txInfo.Logs))
	for i, v := range txInfo.Logs {
		logs[i].Address = trim0x(v.Address)
		logs[i].Data = []byte(trim0x(v.Data))

		index, _ := strconv.Atoi(v.Index)
		logs[i].Index = index

		array := make([]string, len(v.Topics))
		for i, v := range v.Topics {
			array[i] = trim0x(v)
		}
		logs[i].Topics = array
	}


	return TransactionWithLogs{
		BlockNumber:         toBigInt(txInfo.Blocknumber, 10),
		Confirmations:       toBigInt(txInfo.Confirmations, 10),
		Feecurrency:         trim0x(txInfo.Feecurrency),
		From:                trim0x(txInfo.From),
		GasLimit:            toBigInt(txInfo.Gaslimit, 10),
		GasPrice:            toBigInt(txInfo.Gasused, 10),
		GasUsed:             gasUsed,
		GatewayFee:          toBigInt(txInfo.Gatewayfee, 10),
		GatewayFeeRecipient: trim0x(txInfo.Gatewayfeerecipient),
		Hash:                trim0x(txInfo.Hash),
		Input:               []byte(trim0x(txInfo.Input)),
		Logs:                logs,
		RevertReason:        txInfo.Revertreason,
		Success:             txInfo.Success,
		Timestamp:           timestamp,
		To:                  trim0x(txInfo.To),
		Value:               toBigInt(txInfo.Value, 10),
	}, nil
}

// Get transaction receipt status. 
func (c *Client) GetTxReceiptStatus(txHash string) (bool, error) {
	status, err := c.req.GetTxReceiptStatus(txHash)
	if err != nil {
		return false, err
	}

	if status.Status == "1" {
		return true, nil
	}
	return false, nil
}

// Get error status and error message. 
func (c *Client) GetStatus(txHash string) (bool, string, error) {
	status, err := c.req.GetStatus(txHash)
	if err != nil {
		return false, "", err
	}

	if status.Iserror == "0" {
		return true, status.Errdescription, nil
	}
	return false, status.Errdescription, nil
}