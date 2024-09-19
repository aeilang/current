package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "postgres://lang:password@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("PUT /student/{id}", updateStudent)

	http.ListenAndServe(":8888", mux)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	// 提取id
	id := r.PathValue("id")

	if len(id) == 0 {
		http.Error(w, "路径id未提供", http.StatusBadRequest)
		return
	}

	// 把id 转未int类型
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "错误的id", http.StatusBadRequest)
		return
	}

	// 查找id对应的学生
	query := `SELECT id, name, age, version from students where id = $1;`
	row := db.QueryRowContext(context.Background(), query, idInt)

	var s Student
	if err := row.Scan(&s.Id, &s.Name, &s.Age, &s.Version); err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// 从r.Body中反序列要更改的数据

	type input struct {
		Name *string `json:"name"`
		Age  *int    `json:"age"`
	}

	var in input
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "反序列化错误", http.StatusInternalServerError)
		return
	}

	// 对照in 更改学生数据

	if in.Age != nil {
		s.Age = *in.Age
	}

	if in.Name != nil {
		s.Name = *in.Name
	}

	// 将更改后的学生保存的数据库
	update := `update students set name = $1, age = $2, version = version + 1 where id = $3 and version = $4;`
	result, err := db.ExecContext(context.Background(), update, s.Name, s.Age, s.Id, s.Version)
	if err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		http.Error(w, "更新失败", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("更改成功\n"))
}

type Student struct {
	Id      string
	Name    string
	Age     int
	Version int
}
