package models

import "io"

//Prod ...
type Prod struct {
	Codigo string `xml:"cProd"`
	CEAN   string `xml:"cEAN"`
}

//Det ...
type Det struct {
	Prod Prod `xml:"prod"`
}

//InfNFe ...
type InfNFe struct {
	Det []Det `xml:"det"`
}

//NFe ...
type NFe struct {
	InfNFe InfNFe `xml:"infNFe"`
}

//DataFormat ...
type DataFormat struct {
	NFe NFe `xml:"NFe"`
}

//Date ...
type Date struct {
	DataInicial string `json:"dataInicial"`
	DataFinal   string `json:"dataFinal"`
}

//Excel ...
type Excel struct {
	Name  string
	Value io.Reader
}
