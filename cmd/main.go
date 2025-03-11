package main

import(
	//"net/http"

	"net/http"
	"fmt"
	"urlshozim/config"
	"database/sql"
	"github.com/go-chi/chi/v5"
	"urlshozim/internal/repo"
	_ "github.com/lib/pq"// I would revise this place
)

func main(){
	//creating a router
	r := chi.NewRouter()
	_ =r

	//database connection
	cfg := config.Configure()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User,cfg.Password, cfg.DBname)
	fmt.Println(psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	defer db.Close()
	if err!=nil {
		//cheta tam cheta tam
		fmt.Println("database creation\n")
		panic(err)
	}
	err = db.Ping()
	if err!=nil {
		//fmt.Println("here")
		fmt.Println("database check\n")
		panic(err)
	}
	
	aa := repo.Storage{Db: db}

	r.Post("/url", aa.Create())
	r.Get("/{alias}", aa.Getter())

	err = http.ListenAndServe(":8080", r)
	if err!=nil{
		fmt.Println("Server linking\n")
		panic(err)
	}
}