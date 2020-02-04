package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/zserge/webview"
)

const (
	windowWidth  = 480
	windowHeight = 320
)

var indexHTML = `
<!doctype html>
<html>
	<head>
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
	</head>
	<style>
	#result {
		text-align:center;
		font-size:24px;
	}
	</style>
	<body>
		
		<input id="nome" value="" name="nome" type="number" />
		<input id="up" value="1" type="checkbox" />Uppercase
		<input id="down" value="1" type="checkbox" />Lowercase
		<button onclick="external.invoke(document.getElementById('nome').value +',' 
		+ document.getElementById('up').checked + ',' +  document.getElementById('down').checked )">
			Generate
		</button>
		<div id="result"></div>
	</body>
</html>
`

func randSeq(n int, up string, down string) string {
	var lettersUP = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var lettersDoWN = []rune("abcdefghijklmnopqrstuvwxyz")
	var lettersAll = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	letters := lettersUP
	if up == "true" && down == "false" {
		letters = lettersUP
	}
	if up == "false" && down == "true" {
		letters = lettersDoWN
	}

	if up == "true" && down == "true" {
		letters = lettersAll
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func startServer() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer ln.Close()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(indexHTML))
		})
		log.Fatal(http.Serve(ln, nil))
	}()
	return "http://" + ln.Addr().String()
}

func handleRPC(w webview.WebView, data string) {

	s := strings.Split(data, ",")
	dimension := s[0]
	up := s[1]
	down := s[2]

	i, num := strconv.Atoi(dimension)
	fmt.Println(num)

	password := randSeq(i, up, down)
	fmt.Println(password)

	w.Eval(`document.getElementById('result').innerHTML="` + password + `" ;`)
}

func main() {
	url := startServer()
	w := webview.New(webview.Settings{
		Width:                  windowWidth,
		Height:                 windowHeight,
		Title:                  "Generate Random Password",
		Resizable:              true,
		URL:                    url,
		ExternalInvokeCallback: handleRPC,
	})
	w.SetColor(255, 255, 255, 255)
	defer w.Exit()
	w.Run()
}
