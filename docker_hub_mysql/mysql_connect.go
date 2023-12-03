package main

import (
	"database/sql"
	"fmt"

	// 下のコマンドでライブラリをインストールする
	// * go get -u github.com/go-sql-driver/mysql
	_ "github.com/go-sql-driver/mysql"
)

type All struct {
	Id   int
	Name string
}

type One struct {
	ID   int
	NAME string
}

func connectionDB() *sql.DB {
	// "[ ユーザー ]:[ パスワード ]@/[ データベース名 ]"
	dsn := "root:docker@tcp(localhost:13306)/docker_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func getAll(db *sql.DB) *sql.Rows {
	all, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println("データベース接続失敗-データを取得できません")
		panic(err.Error())
	} else {
		fmt.Println("データベース接続成功-データを取得を行います")
	}

	return all
}

func addData(db *sql.DB) *sql.Stmt {
	// プリペアードステートメントを使用
	in, err := db.Prepare("INSERT INTO users(name) VALUES(?)")
	if err != nil {
		fmt.Println("データベース接続失敗-データを追加できません")
		panic(err.Error())
	} else {
		fmt.Println("データベース接続成功-データを追加を行います")
	}

	return in
}

func updateData(db *sql.DB) *sql.Stmt {
	up, err := db.Prepare("UPDATE users SET name=? WHERE id=?")
	if err != nil {
		fmt.Println("データベース接続失敗-データを更新できません")
		panic(err.Error())
	} else {
		fmt.Println("データベース接続成功-データを更新を行います")
	}

	return up
}

func deleteData(db *sql.DB) *sql.Stmt {
	del, err := db.Prepare("DELETE FROM users WHERE id=?")
	if err != nil {
		fmt.Println("データベース接続失敗-データを削除できません")
		panic(err.Error())
	} else {
		fmt.Println("データベース接続成功-データを削除を行います")
	}

	return del
}

func main() {
	db := connectionDB()

	defer db.Close()

	// 一行データ取得
	one := One{}
	err := db.QueryRow("SELECT * FROM users WHERE id = ?", 2).Scan(&one.ID, &one.NAME)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(one.ID, one.NAME)

	in := addData(db)
	// データベースに追加するデータ
	result_in, err := in.Exec("goro")
	if err != nil {
		panic(err.Error())
	}
	lastId, err := result_in.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(lastId)

	up := updateData(db)
	result_up, err := up.Exec("一郎", 1)
	rowsAffect_up, err := result_up.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(rowsAffect_up)

	del := deleteData(db)
	result_del, err := del.Exec(6)
	rowsAffect_del, err := result_del.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(rowsAffect_del)

	rows := getAll(db)
	// 行データ取得
	all := All{}
	var result_sel []All
	for rows.Next() {
		error := rows.Scan(&all.Id, &all.Name)
		if error != nil {
			fmt.Println("スキャン失敗")
		} else {
			result_sel = append(result_sel, all)
		}
	}
	fmt.Println(result_sel)
}