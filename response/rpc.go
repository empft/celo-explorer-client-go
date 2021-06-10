package response

import "encoding/json"

type BaseResponse struct {
	Message string          `json:"message"`	
	Result  json.RawMessage `json:"result"`
	Status  string          `json:"status"`
}

type BaseEthResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
}

type EthResult struct {
	Result string `json:"result"`
}

type EthError struct {
	Error string `json:"error"`
}

type Balance string

type BalanceMulti struct {
	Account string `json:"account"`
	Balance string `json:"balance"`
	Stale   bool   `json:"stale"`
}

type PendingTxList struct {
	Contractaddress   string `json:"contractAddress"`
	Cumulativegasused string `json:"cumulativeGasUsed"`
	From              string `json:"from"`
	Gas               string `json:"gas"`
	Gasprice          string `json:"gasPrice"`
	Gasused           string `json:"gasUsed"`
	Hash              string `json:"hash"`
	Input             string `json:"input"`
	Nonce             string `json:"nonce"`
	To                string `json:"to"`
	Value             string `json:"value"`
}

type TxList struct {
	Blockhash           string `json:"blockHash"`
	Blocknumber         string `json:"blockNumber"`
	Confirmations       string `json:"confirmations"`
	Contractaddress     string `json:"contractAddress"`
	Cumulativegasused   string `json:"cumulativeGasUsed"`
	Feecurrency         string `json:"feeCurrency"`
	From                string `json:"from"`
	Gas                 string `json:"gas"`
	Gasprice            string `json:"gasPrice"`
	Gasused             string `json:"gasUsed"`
	Gatewayfee          string `json:"gatewayFee"`
	Gatewayfeerecipient string `json:"gatewayFeeRecipient"`
	Hash                string `json:"hash"`
	Input               string `json:"input"`
	Iserror             string `json:"isError"`
	Nonce               string `json:"nonce"`
	Timestamp           string `json:"timeStamp"`
	To                  string `json:"to"`
	Transactionindex    string `json:"transactionIndex"`
	TxreceiptStatus     string `json:"txreceipt_status"`
	Value               string `json:"value"`
}

type TxListInternal struct {
	Blocknumber     string `json:"blockNumber"`
	Contractaddress string `json:"contractAddress"`
	Errcode         string `json:"errCode"`
	From            string `json:"from"`
	Gas             string `json:"gas"`
	Gasused         string `json:"gasUsed"`
	Index           string `json:"index"`
	Input           string `json:"input"`
	Iserror         string `json:"isError"`
	Timestamp       string `json:"timeStamp"`
	To              string `json:"to"`
	Transactionhash string `json:"transactionHash"`
	Type            string `json:"type"`
	Value           string `json:"value"`
}

type TokenTx struct {
	Blockhash         string `json:"blockHash"`
	Blocknumber       string `json:"blockNumber"`
	Confirmations     string `json:"confirmations"`
	Contractaddress   string `json:"contractAddress"`
	Cumulativegasused string `json:"cumulativeGasUsed"`
	From              string `json:"from"`
	Gas               string `json:"gas"`
	Gasprice          string `json:"gasPrice"`
	Gasused           string `json:"gasUsed"`
	Hash              string `json:"hash"`
	Input             string `json:"input"`
	Logindex          string `json:"logIndex"`
	Nonce             string `json:"nonce"`
	Timestamp         string `json:"timeStamp"`
	To                string `json:"to"`
	Tokendecimal      string `json:"tokenDecimal"`
	Tokenname         string `json:"tokenName"`
	Tokensymbol       string `json:"tokenSymbol"`
	Transactionindex  string `json:"transactionIndex"`
	Value             string `json:"value"`
}

type TokenBalance string

type TokenList struct {
	Balance         string `json:"balance"`
	Contractaddress string `json:"contractAddress"`
	Decimals        string `json:"decimals"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	Type            string `json:"type"`
}

type GetMinedBlocks struct {
	Blocknumber string `json:"blockNumber"`
	Blockreward string `json:"blockReward"`
	Timestamp   string `json:"timeStamp"`
}

type ListAccounts struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

type GetLogs struct {
	Address             string   `json:"address"`
	Blocknumber         string   `json:"blockNumber"`
	Data                string   `json:"data"`
	Feecurrency         string   `json:"feeCurrency"`
	Gasprice            string   `json:"gasPrice"`
	Gasused             string   `json:"gasUsed"`
	Gatewayfee          string   `json:"gatewayFee"`
	Gatewayfeerecipient string   `json:"gatewayFeeRecipient"`
	Logindex            string   `json:"logIndex"`
	Timestamp           string   `json:"timeStamp"`
	Topics              []string `json:"topics"`
	Transactionhash     string   `json:"transactionHash"`
	Transactionindex    string   `json:"transactionIndex"`
}

type GetToken struct {
	Cataloged       bool   `json:"cataloged"`
	Contractaddress string `json:"contractAddress"`
	Decimals        string `json:"decimals"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	Totalsupply     string `json:"totalSupply"`
	Type            string `json:"type"`
}

type GetTokenHolders struct {
	Address string `json:"address"`
	Value   string `json:"value"`
}

type TokenSupply string

type EthSupplyExchange string

type EthSupply string

type CoinSupply float64

type EthPrice struct {
	Ethbtc          string `json:"ethbtc"`
	EthbtcTimestamp string `json:"ethbtc_timestamp"`
	Ethusd          string `json:"ethusd"`
	EthusdTimestamp string `json:"ethusd_timestamp"`
}

type TotalTransactions string

type GetBlockReward struct {
	Blockminer           string      `json:"blockMiner"`
	Blocknumber          string      `json:"blockNumber"`
	Blockreward          string      `json:"blockReward"`
	Timestamp            string      `json:"timeStamp"`
	Uncleinclusionreward interface{} `json:"uncleInclusionReward"`
	Uncles               interface{} `json:"uncles"`
}

type ListContracts struct {
	Abi              string `json:"ABI"`
	Compilerversion  string `json:"CompilerVersion"`
	Contractname     string `json:"ContractName"`
	Optimizationused string `json:"OptimizationUsed"`
	Sourcecode       string `json:"SourceCode"`
}

type GetAbi string

type GetSourceCode struct {
	Abi              string `json:"ABI"`
	Compilerversion  string `json:"CompilerVersion"`
	Contractname     string `json:"ContractName"`
	Optimizationused string `json:"OptimizationUsed"`
	Sourcecode       string `json:"SourceCode"`
}

type Verify struct {
	Abi              string `json:"ABI"`
	Compilerversion  string `json:"CompilerVersion"`
	Contractname     string `json:"ContractName"`
	Optimizationused string `json:"OptimizationUsed"`
	Sourcecode       string `json:"SourceCode"`
}

type GetTxInfo struct {
	Revertreason        string `json:"revertReason"`
	Blocknumber         string `json:"blockNumber"`
	Confirmations       string `json:"confirmations"`
	Feecurrency         string `json:"feeCurrency"`
	From                string `json:"from"`
	Gaslimit            string `json:"gasLimit"`
	Gasprice            string `json:"gasPrice"`
	Gasused             string `json:"gasUsed"`
	Gatewayfee          string `json:"gatewayFee"`
	Gatewayfeerecipient string `json:"gatewayFeeRecipient"`
	Hash                string `json:"hash"`
	Input               string `json:"input"`
	Logs                []struct {
		Address string   `json:"address"`
		Data    string   `json:"data"`
		Index	string	 `json:"index"`
		Topics  []string `json:"topics"`
	} `json:"logs"`
	Success   bool   `json:"success"`
	Timestamp string `json:"timeStamp"`
	To        string `json:"to"`
	Value     string `json:"value"`
}

type GetTxReceiptStatus struct {
	Status string `json:"status"`
}

type GetStatus struct {
	Errdescription string `json:"errDescription"`
	Iserror        string `json:"isError"`
}