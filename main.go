package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/urfave/negroni"

	"github.com/valuechaintool/valuechaintool/models"
	"github.com/valuechaintool/valuechaintool/web"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("starting")

	// Init DB
	if len(os.Getenv("DBC")) == 0 {
		panic("DBC has to be set in the environment")
	}
	session, err := gorm.Open("postgres", os.Getenv("DBC"))
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.DB().SetMaxOpenConns(5)

	if err := models.Init(session); err != nil {
		panic("Error while initializing the Internals module")
	}

	// Router
	r := mux.NewRouter()

	// Webpages
	r.HandleFunc("/", web.Home).Methods("GET")
	r.HandleFunc("/companies", web.CompaniesList).Methods("GET")
	r.HandleFunc("/companies/new", web.CompaniesCreate).Methods("GET")
	r.HandleFunc("/companies/new", web.CompaniesCreatePost).Methods("POST")
	r.HandleFunc("/companies/{id:[a-z0-9/-]{36}}", web.CompaniesRead).Methods("GET")
	r.HandleFunc("/companies/{id:[a-z0-9/-]{36}}/edit", web.CompaniesUpdate).Methods("GET")
	r.HandleFunc("/companies/{id:[a-z0-9/-]{36}}/edit", web.CompaniesUpdatePost).Methods("POST")
	r.HandleFunc("/companies/{id:[a-z0-9/-]{36}}/delete", web.CompaniesDelete).Methods("GET")

	// Static contents
	css := http.StripPrefix("/css/", http.FileServer(http.Dir("./css/")))
	r.PathPrefix("/css/").Handler(css)
	js := http.StripPrefix("/js/", http.FileServer(http.Dir("./js/")))
	r.PathPrefix("/js/").Handler(js)
	wf := http.StripPrefix("/webfonts/", http.FileServer(http.Dir("./webfonts/")))
	r.PathPrefix("/webfonts/").Handler(wf)

	// Health
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("ok")); err != nil {
			log.Println(err)
			return
		}
	}).Methods("GET")
	r.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("ok")); err != nil {
			log.Println(err)
			return
		}
	}).Methods("GET")

	// Negroni
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(r)

	log.Println("ready to serve")
	if err := http.ListenAndServe(":10080", n); err != nil {
		log.Println(err)
	}
}
