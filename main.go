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

	// Webpages
	router.GET("/", web.Home)
	router.GET("/company/new", web.CompaniesCreate)
	router.POST("/company/new", web.CompaniesCreatePost)
	router.GET("/companies", web.CompaniesList)
	router.GET("/companies/:id", web.CompaniesRead)
	router.GET("/companies/:id/edit", web.CompaniesUpdate)
	router.POST("/companies/:id/edit", web.CompaniesUpdatePost)
	router.GET("/companies/:id/delete", web.CompaniesDelete)
	router.POST("/companies/:id/relationships", web.RelationshipsCreatePost)
	router.GET("/companies/:id/relationships/:rid/delete", web.RelationshipsDelete)

	// Health
	router.GET("/healthz", web.Healthz)
	router.GET("/readiness", web.Readiness)

	log.Println("ready to serve")
	if err := router.Run(":10080"); err != nil {
		log.Println(err)
	}
}
