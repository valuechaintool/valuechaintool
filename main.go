package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"

	"github.com/valuechaintool/valuechaintool/api"
	"github.com/valuechaintool/valuechaintool/middleware"
	"github.com/valuechaintool/valuechaintool/models"
	"github.com/valuechaintool/valuechaintool/web"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("starting")

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("fatal error reading the config file: %s \n", err)
		return
	}

	// Init DB
	if len(viper.GetString("dbc")) == 0 {
		panic("DBC has to be set in the environment")
	}
	session, err := gorm.Open("postgres", viper.GetString("dbc"))
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.DB().SetMaxOpenConns(5)

	if err := models.Init(session); err != nil {
		fmt.Printf("fatal error while initializing the internal modules: %s \n", err)
		return
	}

	if viper.GetBool("DEBUG") {
		session.LogMode(true)
	}
	// Router
	router := gin.Default()
	router.POST("/api/v1/login", api.Login)
	router.POST("/api/v1/register", api.Register)
	router.GET("/api/v1/user", middleware.Auth(), api.User)
	router.GET("/api/v1/users", middleware.Auth(), api.UsersList)
	router.GET("/api/v1/users/:user", middleware.Auth(), api.UsersRead)
	router.PATCH("/api/v1/users/:user", middleware.Auth(), api.UsersUpdate)
	router.DELETE("/api/v1/users/:user", middleware.Auth(), api.UsersDelete)
	router.POST("/api/v1/users/:user/permissions", middleware.Auth(), api.PermissionsCreate)
	router.GET("/api/v1/users/:user/permissions", middleware.Auth(), api.PermissionsList)
	router.DELETE("/api/v1/users/:user/permissions/:permission", middleware.Auth(), api.PermissionsDelete)
	log.Fatal(router.Run(":8080"))

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
	r.HandleFunc("/relationships/new", web.RelationshipsCreate).Methods("GET")
	r.HandleFunc("/relationships/new", web.RelationshipsCreatePost).Methods("POST")
	r.HandleFunc("/relationships/{id:[a-z0-9/-]{36}}/delete", web.RelationshipsDelete).Methods("GET")

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
