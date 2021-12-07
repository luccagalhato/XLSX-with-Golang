package apihandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"social-planilha/sql"
	"strconv"
)

//XML ...
func XML(s *sql.Str) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		b, err := ioutil.ReadAll(file)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data := xmlUnMarshal(b)
		gtin := make([]string, 0)
		for _, product := range data.NFe.InfNFe.Det {

			if _, err := strconv.Atoi(product.Prod.CEAN); err == nil {
				gtin = append(gtin, product.Prod.CEAN)
				continue
			}

			gtin = append(gtin, product.Prod.Codigo)
		}
		excels := gerarExcel(s.GetGtin(gtin))

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("access-control-expose-headers", "*")
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(excels)
	}
}
