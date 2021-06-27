package mongo

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var sysDBSet = map[string]struct{}{
	"admin":  {},
	"config": {},
	"local":  {},
}

func createClient(tx context.Context, addr string, db string, username string, pswd string) (client *mongo.Client, err error) {
	// 初始化 mongo client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化连接参数
	var auth string
	if username != "" && pswd != "" {
		auth = username + ":" + pswd + "@"
	}
	uri := "mongodb://" + auth + addr + "/" + db
	opt := options.Client().ApplyURI(uri)

	// 发起连接
	log.Println("dial: ", uri)
	client, err = mongo.Connect(ctx, opt)
	if err != nil {
		err = errors.WithMessage(err, addr)
		err = errors.WithStack(err)
		return
	}

	// 验证连接
	log.Println("ping")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

func isSysDB(dbName string) (is bool) {
	_, is = sysDBSet[dbName]
	return
}

// GetAllDBsAndCols 从 mongo 获取所有的 db 和 col
func GetAllDBsAndCols(ctx context.Context, addr string, db string, username string, pswd string) (dbName2colNames map[string][]string, dbCount int, colCount int, err error) {
	dbName2colNames = make(map[string][]string)

	// 创建连接
	client, err := createClient(ctx, addr, db, username, pswd)
	if err != nil {
		return
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// 拉取所有 db
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbNames, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	// 拉取所有 col
	for _, dbName := range dbNames {
		if isSysDB(dbName) {
			continue
		}

		dbCount++

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var colNames []string
		colNames, err = client.Database(dbName).ListCollectionNames(ctx, bson.M{})
		if err != nil {
			err = errors.WithMessage(err, dbName)
			err = errors.WithStack(err)
			return
		}

		dbName2colNames[dbName] = colNames
		colCount += len(colNames)

		time.Sleep(time.Millisecond * 100)
	}

	return
}
