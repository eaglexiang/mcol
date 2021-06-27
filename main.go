package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/eaglexiang/mcol/cache"
	"github.com/eaglexiang/mcol/config"
	"github.com/eaglexiang/mcol/mongo"
)

func main() {
	ctx := context.Background()

	if len(os.Args) < 2 {
		return
	}

	if os.Args[1] == "--cache" {
		log.Println("start to cache")
		err := createCache(ctx)
		if err != nil {
			log.Printf("%+v", err)
		}
		return
	}

	err := cache.Load()
	if err != nil {
		log.Printf("%+v", err)
		return
	}

	keys := os.Args[1:]
	results := cache.Search(keys)
	for _, result := range results {
		format, args := result.Fmt()
		fmt.Printf(format+"\n", args...)
	}
}

func createCache(ctx context.Context) (err error) {
	err = config.Load()
	if err != nil {
		return
	}

	c, dbCount, colCount, err := mongo.GetAllDBsAndCols(ctx, config.C.Addr, config.C.DB, config.C.UserName, config.C.Password)
	if err != nil {
		return
	}
	log.Printf("%d dbs and %d cols cached", dbCount, colCount)

	err = cache.Save(c)

	return
}
