package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
)

// Примечание: Названия у функций начинаются со строчной буквы  по причине того,
// что они используются только в текущем пакете и нет смысла делать их публичными

func panicWrapper(err error) {
	if err != nil {
		panic(err)
	}
}

func initConfig() (string, string, string) {

	target := flag.String("target", "http://habrahabr.ru/", "Proxy target")
	search := flag.String("search", "", "Seeking sentence")
	replace := flag.String("replace", "", "Replace sentence")

	flag.Parse()

	return *target, *search, *replace
}

func getRequestString(targetLink string) string {

	resp, err := http.Get(targetLink)
	panicWrapper(err)
	defer resp.Body.Close()
	bodyData, err := ioutil.ReadAll(resp.Body)
	panicWrapper(err)
	result := string(bodyData)
	return result
}

func main() {
	var serverPort int = 3001
	var serverHost string = ""
	var listenAddress string = serverHost + ":" + strconv.Itoa(serverPort)
	target, search, replace := initConfig()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		outsidePath := "http://" + path.Clean(strings.Replace(target, "http://", "", 1)+r.URL.Path)
		result := getRequestString(outsidePath)
		result = strings.Replace(result, search, replace, -1)

		fmt.Fprint(w, result)
	})

	fmt.Printf(`
        Server started %s
        Target host: %s
        Seeking sentence: %s
        Replace for next sentence: %s

    `, listenAddress, target, search, replace)

	err := http.ListenAndServe(listenAddress, nil)
	panicWrapper(err)
}
