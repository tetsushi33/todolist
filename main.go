package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"todolist.go/db"
	"todolist.go/service"
)

const port = 8000

func main() {
	// initialize DB connection
	/*
		MySQLサーバとの接続に必要な設定を環境変数から読み込む処理．
		docker-compose.yml 内で Docker コンテナ上における環境変数それぞれの値を設定している。
		ローカルで動かしたい場合には，適切な環境変数を設定することでプログラムを変更することなく接続先データベースを変更可能
	*/
	dsn := db.DefaultDSN(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	//データベースとの接続　失敗したらlog.Fatal関数を呼び出す
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}

	// initialize Gin engine
	engine := gin.Default()
	engine.LoadHTMLGlob("views/*.html")

	// routing
	engine.Static("/assets", "./assets")
	//初期画面
	engine.GET("/", service.Home)
	//タスクの一覧画面
	engine.GET("/list", service.TaskList)
	//タスクの個別の情報表示画面
	engine.GET("/task/:id", service.ShowTask) // ":id"はパラメータ記法
	//タスクの新規登録画面
	engine.GET("task/new", service.NewTaskForm)
	engine.POST("task/new", service.RegisterTask)

	// start server
	engine.Run(fmt.Sprintf(":%d", port))
}
