package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"html/template"
	"log"
	"main.go/domain"
	"net/http"
)

var (
	rd  *render.Render
	db  *gorm.DB
	tpl *template.Template
)

//main 함수
//1. render 변수 선언 : html 확장자 옵션 처리
//2. 사용자 핸들러 함수 담기
//3. negroni 기본 핸들러 선언 + 사용자 핸들러 담기
func main() {
	rd = render.New(render.Options{ //-- 1
		Directory:  "view/templates",
		Extensions: []string{".html", ".tmpl"},
	})
	mux := MakeWebHandler() //-- 2

	n := negroni.Classic() //-- 3
	n.UseHandler(mux)

	//mysql 연결
	db, err := ConnectDB()
	if err != nil {
		err := fmt.Errorf("연결실패 : %v", err)
		log.Println(err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()
	//테이블 생성
	if err := db.AutoMigrate(&domain.Board{}); err != nil {
		fmt.Println("Board table Err")
	} else {
		fmt.Println("Board table Suc")
	}

	log.Println("Started App")
	err = http.ListenAndServe(":3000", n)
	if err != nil {
		panic(err)
	}
}

//사용자 핸들러 함수
//return 값 : 핸들러 인스턴스
func MakeWebHandler() http.Handler {
	router := mux.NewRouter()

	router.Handle("/static/{dir}/{file}", http.StripPrefix("/static/", http.FileServer(http.Dir("view/static/"))))
	router.HandleFunc("/", MainHandler).Methods("GET")
	router.HandleFunc("/write", WritePageHandler).Methods("GET")
	router.HandleFunc("/write", writeHandler).Methods("POST")
	//router.HandleFunc("/read", controller.readAllHandler).Methods("GET")
	//router.HandleFunc("read/{id}", controller.ReadHandler).Methods("GET")
	return router
}

func ConnectDB() (*gorm.DB, error) {
	dsn := "root:root@tcp(10.28.3.180:3307)/jamie?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	return db, err
}

//메인
func MainHandler(w http.ResponseWriter, r *http.Request) {
	rd.HTML(w, http.StatusOK, "home", nil)
}

//글 작성 페이지
func WritePageHandler(w http.ResponseWriter, r *http.Request) {
	rd.HTML(w, http.StatusOK, "write", nil)
}

//글 작성
func writeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	if r.Method == http.MethodPost {
		title := r.PostFormValue("title")
		author := r.PostFormValue("author")
		content := r.PostFormValue("content")

		newPost := domain.Board{Title: title, Author: author, Content: content}
		db.Create(&newPost)
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}
	tpl.ExecuteTemplate(w, "write", nil)
}

//특정 글 조회
func ReadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["id"]
	fmt.Println("id는 ", Id)
	fmt.Fprintf(w, "게시판 입니다. 게시글 %s 번", Id)
}
