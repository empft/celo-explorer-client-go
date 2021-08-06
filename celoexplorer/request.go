package celoexplorer

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	ethGetBalanceUrl      string = "?module=account&action=eth_get_balance&address={addressHash}"
	balanceUrl            string = "?module=account&action=balance&address={addressHash}"
	balanceMultiUrl       string = "?module=account&action=balancemulti&address={addressHash1,addressHash2,addressHash3}"
	pendingTxListUrl      string = "?module=account&action=pendingtxlist&address={addressHash}"
	txListUrl             string = "?module=account&action=txlist&address={addressHash}"
	txListInternalUrl     string = "?module=account&action=txlistinternal&txhash={transactionHash}"
	tokenTxUrl            string = "?module=account&action=tokentx&address={addressHash}"
	tokenBalanceUrl       string = "?module=account&action=tokenbalance&contractaddress={contractAddressHash}&address={addressHash}"
	tokenListUrl          string = "?module=account&action=tokenlist&address={addressHash}"
	getMinedBlocksUrl     string = "?module=account&action=getminedblocks&address={addressHash}"
	listAccountsUrl       string = "?module=account&action=listaccounts"
	getLogsUrl            string = "?module=logs&action=getLogs&fromBlock={blockNumber}&toBlock={blockNumber}&address={addressHash}&topic0={firstTopic}"
	getTokenUrl           string = "?module=token&action=getToken&contractaddress={contractAddressHash}"
	getTokenHoldersUrl    string = "?module=token&action=getTokenHolders&contractaddress={contractAddressHash}"
	tokenSupplyUrl        string = "?module=stats&action=tokensupply&contractaddress={contractAddressHash}"
	ethSupplyExchangeUrl  string = "?module=stats&action=ethsupplyexchange"
	ethSupplyUrl          string = "?module=stats&action=ethsupply"
	coinSupplyUrl         string = "?module=stats&action=coinsupply"
	ethPriceUrl           string = "?module=stats&action=ethprice"
	totalTransactionsUrl  string = "?module=stats&action=totaltransactions"
	getBlockRewardUrl     string = "?module=block&action=getblockreward&blockno={blockNumber}"
	ethBlockNumberUrl     string = "?module=block&action=eth_block_number"
	listContractsUrl      string = "?module=contract&action=listcontracts"
	getAbiUrl             string = "?module=contract&action=getabi&address={addressHash}"
	getSourceCodeUrl      string = "?module=contract&action=getsourcecode&address={addressHash}"
	verifyUrl             string = "?module=contract&action=verify&addressHash={addressHash}&name={name}&compilerVersion={compilerVersion}&optimization={false}&contractSourceCode={contractSourceCode}"
	getTxInfoUrl          string = "?module=transaction&action=gettxinfo&txhash={transactionHash}"
	getTxReceiptStatusUrl string = "?module=transaction&action=gettxreceiptstatus&txhash={transactionHash}"
	getStatusUrl          string = "?module=transaction&action=getstatus&txhash={transactionHash}"
)

type RequestClient struct {
	http *http.Client
	base string
}

func NewRequestClientWithHttp(url string, http *http.Client) *RequestClient {
	return &RequestClient{
		http: http,
		base: url,
	}
}

func add0x(s string) string {
	var sb strings.Builder
	sb.WriteString("0x")
	sb.WriteString(s)
	return sb.String()
}

func buildUrl(ss ...string) *url.URL {
	var sb strings.Builder
	for _, s := range ss {
		sb.WriteString(s)
	}

	u, err := url.Parse(sb.String())
	if err != nil {
		panic(err)
	}
	return u
}

func (r *RequestClient) jsonResponse(u *url.URL, respObject interface{}) error {
	resp, err := r.http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var baseResp BaseResponse
	json.Unmarshal(body, &baseResp)

	if baseResp.Status == "1" {
		json.Unmarshal(baseResp.Result, respObject)
	} else {
		return errors.New(baseResp.Message)
	}
	return nil
}

// Use different json parser for different response code
// return true if success
func (r *RequestClient) jsonResponseDiff(u *url.URL, respSuccess, respFailure interface{}) (bool, error) {
	resp, err := r.http.Get(u.String())
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		json.Unmarshal(body, respSuccess)
		return true, nil
	}

	json.Unmarshal(body, respFailure)
	return false, nil
}

type queryBuilder struct {
	url *url.URL
}

func newQueryBuilder(u *url.URL) *queryBuilder {
	return &queryBuilder{
		url: u,
	}
}

func (qb *queryBuilder) set(key, value string) {
	q := qb.url.Query()
	q.Set(key, value)
	qb.url.RawQuery = q.Encode()
}

// set after nil check
func (qb *queryBuilder) setIfExist(key string, value *string) {
	if value != nil {
		qb.set(key, *value)
	}
}

func (qb *queryBuilder) address(address string) {
	qb.set("address", add0x(address))
}

func (qb *queryBuilder) txHash(hash string) {
	qb.set("txhash", add0x(hash))
}

func (qb *queryBuilder) contractAddress(address string) {
	qb.set("contractaddress", add0x(address))
}


func (qb *queryBuilder) addressMulti(address []string) {
	formatAddress := make([]string, len(address))
	for i, v := range address {
		formatAddress[i] = add0x(v)
	}
	qb.set("address", strings.Join(formatAddress, ","))
}

func (qb *queryBuilder) block(number *big.Int) {
	if number != nil {
		qb.set("block", number.String())
	}
}

func (qb *queryBuilder) blockNo(number *big.Int) {
	if number != nil {
		qb.set("blockno", number.String())
	}
}

func (qb *queryBuilder) pageRange(pages *PageRange) {
	if pages != nil {
		qb.set("page", strconv.Itoa(pages.Page))
		qb.set("offset", strconv.Itoa(pages.Offset))
	}
}

func (qb *queryBuilder) sort(direction *sortDirection) {
	if direction != nil {
		qb.set("sort", string(*direction))
	}
}

func (qb *queryBuilder) blockRange(block *BlockRange) {
	if block != nil {
		if block.StartBlock != nil {
			qb.set("startblock", block.StartBlock.String())
		}
		if block.EndBlock != nil {
			qb.set("endblock", block.EndBlock.String())
		}
	}
}

func (qb *queryBuilder) filterByDirection(filter *filterDirection) {
	if filter != nil {
		qb.set("filterby", string(*filter))
	}
}

func (qb *queryBuilder) filterContract(filter *filterContract) {
	if filter != nil {
		qb.set("filter", string(*filter))
	}
}

func (qb *queryBuilder) notDecompiledWithVersion(version *string) {
	qb.setIfExist("not_decompiled_with_version", version)
}

func (qb *queryBuilder) ignoreProxy(ignore *bool) {
	if ignore != nil {
		var text string
		if *ignore {
			text = "true"
		} else {
			text = "false"
		}

		qb.set("ignoreProxy", text)
	}
}

func (qb *queryBuilder) timeRange(timeRange *TimeRange) {
	if timeRange != nil {
		qb.set("starttimestamp", strconv.FormatInt(timeRange.Start.Unix(), 10))
		qb.set("endtimestamp", strconv.FormatInt(timeRange.End.Unix(), 10))
	}
}

func (qb *queryBuilder) blockRangeAdv(block BlockRangeAdv) {
	qb.set("fromBlock", block.FromBlock.String())

	if block.ToLatest {
		qb.set("toBlock", "latest")
	} else {
		qb.set("toBlock", block.ToBlock.String())
	}
}

func (qb *queryBuilder) topics(topics Topics) {
	// optional topic 1,2 and 3
	qb.set("topic0", add0x(topics.Topic0))

	// (1,0,0)
	if topics.Topic1 != nil {
		qb.set("topic1", add0x(*topics.Topic1))
		qb.set("topic0_1_opr", add0x(string(*topics.Opr01)))

		// (1,1,0)
		if topics.Topic2 != nil {
			qb.set("topic2", add0x(*topics.Topic2))
			qb.set("topic0_2_opr", add0x(string(*topics.Opr02)))
			qb.set("topic1_2_opr", add0x(string(*topics.Opr12)))
			
			// (1,1,1)
			if topics.Topic3 != nil {
				qb.set("topic3", add0x(*topics.Topic3))
				qb.set("topic0_3_opr", add0x(string(*topics.Opr03)))
				qb.set("topic1_3_opr", add0x(string(*topics.Opr13)))
				qb.set("topic2_3_opr", add0x(string(*topics.Opr23)))
			}
		
		// (1,0,0)
		} else {
			// (1,0,1)
			if topics.Topic3 != nil {
				qb.set("topic3", add0x(*topics.Topic3))
				qb.set("topic0_3_opr", add0x(string(*topics.Opr03)))
				qb.set("topic1_3_opr", add0x(string(*topics.Opr13)))
			} 
		}
	// (0,0,0)
	} else {
		// (0,1,0)
		if topics.Topic2 != nil {
			qb.set("topic2", add0x(*topics.Topic2))
			qb.set("topic0_2_opr", add0x(string(*topics.Opr02)))

			// (0,1,1)
			if topics.Topic3 != nil {
				qb.set("topic3", add0x(*topics.Topic3))
				qb.set("topic0_3_opr", add0x(string(*topics.Opr03)))
				qb.set("topic2_3_opr", add0x(string(*topics.Opr23)))
			}	

		// (0,0,0)
		} else {
			// (0,0,1)
			if topics.Topic3 != nil {
				qb.set("topic3", add0x(*topics.Topic3))
				qb.set("topic0_3_opr", add0x(string(*topics.Opr03)))
			} 
		}
	}
}

func (qb *queryBuilder) verify(contract ContractInfo) {
	qb.set("addressHash", add0x(contract.AddressHash))
	qb.set("name", contract.Name)
	qb.set("compilerVersion", contract.CompilerVersion)

	if contract.Optimization {
		qb.set("optimization", "true")
	} else {
		qb.set("optimization", "false")
	}
	
	qb.set("contractSourceCode", contract.ContractSourceCode)

	qb.setIfExist("constructorArguments", contract.ConstructorArguments)

	if contract.AutodetectConstructorArguments != nil {
		if *contract.AutodetectConstructorArguments {
			qb.set("autoDetectConstructorArguments", "true")
		} else {
			qb.set("autoDetectConstructorArguments", "false")
		}
	}

	qb.setIfExist("evmVersion", contract.EvmVersion)

	if contract.OptimizationRuns != nil {
		qb.set("optimizationRuns", strconv.Itoa(*contract.OptimizationRuns))
	}

	qb.setIfExist("proxyAddress", contract.ProxyAddress)
	qb.setIfExist("library1Name", contract.Library1Name)
	qb.setIfExist("library2Name", contract.Library2Name)
	qb.setIfExist("library3Name", contract.Library3Name)
	qb.setIfExist("library4Name", contract.Library4Name)
	qb.setIfExist("library5Name", contract.Library5Name)

	if contract.Library1Address != nil {
		qb.set("library1Address", add0x(*contract.Library1Address))
	}

	if contract.Library2Address != nil {
		qb.set("library2Address", add0x(*contract.Library2Address))
	}

	if contract.Library3Address != nil {
		qb.set("library3Address", add0x(*contract.Library3Address))
	}

	if contract.Library4Address != nil {
		qb.set("library4Address", add0x(*contract.Library4Address))
	}

	if contract.Library5Address != nil {
		qb.set("library5Address", add0x(*contract.Library5Address))
	}
}

// A nonnegative integer that represents the log index to be used for pagination.
func (qb *queryBuilder) index(num *int) {
	if num != nil {
		qb.set("index", strconv.Itoa(*num))
	} 
}

type ContractInfo struct {
	// I don't know how to align
	AddressHash string
	Name string
	CompilerVersion string
	Optimization bool
	ContractSourceCode string
	ConstructorArguments *string
	AutodetectConstructorArguments *bool
	EvmVersion *string
	OptimizationRuns *int
	ProxyAddress *string
	Library1Name *string
	Library1Address *string
	Library2Name *string
	Library2Address *string
	Library3Name *string
	Library3Address *string
	Library4Name *string
	Library4Address *string
	Library5Name *string
	Library5Address *string
}

type sortDirection string

const (
	Asc  sortDirection = "asc"
	Desc sortDirection = "desc"
)

type topicOperator string

const (
	And topicOperator = "and"
	Or topicOperator  = "or"
)

type Topics struct {
	Topic0 string
	Topic1 *string
	Topic2 *string
	Topic3 *string
	Opr01  *topicOperator
	Opr02  *topicOperator
	Opr03  *topicOperator
	Opr12  *topicOperator
	Opr13  *topicOperator
	Opr23  *topicOperator
}

type BlockRangeAdv struct {
	FromBlock *big.Int
	ToBlock   *big.Int
	ToLatest  bool
}

type BlockRange struct {
	StartBlock *big.Int
	EndBlock   *big.Int
}

type PageRange struct {
	Page   int
	Offset int
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

type filterDirection string

const (
	To   filterDirection = "to"
	From filterDirection = "from"
)

type filterContract string

const (
	Verified      filterContract = "verified"
	Decompiled 	  filterContract = "decompiled"
	Unverified 	  filterContract = "unverified"
	NotDecompiled filterContract = "not_decompiled"
	Empty 		  filterContract = "empty"
)

// Mimics Ethereum JSON RPC's eth_getBalance.
// Returns the wei balance (1 Celo = 10^18 wei) for an address as of the provided block (defaults to latest).
func (r *RequestClient) EthGetBalance(address string, block *big.Int) (string, error) {
	u := buildUrl(r.base, ethGetBalanceUrl)
	qb := newQueryBuilder(u)
	qb.address(address)
	qb.block(block)

	var ethResult EthResult
	var ethError EthError
	ok, err := r.jsonResponseDiff(u, ethResult, ethError)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", errors.New(ethError.Error)
	}

	return ethResult.Result, nil
}

// Get balance for address.
func (r *RequestClient) Balance(address string) (Balance, error) {
	u := buildUrl(r.base, balanceUrl)
	qb := newQueryBuilder(u)
	qb.address(address)

	var balance Balance
	err := r.jsonResponse(u, &balance)
	return balance, err
}

// Get balance for multiple addresses.
// If the balance hasn't been updated in a long time, we will double check with the node to fetch the absolute latest balance. This will not be reflected in the current request, but once it is updated, subsequent requests will show the updated balance. You can know that this is taking place via the `stale` attribute, which is set to `true` if a new balance is being fetched.
func (r *RequestClient) BalanceMulti(address []string) ([]BalanceMulti, error) {
	u := buildUrl(r.base, balanceMultiUrl)
	qb := newQueryBuilder(u)
	qb.addressMulti(address)

	var balanceMulti []BalanceMulti
	err := r.jsonResponse(u, &balanceMulti)
	return balanceMulti, err
}

// Get pending transactions by address.
func (r *RequestClient) PendingTxList(address string, page *PageRange) ([]PendingTxList, error) {
	u := buildUrl(r.base, pendingTxListUrl)
	qb := newQueryBuilder(u)
	qb.address(address)
	qb.pageRange(page)

	var pendingtxlist []PendingTxList
	err := r.jsonResponse(u, &pendingtxlist)
	return pendingtxlist, err
}

// Get transactions sent by an address. Up to a maximum of 10,000 transactions.
func (r *RequestClient) TxList(address string, sort *sortDirection, block *BlockRange, page *PageRange, filter *filterDirection, timeRange *TimeRange) ([]TxList, error) {
	u := buildUrl(r.base, txListUrl)
	qb := newQueryBuilder(u)
	qb.address(address)
	qb.sort(sort)
	qb.blockRange(block)
	qb.pageRange(page)
	qb.filterByDirection(filter)
	qb.timeRange(timeRange)

	var txList []TxList
	err := r.jsonResponse(u, &txList)
	return txList, err
}

// Get internal transactions by transaction or address hash. Up to a maximum of 10,000 internal transactions.
func (r *RequestClient) TxListInternal(txhash string, address *string, sort *sortDirection, block *BlockRange, page *PageRange) ([]TxListInternal, error) {
	u := buildUrl(r.base, txListInternalUrl)
	qb := newQueryBuilder(u)
	qb.txHash(txhash)
	if address != nil {
		qb.address(*address)
	}
	qb.sort(sort)
	qb.blockRange(block)
	qb.pageRange(page)
	
	var txListInternal []TxListInternal
	err := r.jsonResponse(u, &txListInternal)
	return txListInternal, err
}

// Get token transfer events by address. Up to a maximum of 10,000 token transfer events.
func (r *RequestClient) TokenTx(address string, contractAddress *string, sort *sortDirection, block *BlockRange, page *PageRange) ([]TokenTx, error) {
	u := buildUrl(r.base, tokenTxUrl)
	qb := newQueryBuilder(u)
	qb.address(address)
	if contractAddress != nil {
		qb.contractAddress(*contractAddress)
	}
	qb.sort(sort)
	qb.blockRange(block)
	qb.pageRange(page)

	var tokenTx []TokenTx
	err := r.jsonResponse(u, &tokenTx)
	return tokenTx, err
}

// Get token account balance for token contract address.
func (r *RequestClient) TokenBalance(contractAddress, address string) (TokenBalance, error) {
	u := buildUrl(r.base, tokenBalanceUrl)
	qb := newQueryBuilder(u)
	qb.contractAddress(contractAddress)
	qb.address(address)

	var tokenBalance TokenBalance
	err := r.jsonResponse(u, &tokenBalance)
	return tokenBalance, err
}

// Get list of tokens owned by address.
func (r *RequestClient) TokenList(address string) ([]TokenList, error) {
	u := buildUrl(r.base, tokenListUrl)
	qb := newQueryBuilder(u)
	qb.address(address)

	var tokenList []TokenList
	err := r.jsonResponse(u, &tokenList)
	return tokenList, err
}

// Get list of blocks mined by address.
func (r *RequestClient) GetMinedBlocks(address string, page *PageRange) ([]GetMinedBlocks, error) {
	u := buildUrl(r.base, getMinedBlocksUrl)
	qb := newQueryBuilder(u)
	qb.address(address)
	qb.pageRange(page)

	var getMinedBlocks []GetMinedBlocks
	err := r.jsonResponse(u, &getMinedBlocks)
	return getMinedBlocks, err
}

// Get a list of accounts and their balances, sorted ascending by the time they were first seen by the explorer.
func (r *RequestClient) ListAccounts(page *PageRange) ([]ListAccounts, error) {
	u := buildUrl(r.base, listAccountsUrl)
	qb := newQueryBuilder(u)
	qb.pageRange(page)

	var listAccounts []ListAccounts
	err := r.jsonResponse(u, &listAccounts)
	return listAccounts, err
}

// Get event logs for an address and/or topics. Up to a maximum of 1,000 event logs.
func (r *RequestClient) GetLogs(block BlockRangeAdv, contractAddress string, topics Topics) ([]GetLogs, error) {
	u := buildUrl(r.base, getLogsUrl)
	qb := newQueryBuilder(u)
	qb.blockRangeAdv(block)
	qb.address(contractAddress)
	qb.topics(topics)

	var getLogs []GetLogs
	err := r.jsonResponse(u, &getLogs)
	return getLogs, err
}

// Get ERC-20 or ERC-721 token by contract address.
func (r *RequestClient) GetToken(contractAddress string) (GetToken, error) {
	u := buildUrl(r.base, getTokenUrl)
	qb := newQueryBuilder(u)
	qb.contractAddress(contractAddress)

	var getToken GetToken
	err := r.jsonResponse(u, &getToken)
	return getToken, err
}

// Get token holders by contract address.
func (r *RequestClient) GetTokenHolders(contractAddress string, page *PageRange) ([]GetTokenHolders, error) {
	u := buildUrl(r.base, getTokenHoldersUrl)
	qb := newQueryBuilder(u)
	qb.contractAddress(contractAddress)
	qb.pageRange(page)

	var getTokenHolders []GetTokenHolders
	err := r.jsonResponse(u, &getTokenHolders)
	return getTokenHolders, err
}

// Get ERC-20 or ERC-721 token total supply by contract address.
func (r *RequestClient) TokenSupply(contractAddress string) (TokenSupply, error) {
	u := buildUrl(r.base, tokenSupplyUrl)
	qb := newQueryBuilder(u)
	qb.contractAddress(contractAddress)

	var tokenSupply TokenSupply
	err := r.jsonResponse(u, &tokenSupply)
	return tokenSupply, err
}

// Get total supply in Wei from exchange.
func (r *RequestClient) EthSupplyExchange() (EthSupplyExchange, error) {
	u := buildUrl(r.base, ethSupplyExchangeUrl)

	var ethSupplyExchange EthSupplyExchange
	err := r.jsonResponse(u, &ethSupplyExchange)
	return ethSupplyExchange, err
}

// Get total supply in Wei from DB.
func (r *RequestClient) EthSupply() (EthSupply, error) {
	u := buildUrl(r.base, ethSupplyUrl)

	var ethSupply EthSupply
	err := r.jsonResponse(u, &ethSupply)
	return ethSupply, err
}

// Get total coin supply from DB minus burnt number.
func (r *RequestClient) CoinSupply() (CoinSupply, error) {
	u := buildUrl(r.base, coinSupplyUrl)

	var coinSupply CoinSupply
	err := r.jsonResponse(u, &coinSupply)
	return coinSupply, err
}

// Get latest price in USD and BTC.
func (r *RequestClient) EthPrice() (EthPrice, error) {
	u := buildUrl(r.base, ethPriceUrl)

	var ethPrice EthPrice
	err := r.jsonResponse(u, &ethPrice)
	return ethPrice, err
}

// Get estimated total number of transactions.
func (r *RequestClient) TotalTransactions() (TotalTransactions, error) {
	u := buildUrl(r.base, totalTransactionsUrl)

	var totalTransactions TotalTransactions
	err := r.jsonResponse(u, &totalTransactions)
	return totalTransactions, err
}

// Get block reward by block number.
func (r *RequestClient) GetBlockReward(blockNumber *big.Int) (GetBlockReward, error) {
	u := buildUrl(r.base, getBlockRewardUrl)
	qb := newQueryBuilder(u)
	qb.blockNo(blockNumber)

	var getBlockReward GetBlockReward
	err := r.jsonResponse(u, &getBlockReward)
	return getBlockReward, err
}

// Mimics Ethereum JSON RPC's eth_blockNumber. Returns the lastest block number
func (r *RequestClient) EthBlockNumber() (string, error) {
	u := buildUrl(r.base, ethBlockNumberUrl)

	var ethResult EthResult
	var ethError EthError
	ok, err := r.jsonResponseDiff(u, ethResult, ethError)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", errors.New(ethError.Error)
	}

	return ethResult.Result, nil
}

// Get a list of contracts, sorted ascending by the time they were first seen by the explorer. If you provide the filters `not_decompiled`(`4`) or `not_verified(4)` the results will not be sorted for performance reasons.
func (r *RequestClient) ListContracts(page *PageRange, filter *filterContract, notVersion *string) ([]ListContracts, error) {
	u := buildUrl(r.base, listContractsUrl)
	qb := newQueryBuilder(u)
	qb.pageRange(page)
	qb.filterContract(filter)
	qb.notDecompiledWithVersion(notVersion)

	var listContracts []ListContracts
	err := r.jsonResponse(u, &listContracts)
	return listContracts, err
}

// Get ABI for verified contract. 
func (r *RequestClient) GetAbi(address string) (GetAbi, error) {
	u := buildUrl(r.base, getAbiUrl)
	qb := newQueryBuilder(u)
	qb.address(address)

	var getAbi GetAbi
	err := r.jsonResponse(u, &getAbi)
	return getAbi, err
}

// Get contract source code for verified contract.
func (r *RequestClient) GetSourceCode(address string, ignoreProxy *bool) (GetSourceCode, error) {
	u := buildUrl(r.base, getSourceCodeUrl)
	qb := newQueryBuilder(u)
	qb.address(address)
	qb.ignoreProxy(ignoreProxy)

	var getSourceCode GetSourceCode
	err := r.jsonResponse(u, &getSourceCode)
	return getSourceCode, err
}

// Verify a contract with its source code and contract creation information.
func (r *RequestClient) Verify(contract ContractInfo) (Verify, error) {
	u := buildUrl(r.base, verifyUrl)
	qb := newQueryBuilder(u)
	qb.verify(contract)

	var verify Verify
	err := r.jsonResponse(u, &verify)
	return verify, err
}

// Get transaction info.
func (r *RequestClient) GetTxInfo(txhash string, index *int) (GetTxInfo, error) {
	u := buildUrl(r.base, getTxInfoUrl)
	qb := newQueryBuilder(u)
	qb.txHash(txhash)
	qb.index(index)

	var getTxInfo GetTxInfo
	err := r.jsonResponse(u, &getTxInfo)
	return getTxInfo, err
}

// Get transaction receipt status.
func (r *RequestClient) GetTxReceiptStatus(txhash string) (GetTxReceiptStatus, error) {
	u := buildUrl(r.base, getTxReceiptStatusUrl)
	qb := newQueryBuilder(u)
	qb.txHash(txhash)

	var getTxReceiptStatus GetTxReceiptStatus
	err := r.jsonResponse(u, &getTxReceiptStatus)
	return getTxReceiptStatus, err
}

// Get error status and error message.
func (r *RequestClient) GetStatus(txhash string) (GetStatus, error) {
	u := buildUrl(r.base, getStatusUrl)
	qb := newQueryBuilder(u)
	qb.txHash(txhash)

	var getStatus GetStatus
	err := r.jsonResponse(u, &getStatus)
	return getStatus, err
}