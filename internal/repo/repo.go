package repo

import(
	"fmt"
	"database/sql"
	"net/http"
	"encoding/json"
)

type Storage struct {
	Db *sql.DB
}

type Request struct{
	Url string		`json:"url"`
	Alias string	`json:"alias"`
}

func (stor *Storage) Create() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		var fuck Request
		fmt.Println(r.Body)
		//fmt.Println(r.Body.RemoteAddr)
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&fuck)
		fmt.Println(fuck.Url)
		if err!=nil {
			fmt.Println("jsson decoding\n")
			panic(err)
		}
		_, err = stor.Db.Exec("INSERT INTO urls (url, alias) VALUES ($1, $2)", fuck.Url, fuck.Alias)
		if err != nil {
			fmt.Println("db exec")
			panic(err)
		}
		w.Write([]byte("Success"))
	});
}

func (stor *Storage) Getter() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		alias := r.URL.Path
		alias = alias[1:]
		var url string
		err := stor.Db.QueryRow("SELECT url FROM urls WHERE alias = $1", alias).Scan(&url)
		if err!=nil {
			fmt.Println("getter\n")
			panic(err)
		}
		http.Redirect(w,r,url,http.StatusSeeOther)
	});
}