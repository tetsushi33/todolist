package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)

// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Get tasks in DB
	var tasks []database.Task                      //型はdatabase.Taskのスライス
	err = db.Select(&tasks, "SELECT * FROM tasks") // Use DB#Select for multiple entries
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks})
}

// ShowTask renders a task with given ID
func ShowTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// parse ID given as a parameter
	//Ginのメソッド.Paramによりパスからしてのパラメータの値を取得する
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Get a task with given ID
	//パラメータから取得したIDの値をもとにデータベースからタスクを抽出する
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id) // Use DB#Get for one entry
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Render task
	ctx.HTML(http.StatusOK, "task_info.html", task)
}

// タスクの新規登録---------------------------------------------------------------------
// GET
func NewTaskForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "form_newTask.html", gin.H{"Title": "新規タスクの登録"})
}

// POST
func RegisterTask(ctx *gin.Context) {
	//POSTに送られているデータからtitleとmessageを取得
	title, exist := ctx.GetPostForm("title")
	if !exist {
		Error(http.StatusBadRequest, "No title is given")(ctx)
		return
	}
	message, exist := ctx.GetPostForm("message")
	if !exist {
		Error(http.StatusBadRequest, "No message is given")(ctx)
		return
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	//DBに追加データを挿入する
	result, err := db.Exec("INSERT INTO tasks (title, message) VALUES (?, ?)", title, message)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	//ページの遷移先の決定(リダイレクト)
	path := "/list" // デフォルトではタスク一覧ページへ戻る
	if id, err := result.LastInsertId(); err == nil {
		path = fmt.Sprintf("/task/%d", id) // 正常にIDを取得できた場合は /task/<id> へ戻る
	}
	ctx.Redirect(http.StatusFound, path)
}

// タスクの編集-------------------------------------------------------------------------
// GET
func EditTaskForm(ctx *gin.Context) {
	// ID の取得
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Get target task
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id)
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}
	// Render edit form
	ctx.HTML(http.StatusOK, "form_editTask.html",
		gin.H{"Title": fmt.Sprintf("タスクの編集　ID: %d", task.ID), "Task": task})
}

// POST
func UpdateTask(ctx *gin.Context) {
	// ID の取得(パスパラメータから)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	//POSTに送られたデータを取得(title, is_done, messageの3つ)
	title, exist := ctx.GetPostForm("title")
	if !exist {
		Error(http.StatusBadRequest, "No title is given")(ctx)
		return
	}
	is_done, exist := ctx.GetPostForm("is_done")
	if !exist {
		Error(http.StatusBadRequest, "No state is given")(ctx)
		return
	}
	is_done_bool, err0 := strconv.ParseBool(is_done)
	if err0 != nil {
		fmt.Printf("文字列\"%s\"を論理値に変換できませんでした\n", is_done)
	}
	message, exist := ctx.GetPostForm("message")
	if !exist {
		Error(http.StatusBadRequest, "No message is given")(ctx)
		return
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	//DB上の対象のテーブルの情報を更新する
	_, err1 := db.Exec("UPDATE tasks SET title = ?, is_done = ?, message = ? WHERE id = ?",
		title, is_done_bool, message, id)
	if err1 != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	//ページの遷移先 タスクの個別画面に戻る
	path := fmt.Sprintf("/task/%d", id)
	ctx.Redirect(http.StatusFound, path)
}
