package sql

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/denisenkom/go-mssqldb" //bblablalba
)

type Nencontrados struct {
	Codgtin string
}

const (
	selectFormat string = `SELECT LTRIM(RTRIM(COALESCE(P.PRODUTO, ''))) AS PRODUTO, LTRIM(RTRIM(COALESCE(P.COR_PRODUTO,''))) AS COR_PRODUTO,LTRIM(RTRIM(COALESCE(A.GRIFFE, ''))) AS GRIFFE,LTRIM(RTRIM(COALESCE(A.DESC_PRODUTO, ''))) AS DESC_PRODUTO,LTRIM(RTRIM(COALESCE(B.DESC_COR,''))) AS DESC_COR,LTRIM(RTRIM(COALESCE(A.PESO,'0.00'))) AS PESO,LTRIM(RTRIM(COALESCE(A.CLASSIF_FISCAL, ''))) AS CLASSIF_FISCAL,LTRIM(RTRIM(COALESCE(P.GRADE, ''))) AS GRADE,LTRIM(RTRIM(COALESCE(E.PRECO1,'0.00'))) AS PRECO_VENDA,LTRIM(RTRIM(COALESCE(D.PRECO1,'0.00'))) AS PRECO_FATURAMENTO,LTRIM(RTRIM(COALESCE(P.[0],''))) AS COD_MILLENNIUM, LTRIM(RTRIM(COALESCE(P.[1],''))) AS COD_GTIN, LTRIM(RTRIM(COALESCE(P.[3],''))) AS COD_LINX,LTRIM(RTRIM(COALESCE(C.DESC_COMPOSICAO, ''))) AS DESC_COMPOSICAO, LTRIM(RTRIM(COALESCE(F.DESC_SEXO_TIPO,''))) AS GENERO, LTRIM(RTRIM(G.DESC_COLECAO)) AS COLECAO, LTRIM(RTRIM(A.GRUPO_PRODUTO)) AS GRUPO_PRODUTO
	FROM ( SELECT ROW_NUMBER() OVER(PARTITION BY PRODUTO,COR_PRODUTO,GRADE,TIPO_COD_BAR ORDER BY CODIGO_BARRA DESC, TIPO_COD_BAR) AS ID,
        CODIGO_BARRA,PRODUTO,COR_PRODUTO,GRADE, TIPO_COD_BAR,TAMANHO
        FROM LINX_TBFG..PRODUTOS_BARRA
        WHERE INATIVO='0'
        ) A
        PIVOT ( MAX(A.CODIGO_BARRA)
        FOR A.TIPO_COD_BAR IN ([0], [1], [3])) P
	LEFT JOIN LINX_TBFG..PRODUTOS A ON P.PRODUTO=A.PRODUTO
	LEFT JOIN LINX_TBFG..CORES_BASICAS B ON P.COR_PRODUTO=B.COR
    LEFT JOIN LINX_TBFG..MATERIAIS_COMPOSICAO C ON A.COMPOSICAO=C.COMPOSICAO
	LEFT JOIN LINX_TBFG..PRODUTOS_PRECOS D ON P.PRODUTO=D.PRODUTO AND D.CODIGO_TAB_PRECO='51'
	LEFT JOIN LINX_TBFG..PRODUTOS_PRECOS E ON P.PRODUTO=E.PRODUTO AND E.CODIGO_TAB_PRECO='01'
	LEFT JOIN LINX_TBFG..W_SEXO_TIPO F ON A.SEXO_TIPO=F.SEXO_TIPO
	LEFT JOIN LINX_TBFG..COLECOES G ON A.COLECAO=G.COLECAO
	WHERE A.STATUS_PRODUTO='04' AND P.ID='1'`
)

type Str struct {
	url *url.URL
	db  *sql.DB
}

//Product ...
type Product struct {
	Produto          string  `json:"PRODUTO,omitempty"`
	CorProduto       string  `json:"COR_PRODUTO,omitempty"`
	Marca            string  `json:"GRIFFE,omitempty"`
	DescProduto      string  `json:"DESC_PRODUTO,omitempty"`
	DescColorProd    string  `json:"DESC_COR,omitempty"`
	Peso             float64 `json:"PESO,omitempty"`
	Ncm              string  `json:"CLASSIF_FISCAL,omitempty"`
	Tamanho          string  `json:"GRADE,omitempty"`
	PrecoFaturamento float64 `json:"PRECO_FATURAMENTO,omitempty"`
	PrecoVenda       float64 `json:"PRECO_VENDA,omitempty"`
	CodLinx          string  `json:"COD_LINX,omitempty"`
	CodMillenium     string  `json:"COD_MILLENIUM,omitempty"`
	CodGtin          string  `json:"COD_GTIN,omitempty"`
	Composicao       string  `json:"DESC_COMPOSICAO,omitempty"`
	Genero           string  `json:"GENERO,omitempty"`
	DescColecao      string  `json:"GRUPO_PRODUTO,omitempty"`
	DescCategoria    string  `json:"COLECAO,omitempty"`
}

//GetGtin ...
func (s *Str) GetGtin(gtins []string) []*Product {
	if len(gtins) == 0 {
		return nil
	}
	rst := make([]*Product, 0)
	gts := strings.Join(gtins, "', '")
	sel := fmt.Sprintf(selectFormat+` AND A.PRODUTO IN (SELECT DISTINCT PRODUTO FROM LINX_TBFG..PRODUTOS_BARRA
		WHERE CODIGO_BARRA IN ('%s') OR PRODUTO IN ('%s'))
	ORDER BY P.PRODUTO, P.COR_PRODUTO, P.TAMANHO`, gts, gts)
	rows, err := s.db.QueryContext(context.Background(), sel, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for rows.Next() {
		Product := Product{}

		if err := rows.Scan(&Product.Produto, &Product.CorProduto, &Product.Marca, &Product.DescProduto, &Product.DescColorProd, &Product.Peso, &Product.Ncm, &Product.Tamanho, &Product.PrecoVenda, &Product.PrecoFaturamento, &Product.CodMillenium, &Product.CodGtin, &Product.CodLinx, &Product.Composicao, &Product.Genero, &Product.DescColecao, &Product.DescCategoria); err != nil {
			fmt.Println(err)
			continue
		}
		rst = append(rst, &Product)

	}

	for _, gtin := range gtins {
		encontrado := false
		for _, product := range rst {
			if gtin == product.CodGtin || gtin == product.CodLinx || gtin == product.CodMillenium || gtin == product.Produto {
				encontrado = true
			}
		}
		if !encontrado {
			rst = append(rst, &Product{Produto: gtin, Marca: "NAO ENCONTRADO"})
		}
	}
	return rst
}

//GetDate
func (s *Str) GetDate(dataInicial string, dataFinal string) []*Product {

	fmt.Println("Buscando itens......")
	rst := make([]*Product, 0)
	sel := fmt.Sprintf(selectFormat+` AND A.DATA_PARA_TRANSFERENCIA BETWEEN '%s' AND '%s'
	ORDER BY P.PRODUTO, P.COR_PRODUTO, P.TAMANHO`, dataInicial, dataFinal)
	rows, err := s.db.QueryContext(context.Background(), sel, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for rows.Next() {

		Product := Product{}

		if err := rows.Scan(&Product.Produto, &Product.CorProduto, &Product.Marca, &Product.DescProduto, &Product.DescColorProd, &Product.Peso, &Product.Ncm, &Product.Tamanho, &Product.PrecoVenda, &Product.PrecoFaturamento, &Product.CodMillenium, &Product.CodGtin, &Product.CodLinx, &Product.Composicao, &Product.Genero, &Product.DescColecao, &Product.DescCategoria); err != nil {
			fmt.Println(err)
			continue
		}
		rst = append(rst, &Product)

	}
	fmt.Println("Busca de itens conluída")
	return rst

}

//MakeSQL
func MakeSQL(host, port, username, password string) (*Str, error) {
	s := &Str{}
	s.url = &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(username, password),
		Host:     fmt.Sprintf("%s:%s", host, port),
		RawQuery: url.Values{}.Encode(),
	}
	return s, s.connect()
}

//Ping ...
func (s *Str) Ping() error {
	return s.db.Ping()
}
func (s *Str) connect() error {
	var err error
	if s.db, err = sql.Open("sqlserver", s.url.String()); err != nil {
		return err
	}
	return s.db.PingContext(context.Background())
}

// func (s *Str) disconnect() error {
// 	return s.db.Close()
// }
