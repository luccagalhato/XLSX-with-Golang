package apihandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	modelss "social-planilha/models"
	"social-planilha/sql"
)

//Date ...
func Date(s *sql.Str) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		date := new(modelss.Date)
		if err := json.NewDecoder(r.Body).Decode(date); err != nil {
			fmt.Println(err)
		}
		excels := gerarExcel(s.GetDate(date.DataInicial, date.DataFinal))
		fmt.Println("Arquivos", excels)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("access-control-expose-headers", "*")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(excels)
	}
}
