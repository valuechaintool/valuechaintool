package main

import (
	"fmt"
	"log"

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

	// Authentication bits
	router.GET("/login", web.Login)
	router.POST("/login", web.LoginPost)
	router.GET("/logout", web.MiddlewareAuth(), web.Logout)

	// Application bits
	r := router.Group("/")
	r.Use(web.MiddlewareAuth())
	r.GET("/", web.Home)
	r.GET("/companies", web.CompaniesList)
	r.GET("/companies/:id", web.CompaniesRead)
	r.GET("/companies/:id/new", web.CompaniesCreate)
	r.POST("/companies/:id/new", web.CompaniesCreatePost)
	r.GET("/companies/:id/edit", web.CompaniesUpdate)
	r.POST("/companies/:id/edit", web.CompaniesUpdatePost)
	r.GET("/companies/:id/delete", web.CompaniesDelete)
	r.POST("/companies/:id/relationships", web.RelationshipsCreatePost)
	r.GET("/companies/:id/relationships/:rid/delete", web.RelationshipsDelete)

	// Health
	router.GET("/healthz", web.Healthz)
	router.GET("/readiness", web.Readiness)

	log.Println("ready to serve")
	if err := router.Run(":10080"); err != nil {
		log.Println(err)
	}
}
