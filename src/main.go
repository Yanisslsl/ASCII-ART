package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string

	Body []byte
}

func (p *Page) save() error {

	filename := p.Title + ".txt"

	return ioutil.WriteFile(filename, p.Body, 0600)

}

func loadPage(title string) (*Page, error) {

	filename := title + ".txt"

	body, err := ioutil.ReadFile(filename)

	if err != nil {

		return nil, err

	}

	return &Page{Title: title, Body: body}, nil

}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func serv(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "form.html")
	case "POST":

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		name := r.FormValue("ascci-input")
		typ := r.FormValue("fonts")
		Arg := []rune(name)
		var nbline int
		var line int
		var output string

		for index1 := 0; index1 < 8; index1++ {
			for index2 := 0; index2 < len(Arg); index2++ {
				nbline = 0
				line = getLine(Arg[index2])
				file, erreur := os.Open(typ + ".txt")

				if erreur != nil {
					errortake(erreur.Error())
				} else {
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						if nbline == line+index1 {

							output += scanner.Text()

						}
						nbline++
					}
				}

			}

			output += ("\n")

		}
		p := &Page{Title: output}
		tmpl, _ := template.ParseFiles("form.html")
		tmpl.Execute(w, p)

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func getLine(char rune) int {
	var line int
	for index := 0; index < 95; index++ {
		if rune(index+32) == char {
			line = index
			break
		}
	}
	line = line*9 + 1
	return line

}

func errortake(str string) {
	fmt.Println("ERREUR: " + str)
}

func main() {
	http.HandleFunc("/", serv)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
