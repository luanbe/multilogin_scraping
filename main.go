package main

import (
	"fmt"
	"github.com/luanbe/golang-web-app-structure/initialization"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func init() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool("debug") {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	// Init Db connection
	db, err := initialization.InitDb()
	if err != nil {
		fmt.Errorf("error Db connection: %v", err.Error())
		panic(err.Error())
	}

	// Int Session Manager
	sessionManager := initialization.IntSessionManager()

	// Int Router
	router := initialization.InitRouting(db, sessionManager)

	fmt.Printf("Server START on port%v\n", viper.GetString("server.address"))
	log.Fatal(http.ListenAndServe(
		viper.GetString("server.address"),
		sessionManager.LoadAndSave(router),
	))
}
