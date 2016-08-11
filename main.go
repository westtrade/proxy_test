package main

import (
	"bytes"
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

func initConfig() (string, []byte, []byte) {

	target := flag.String("target", "http://habrahabr.ru/", "Proxy target")
	search := flag.String("search", "", "Seeking sentence")
	replace := flag.String("replace", "", "Replace sentence")

	flag.Parse()

	return *target, []byte(*search), []byte(*replace)
}

func getMirrorData(targetLink string) ([]byte, http.Header) {

	resp, err := http.Get(targetLink)
	panicWrapper(err)
	defer resp.Body.Close()
	bodyData, err := ioutil.ReadAll(resp.Body)

	panicWrapper(err)
	return bodyData, resp.Header
}

func main() {
	var serverPort int = 3001
	var serverHost string = ""
	var listenAddress string = serverHost + ":" + strconv.Itoa(serverPort)
	target, search, replace := initConfig()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		outsidePath := "http://" + path.Clean(strings.Replace(target, "http://", "", 1)+r.URL.Path)
		bodyData, reqHeaders := getMirrorData(outsidePath)

		w.Header().Add("Content-Type", reqHeaders.Get("Content-Type"))
		bodyData = bytes.Replace(bodyData, search, replace, -1)
		w.Write(bodyData)

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
