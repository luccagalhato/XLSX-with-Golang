package apihandler

import (
	"net/http"
)

//GtinandCor ...
func GtinandCor(w http.ResponseWriter, r *http.Request) {
	// rst := connection.GetGtin(s)
	//excelIds := gerarExcel(connection.GetGtinandCor(genStrsbyCor(r.Body)))

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("access-control-expose-headers", "*")
	w.Header().Set("Content-Type", "application/json")

	//json.NewEncoder(w).Encode(excelIds)
}
