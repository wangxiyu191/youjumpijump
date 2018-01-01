package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

var basePath string
var files map[string]*os.File

func init() {
	files = map[string]*os.File{}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	basePath = filepath.Dir(ex)
}

func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rval := recover(); rval != nil {
				debug.PrintStack()
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		log.Printf("%s %s %s %s\n", r.RemoteAddr, r.Method, r.URL, elapsed)
	}
	return http.HandlerFunc(fn)
}

func getFile(fileName string) *os.File {
	if file, ok := files[fileName]; ok {
		return file
	}
	file, _ := os.OpenFile(basePath+"/"+fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	return file
}

func main() {
	router := httprouter.New()
	router.GET("/:phoneos/:width/:height", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		body, _ := ioutil.ReadFile(fmt.Sprintf("%s/%s-%sx%s.txt", basePath, ps.ByName("phoneos"), ps.ByName("width"), ps.ByName("height")))
		w.Write(body)
	})
	router.POST("/:phoneos/:width/:height", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		file := getFile(fmt.Sprintf("%s-%sx%s.txt", ps.ByName("phoneos"), ps.ByName("width"), ps.ByName("height")))
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			line := strings.Split(string(body), ",")
			if len(line) == 2 {
				distance, err1 := strconv.ParseFloat(line[0], 64)
				ratio, err2 := strconv.ParseFloat(line[1], 64)
				if err1 == nil && err2 == nil {
					file.Write([]byte(fmt.Sprintf("%v,%v\n", distance, ratio)))
				}
			}

		}
		w.Write([]byte{})
	})

	log.Fatal(http.ListenAndServe(":8213", Middleware(router)))
}
