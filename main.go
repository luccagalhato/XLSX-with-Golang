package main

import (
	"embed"
	"flag"
	"io"
	"io/fs"
	"log"
	"net/http"
	"social-planilha/apihandler"
	c "social-planilha/config"
	"social-planilha/sql"
	"time"

	_ "embed"

	"github.com/gorilla/mux"
)

//go:embed html
var content embed.FS
var createConfig bool

func main() {

	flag.BoolVar(&createConfig, "c", false, "create config.yaml file")
	flag.Parse()

	if createConfig {
		c.CreateConfigFile()
		return
	}

	log.Print("loading config file")
	if err := c.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	connection, err := sql.MakeSQL(c.Config.SQL.Host, c.Config.SQL.Port, c.Config.SQL.User, c.Config.SQL.Password)
	if err != nil {
		log.Println(err)
		return
	}
	apihandler.LoadConfig(c.Config.DePara)
	htmlFS, err := fs.Sub(content, "html")
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("starting server '%s' at port: %s", c.Config.API.Host, c.Config.API.Port)
	go func() {
		tickerPing := time.NewTicker(time.Second * 10)
		for range tickerPing.C {
			connection.Ping()
		}
	}()

	// fs := http.FileServer(http.Dir("html"))
	fs := http.FileServer(http.FS(htmlFS))

	r := mux.NewRouter()
	r.PathPrefix("/html/").Handler(http.StripPrefix("/html/", fs))
	r.Path("/id/{id}").Methods(http.MethodGet).HandlerFunc(apiDownloadHandler)
	r.Path("/xml").Methods(http.MethodPost).HandlerFunc(apihandler.XML(connection))
	r.Path("/gtin").Methods(http.MethodPost).HandlerFunc(apihandler.Gtin(connection))
	r.Path("/date").Methods(http.MethodPost).HandlerFunc(apihandler.Date(connection))
	//r.Path("/gtincor").Methods(http.MethodPost).HandlerFunc(apihandler.GtinandCor)
	r.Path("/").Methods(http.MethodGet).HandlerFunc(redirect)
	//r.Path("/index.html").Methods(http.MethodGet).HandlerFunc(homeHandler)

	server := &http.Server{
		Handler:      r,
		Addr:         ":" + c.Config.API.Port,
		WriteTimeout: 150 * time.Second,
		ReadTimeout:  150 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "http://" + req.Host + "/html/index.html"
	// log.Println(target)
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	// log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
		// see comments below and consider the codes 308, 302, or 301
		http.StatusTemporaryRedirect)
}

func apiDownloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("access-control-expose-headers", "*")
	w.Header().Set("Content-Type", "application/octet-stream")

	id := mux.Vars(r)["id"]
	if file := apihandler.GetFile(id); file != nil {
		w.Header().Set("File-Name", file.Name)
		io.Copy(w, file.Value)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}
