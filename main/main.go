package main

import (
	"log"

	"github.com/skowe/iniparser"
)

func main(){
	log.SetFlags(0)
	c := iniparser.NewINI("./conf.ini")
	c.Parse()
	
}