package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/eaglexiang/mcol/cache"
	"github.com/eaglexiang/mcol/config"
	"github.com/eaglexiang/mcol/env"
	"github.com/eaglexiang/mcol/mongo"
	"github.com/pkg/errors"
)

//go:embed mcol.config.default
var defaultConfig []byte

func main() {
	ctx := context.Background()

	if len(os.Args) < 2 {
		return
	}

	var err error
	switch os.Args[1] {
	case "--cache":
		err = createCache(ctx)
	case "--config":
		err = editConfig()
	default:
		err = cache.Load()
	}
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

func editConfig() (err error) {
	filename, err := config.Try2InitFile(defaultConfig)
	if err != nil {
		return
	}
	fmt.Println("文件路径: ", filename)

	editor := getEditor()

	err = runEditor(editor, filename)

	return
}

func getEditor() (editor string) {
	editor = env.DefaultEditor()
	fmt.Printf("请输入编辑器命令(默认: %s): ", editor)
	fmt.Scanf("%s", &editor)
	return
}

func runEditor(editor string, filename string) (err error) {
	fmt.Printf("%s %s\n", editor, filename)

	cmd := exec.Command(editor, filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}
