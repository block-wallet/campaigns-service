package model

type Event struct {
	BlockNumber     uint64
	Commitment      string
	LeafIndex       uint32
	Timestamp       string
	TransactionHash string
}

type Chain struct {
	Name           string
	Chain          string
	Network        string
	Icon           string
	Rpc            []string
	Faucet         []string
	NativeCurrency *Currency
	InfoURL        string
	ShortName      string
	ChainId        uint64
	NetworkId      uint64
	Ens            *Ens
	Explorers      *[]Explorer
}

type Currency struct {
	Name     string
	Symbol   string
	Decimals uint64
}

type Ens struct {
	Registry string
}

type Explorer struct {
	Name     string
	Url      string
	Standard string
}
