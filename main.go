package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type Mapper struct {
	Mapping map[string]string
	Lock sync.Mutex
}

var urlMapper Mapper

func init(){
	urlMapper = Mapper{
		Mapping: make(map[string]string),
	}
}


func main(){

	r:= chi.NewRouter()

	r.Use(middleware.Logger)
	
	r.Get("/" , func(w http.ResponseWriter , r *http.Request) {
		w.Write([]byte("Server is running..."))
	})

	r.Post("/short" , shortUrlHandler)
	r.Get("/redirect/{key}" , redirectHandler)


	http.ListenAndServe(":8080" , r)
}


func shortUrlHandler(w http.ResponseWriter , r *http.Request){
	r.ParseForm()
	u:= r.Form.Get("url")
	
	if u == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("url files is empty."))
		return

	}

	uuid := uuid.New().String()

	insertMapping(uuid , u)

	w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("http://localhost:8080/redirect/%s" , uuid)))
}


func redirectHandler(w http.ResponseWriter , r *http.Request){
	key:= chi.URLParam(r , "key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("key files is empty."))
		return
	}

	// fetch mapping 

	u:= fetchMapping(key)

	if u == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("key files is empty."))
		return

	}


	http.Redirect(w , r , u , http.StatusFound)
}


func insertMapping (key  , url string){
	urlMapper.Lock.Lock()
	defer urlMapper.Lock.Unlock()

	urlMapper.Mapping[key] = url
	fmt.Println("indert" , urlMapper.Mapping)
}


func fetchMapping(key string) string {
	urlMapper.Lock.Lock()
	defer urlMapper.Lock.Unlock()

	println("fetch" , urlMapper.Mapping[key])
	return urlMapper.Mapping[key]
}
