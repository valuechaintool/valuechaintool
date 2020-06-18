package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"

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
	router.HTMLRender = web.LoadTemplates("static/tpl")
	router.Use(static.Serve("/static", static.LocalFile("static", false)))
	if viper.GetBool("tls.enabled") {
		router.Use(web.MiddlewareHSTS())
	}

	// Authentication bits
	router.GET("/login", web.Login)
	router.POST("/login", web.LoginPost)
	router.GET("/logout", web.MiddlewareAuth(), web.Logout)
	router.GET("/register", web.Register)
	router.POST("/register", web.RegisterPost)

	// Application bits
	r := router.Group("/")
	r.Use(web.MiddlewareAuth())
	r.GET("/", web.Home)
	r.GET("/companies", web.CompaniesList)
	r.GET("/companies/:id", web.CompaniesRead)
	r.GET("/companies/:id/changelog", web.CompaniesChangelog)
	r.GET("/companies/:id/new", web.CompaniesCreate)
	r.POST("/companies/:id/new", web.CompaniesCreatePost)
	r.POST("/companies/:id/edit", web.CompaniesUpdatePost)
	r.GET("/companies/:id/delete", web.CompaniesDelete)
	r.POST("/companies/:id/relationships", web.RelationshipsCreatePost)
	r.POST("/companies/:id/relationships/:rid", web.RelationshipsUpdate)
	r.GET("/companies/:id/relationships/:rid/delete", web.RelationshipsDelete)

	// Users management bits
	r.GET("/users", web.UsersList)
	r.GET("/users/:user_id", web.UsersRead)
	r.GET("/users/:user_id/delete", web.UsersDelete)
	r.POST("/users/:user_id/permissions", web.PermissionsCreatePost)
	r.GET("/users/:user_id/permissions/:permission_id/delete", web.PermissionsDelete)

	// Health
	router.GET("/healthz", web.Healthz)
	router.GET("/readiness", web.Readiness)

	log.Println("ready to serve")

	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}

	srvConfig := &http.Server{
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           router,
	}

	if viper.GetBool("tls.enabled") {
		srvConfig.Addr = fmt.Sprintf(":%d", viper.GetInt("tls.port"))
		srvConfig.TLSConfig = tlsConfig
		go http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("port")), http.HandlerFunc(web.Redirect))
		if err := srvConfig.ListenAndServeTLS(viper.GetString("tls.certificate"), viper.GetString("tls.key")); err != nil {
			log.Println(err)
		}
	} else {
		srvConfig.Addr = fmt.Sprintf(":%d", viper.GetInt("port"))
		if err := srvConfig.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}
}
