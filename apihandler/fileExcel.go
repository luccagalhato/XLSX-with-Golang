package apihandler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"social-planilha/config"
	"social-planilha/models"
	"social-planilha/sql"
	"strconv"
	"strings"
	"time"

	_ "embed" //EMB

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/rivo/uniseg"
)

var (
	//go:embed modelosmk.xlsx
	smk []byte
	//go:embed modelonaoencontrado.xlsx
	naoencontrado []byte
	//go:embed modeloyanni.xlsx
	yanni []byte
	//go:embed modeloplanet.xlsx
	planet []byte
	//go:embed modelopolo.xlsx
	polo []byte

	modelsGriffe map[string][]byte = map[string][]byte{
		"POLO WEAR":      polo,
		"PLANET GIRLS":   planet,
		"SMK":            smk,
		"NAO ENCONTRADO": naoencontrado,
	}
	deparaCatSMK  map[string]string
	deparaDescSMK []class
	deparaCatPG   map[string]string
	deparaCatPW   map[string]string
	deparaColor   map[string]string
	deparaDescPW  []class
	deparaDescPG  []class
)

//LoadConfig ...
func LoadConfig(d config.DePara) {

	pgFile, err := os.Open(d.CatPG)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(pgFile).Decode(&deparaCatPG); err != nil {
		log.Fatal(err)
	}

	pwFile, err := os.Open(d.CatPW)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(pwFile).Decode(&deparaCatPW); err != nil {
		log.Fatal(err)
	}

	colorFile, err := os.Open(d.Color)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(colorFile).Decode(&deparaColor); err != nil {
		log.Fatal(err)
	}

	pgFile, err = os.Open(d.DescPG)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(pgFile).Decode(&deparaDescPG); err != nil {
		log.Fatal(err)
	}

	pwFile, err = os.Open(d.DescPW)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(pwFile).Decode(&deparaDescPW); err != nil {
		log.Fatal(err)
	}
	smkFile, err := os.Open(d.CatSMK)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(smkFile).Decode(&deparaCatSMK); err != nil {
		log.Fatal(err)
	}

	smkFile, err = os.Open(d.DescSMK)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(smkFile).Decode(&deparaDescSMK); err != nil {
		log.Fatal(err)
	}

}

func gerarExcel(rst []*sql.Product) []string {
	fmt.Println("Gerando Marca....")
	pormarca := make(map[string][]*sql.Product)
	for _, produto := range rst {
		if _, ok := pormarca[produto.Marca]; !ok {
			pormarca[produto.Marca] = make([]*sql.Product, 0)
		}
		pormarca[produto.Marca] = append(pormarca[produto.Marca], produto)
	}

	keys := make([]string, 0)
	for marca := range pormarca {

		var f *excelize.File
		switch marca {
		case "YANNI":
			fmt.Println("Gerando excel para:", marca)
			f = excelYanni(pormarca[marca])
		case "POLO WEAR":
			fmt.Println("Gerando excel para:", marca)
			f = excelPW(marca, "PW-", pormarca[marca])
		case "PLANET GIRLS":
			fmt.Println("Gerando excel para:", marca)
			f = excelPG(marca, "", pormarca[marca])
		case "SMK":
			fmt.Println("Gerando excel para:", marca)
			f = excelSMK(marca, "", pormarca[marca])
		case "NAO ENCONTRADO":
			f = excelNencontrado(marca, "", pormarca[marca])
		default:
			continue
		}
		if b, err := f.WriteToBuffer(); err == nil {
			keys = append(keys, addExcel(fmt.Sprintf("%s_%s.xlsx", marca, time.Now().Format("20060102_150405")), b))
		}

	}

	return keys
	// pretty.Print(rst)
}

func excelNencontrado(marca, prefix string, produtos []*sql.Product) *excelize.File {
	f, err := excelize.OpenReader(bytes.NewReader(modelsGriffe[marca]))
	if err != nil {
		log.Println(err)
	}
	sheet := f.GetSheetName(1)
	for i, produto := range produtos {
		f.SetCellStr(sheet, fmt.Sprintf("A%d", i+2), produto.Produto)
	}
	return f
}

func excelYanni(produtos []*sql.Product) *excelize.File {
	f, err := excelize.OpenReader(bytes.NewReader(yanni))
	if err != nil {
		log.Println(err)
	}
	sheet := f.GetSheetName(1)
	for i, produto := range produtos {
		f.SetCellStr(sheet, fmt.Sprintf("A%d", i+2), Title(produto.Produto+"-"+produto.CorProduto))
		f.SetCellStr(sheet, fmt.Sprintf("B%d", i+2), Title(produto.DescProduto+" "+produto.DescColorProd))
		f.SetCellStr(sheet, fmt.Sprintf("D%d", i+2), Title(produto.DescColecao))
		f.SetCellStr(sheet, fmt.Sprintf("E%d", i+2), Title("YANNÌ"))
		f.SetCellStr(sheet, fmt.Sprintf("F%d", i+2), Title(fmt.Sprintf("%.2f", produto.Peso)))
		f.SetCellStr(sheet, fmt.Sprintf("G%d", i+2), Title(produto.Ncm))
		f.SetCellStr(sheet, fmt.Sprintf("H%d", i+2), produto.Tamanho)
		f.SetCellStr(sheet, fmt.Sprintf("I%d", i+2), Title(produto.DescColorProd))
		f.SetCellStr(sheet, fmt.Sprintf("J%d", i+2), Title("15"))
		f.SetCellStr(sheet, fmt.Sprintf("K%d", i+2), Title("15"))
		f.SetCellStr(sheet, fmt.Sprintf("L%d", i+2), Title("10"))
		f.SetCellStr(sheet, fmt.Sprintf("M%d", i+2), Title(fmt.Sprintf("%.2f", produto.PrecoFaturamento)))
		f.SetCellStr(sheet, fmt.Sprintf("N%d", i+2), Title(fmt.Sprintf("%.2f", produto.PrecoVenda)))
		cod := produto.CodMillenium
		if cod == "" {
			cod = produto.CodLinx
		}
		f.SetCellStr(sheet, fmt.Sprintf("O%d", i+2), Title(cod))
		f.SetCellStr(sheet, fmt.Sprintf("P%d", i+2), Title(produto.CodGtin))
		f.SetCellStr(sheet, fmt.Sprintf("Q%d", i+2), produto.Tamanho+":")
		f.SetCellStr(sheet, fmt.Sprintf("R%d", i+2), Title(produto.Composicao))
		f.SetCellStr(sheet, fmt.Sprintf("S%d", i+2), Title(produto.DescProduto))
	}
	return f
}

func excelPG(marca, prefix string, produtos []*sql.Product) *excelize.File {
	f, err := excelize.OpenReader(bytes.NewReader(modelsGriffe[marca]))
	if err != nil {
		log.Println(err)
	}
	sheet := f.GetSheetName(1)
	for i, produto := range produtos {
		f.SetCellStr(sheet, fmt.Sprintf("A%d", i+2), Title(prefix+produto.Produto+"-"+produto.CorProduto))
		f.SetCellStr(sheet, fmt.Sprintf("B%d", i+2), Title(produto.DescProduto+" "+produto.Marca+" "+deparaColor[produto.DescColorProd]))
		f.SetCellStr(sheet, fmt.Sprintf("C%d", i+2), Title("ROUPAS"))
		f.SetCellStr(sheet, fmt.Sprintf("D%d", i+2), Title(deparaCatPG[produto.DescCategoria]))

		for _, c := range deparaDescPG {

			if strings.Contains(mn(produto.DescProduto), c.Name) {
				f.SetCellStr(sheet, fmt.Sprintf("D%d", i+2), Title(c.Categoria))
				if c.Departamento != "" {
					f.SetCellStr(sheet, fmt.Sprintf("C%d", i+2), Title(c.Departamento))
				}
				break
			}
		}

		if isPlusSize(produto.Tamanho) {
			f.SetCellStr(sheet, fmt.Sprintf("C%d", i+2), Title("PLUS SIZE"))
		}

		f.SetCellStr(sheet, fmt.Sprintf("E%d", i+2), Title(produto.DescColecao))
		f.SetCellStr(sheet, fmt.Sprintf("F%d", i+2), Title(marca))
		f.SetCellStr(sheet, fmt.Sprintf("G%d", i+2), Title(produto.Genero))
		f.SetCellStr(sheet, fmt.Sprintf("H%d", i+2), Title(fmt.Sprintf("%.2f", produto.Peso)))
		f.SetCellStr(sheet, fmt.Sprintf("I%d", i+2), Title(produto.Ncm))
		f.SetCellStr(sheet, fmt.Sprintf("J%d", i+2), produto.Tamanho)
		f.SetCellStr(sheet, fmt.Sprintf("K%d", i+2), Title(deparaColor[produto.DescColorProd]))
		f.SetCellStr(sheet, fmt.Sprintf("L%d", i+2), Title("15"))
		f.SetCellStr(sheet, fmt.Sprintf("M%d", i+2), Title("15"))
		f.SetCellStr(sheet, fmt.Sprintf("N%d", i+2), Title("10"))
		f.SetCellStr(sheet, fmt.Sprintf("O%d", i+2), Title(fmt.Sprintf("%.2f", produto.PrecoFaturamento)))
		f.SetCellStr(sheet, fmt.Sprintf("P%d", i+2), Title(fmt.Sprintf("%.2f", produto.PrecoVenda)))
		cod := produto.CodMillenium
		if cod == "" {
			cod = produto.CodLinx
		}
		f.SetCellStr(sheet, fmt.Sprintf("Q%d", i+2), Title(cod))
		f.SetCellStr(sheet, fmt.Sprintf("R%d", i+2), Title(produto.CodGtin))
		f.SetCellStr(sheet, fmt.Sprintf("S%d", i+2), Title(produto.Tamanho+":"))
		f.SetCellStr(sheet, fmt.Sprintf("T%d", i+2), Title(produto.Composicao))
		f.SetCellStr(sheet, fmt.Sprintf("U%d", i+2), Title(produto.DescProduto))
	}
	return f
}
func excelPW(marca, prefix string, produtos []*sql.Product) *excelize.File {
	f, err := excelize.OpenReader(bytes.NewReader(modelsGriffe[marca]))
	if err != nil {
		log.Println(err)
	}
	sheet := f.GetSheetName(1)
	for i, produto := range produtos {
		f.SetCellStr(sheet, fmt.Sprintf("A%d", i+2), prefix+Title(produto.Produto+"-"+produto.CorProduto))
		f.SetCellStr(sheet, fmt.Sprintf("B%d", i+2), Title(produto.DescProduto+" "+produto.Marca+" "+deparaColor[produto.DescColorProd]))
		f.SetCellStr(sheet, fmt.Sprintf("E%d", i+2), Title(produto.DescColecao))
		f.SetCellStr(sheet, fmt.Sprintf("F%d", i+2), Title(marca))
		f.SetCellStr(sheet, fmt.Sprintf("G%d", i+2), Title(produto.Genero))
		f.SetCellStr(sheet, fmt.Sprintf("H%d", i+2), Title(fmt.Sprintf("%.2f", produto.Peso)))
		f.SetCellStr(sheet, fmt.Sprintf("I%d", i+2), Title(produto.Ncm))
		f.SetCellStr(sheet, fmt.Sprintf("C%d", i+2), Title(produto.Genero))
		f.SetCellStr(sheet, fmt.Sprintf("D%d", i+2), Title(deparaCatPW[produto.DescCategoria]))
		for _, c := range deparaDescPW {
			if strings.Contains(mn(produto.DescProduto), c.Name) {
				f.SetCellStr(sheet, fmt.Sprintf("D%d", i+2), Title(c.Categoria))
				if c.Departamento != "" {
					f.SetCellStr(sheet, fmt.Sprintf("C%d", i+2), Title(c.Departamento))
				}
				break
			}
		}

		if isPlusSize(produto.Tamanho) {
			f.SetCellStr(sheet, fmt.Sprintf("C%d", i+2), Title("PLUS SIZE"))
		}

		f.SetCellStr(sheet, fmt.Sprintf("J%d", i+2), produto.Tamanho)
		f.SetCellStr(sheet, fmt.Sprintf("K%d", i+2), Title(deparaColor[produto.DescColorProd]))
		f.SetCellStr(sheet, fmt.Sprintf("L%d", i+2), Title("15"))
		f.SetCellStr(sheet, fmt.Sprintf("M%d", i+2), Title("15"))
		f.SetCellStr(sheet, fmt.Sprintf("N%d", i+2), Title("10"))
		f.SetCellStr(sheet, fmt.Sprintf("O%d", i+2), Title(fmt.Sprintf("%.2f", produto.PrecoFaturamento)))
		f.SetCellStr(sheet, fmt.Sprintf("P%d", i+2), Title(fmt.Sprintf("%.2f", produto.PrecoVenda)))
		cod := produto.CodMillenium
		if cod == "" {
			cod = produto.CodLinx
		}
		f.SetCellStr(sheet, fmt.Sprintf("Q%d", i+2), Title(cod))
		f.SetCellStr(sheet, fmt.Sprintf("R%d", i+2), Title(produto.CodGtin))
		f.SetCellStr(sheet, fmt.Sprintf("S%d", i+2), Title(produto.Tamanho+":"))
		f.SetCellStr(sheet, fmt.Sprintf("T%d", i+2), Title(produto.Composicao))
		f.SetCellStr(sheet, fmt.Sprintf("U%d", i+2), Title(produto.DescProduto))
	}
	return f
}
func excelSMK(marca, prefix string, produtos []*sql.Product) *excelize.File {
	f, err := excelize.OpenReader(bytes.NewReader(modelsGriffe[marca]))
	if err != nil {
		log.Println(err)
	}
	sheet := f.GetSheetName(2)
	for i, produto := range produtos {
		str := produto.DescProduto
		str2 := strings.Split(str, " ")
		last := str2[len(str2)-1]
		_, err := strconv.Atoi(last)
		var descProduto []string
		if err != nil {
			descProduto = str2[1:]
		} else {
			sl := str2[:len(str2)-1]
			descProduto = sl[1:]
		}

		f.SetCellStr(sheet, fmt.Sprintf("A%d", i+3), prefix+Title(produto.Produto+"-"+produto.CorProduto))
		cod := produto.CodMillenium
		if cod == "" {
			cod = produto.CodLinx
		}

		f.SetCellStr(sheet, fmt.Sprintf("B%d", i+3), Title(produto.CodGtin))
		f.SetCellStr(sheet, fmt.Sprintf("C%d", i+3), Title(cod))
		f.SetCellStr(sheet, fmt.Sprintf("D%d", i+3), Title(strings.Join(descProduto, " ")+" "+produto.Marca+" "+produto.DescColorProd))
		f.SetCellStr(sheet, fmt.Sprintf("E%d", i+3), Title(strconv.Itoa(uniseg.GraphemeClusterCount(produto.DescProduto))))
		f.SetCellStr(sheet, fmt.Sprintf("F%d", i+3), Title(produto.Marca))
		f.SetCellStr(sheet, fmt.Sprintf("H%d", i+3), Title("NÃO"))
		for _, c := range deparaDescSMK {
			if strings.Contains(mn(produto.DescProduto), c.Name) {
				f.SetCellStr(sheet, fmt.Sprintf("k%d", i+3), Title(c.Categoria))
				if c.Departamento != "" {
					f.SetCellStr(sheet, fmt.Sprintf("J%d", i+3), Title(c.Departamento))
				}
				break
			} else {
				f.SetCellStr(sheet, fmt.Sprintf("J%d", i+3), Title("ROUPAS"))
				f.SetCellStr(sheet, fmt.Sprintf("K%d", i+3), Title(produto.DescCategoria))
			}
		}
		f.SetCellStr(sheet, fmt.Sprintf("L%d", i+3), Title(produto.DescColecao))
		f.SetCellStr(sheet, fmt.Sprintf("M%d", i+3), produto.Tamanho)
		f.SetCellStr(sheet, fmt.Sprintf("N%d", i+3), produto.DescColorProd)
		f.SetCellStr(sheet, fmt.Sprintf("P%d", i+3), Title(fmt.Sprintf("%.2f", produto.Peso)))
		f.SetCellStr(sheet, fmt.Sprintf("Q%d", i+3), Title("15"))
		f.SetCellStr(sheet, fmt.Sprintf("R%d", i+3), Title("15"))
		f.SetCellStr(sheet, fmt.Sprintf("S%d", i+3), Title("10"))
		f.SetCellStr(sheet, fmt.Sprintf("T%d", i+3), "DZ")
		f.SetCellStr(sheet, fmt.Sprintf("U%d", i+3), "UN")
		f.SetCellStr(sheet, fmt.Sprintf("V%d", i+3), Title(fmt.Sprintf("%.2f", produto.PrecoFaturamento)))
		f.SetCellStr(sheet, fmt.Sprintf("W%d", i+3), Title(fmt.Sprintf("%.2f", produto.PrecoVenda)))
		f.SetCellStr(sheet, fmt.Sprintf("Y%d", i+3), ("MATERIAL"))
		f.SetCellStr(sheet, fmt.Sprintf("Z%d", i+3), (produto.Composicao))

	}
	return f
}

func genStrs(r io.Reader) []string {
	s := bufio.NewScanner(r)
	s.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		start := 0
		for ; start < len(data); start++ {
			if data[start] > 47 && data[start] < 58 {
				break
			}
		}
		// Scan until space, marking end of word.
		for i := start + 1; i < len(data); i++ {

			if data[i] < 48 || data[i] > 57 {
				return i + 1, data[start:i], nil
			}
		}
		// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
		if atEOF && len(data) > start {
			return len(data), data[start:], nil
		}
		// Request more data.
		return start, nil, nil
	})
	rst := make([]string, 0)
	for s.Scan() {
		rst = append(rst, s.Text())
	}
	return rst
}

func xmlUnMarshal(b []byte) models.DataFormat {

	data := models.DataFormat{}
	err := xml.Unmarshal(b, &data)
	if nil != err {
		fmt.Println("Error unmarshalling from XML", err)
	}
	return data
}
