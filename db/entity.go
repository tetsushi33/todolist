package db

// schema.go provides data models in DB
import (
	"time"
)

// Task corresponds to a row in `tasks` table
/*　`db:"id"`などの部分がタグ　データベース上での属性名を指定
　　IDはタスクテーブル内のid属性に対応する*/
type Task struct {
	ID        uint64    `db:"id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
	IsDone    bool      `db:"is_done"`
	Message   string    `db:"message"`
}
