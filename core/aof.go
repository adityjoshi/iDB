package core

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/adityjoshi/iDB/config"
)

func dumpKey(fp *os.File, key string, obj *Object) {
	cmd := fmt.Sprintf("SET %s %s", key, obj.value)
	tokens := strings.Split(cmd, " ")
	fp.Write(Encode(tokens, false))
}

func DUMPALLAOF() {
	fp, err := os.OpenFile(config.AOFFILE, os.O_CREATE|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		fmt.Print("error", err)
		return
	}

	log.Println("rewriting AOF file at", config.AOFFILE)
	for k, obj := range store {
		dumpKey(fp, k, obj)
	}
	log.Println("AOF file rewrite complex")
}
