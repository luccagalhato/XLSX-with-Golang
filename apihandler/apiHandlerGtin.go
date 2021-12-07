package apihandler

import (
	"encoding/json"
	"net/http"
	"social-planilha/sql"
)

//Gtin ...
func Gtin(s *sql.Str) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// rst := connection.GetGtin(s)
		excelIds := gerarExcel(s.GetGtin(genStrs(r.Body)))

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("access-control-expose-headers", "*")
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(excelIds)
	}
}
