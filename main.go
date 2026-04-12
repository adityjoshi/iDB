package main

import (
	"flag"
	"log"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for iDB")
	flag.IntVar(&config.Port, "Port", 1926, "port for thr iDB")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println(`
	_____ _____ ____  
  |_   _|  __ \  _ \ 
    | | | |  | | |_) |
    | | | |  | |  _ < 
   _| |_| |__| | |_) |
  |_____|_____/|____/ 
	
`)
	log.Println("If Steve Jobs ever built a database, he’d call it iDB. iDB started.")

}
