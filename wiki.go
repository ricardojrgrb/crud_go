package main

import (
	//"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
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

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil { 
    	http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    	return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	
	log.Fatal(http.ListenAndServe(":8080", nil))

	/*
		the downs call is not use, they are by commented downs code
	*/

	// http.HandleFunc("/view/", viewHandler)
	// log.Fatal(http.ListenAndServe(":8080", nil))

	// http.HandleFunc("/", handler)
	// log.Fatal(http.ListenAndServe(":8080", nil))

	// p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	// p1.save()
	// p2, _ := loadPage("TestPage")
	// fmt.Println(string(p2.Body))
}

/*
	the downs code is not use at mommet, isn`t optimized
*/

// func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
// 	m := validPath.FindStringSubmatch(r.URL.Path)
// 	if m == nil {
// 		http.NotFound(w, r)
// 		return "", errors.New("Invalid Page Title")
// 	}
// 	return m[2], nil // The title is the second subexpression
// }

//	func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
// 	t, err := template.ParseFiles(tmpl + ".html")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusIntervalServerError)
// 		return
// 	}
// 	err = t.Execute(w, p)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusIntervalServerError)
// 	}
// }

//	func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
// 	t, _ := template.ParseFiles(tmpl + ".html")
// 	t.Execute(w, p)	
// }	

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }

// func viewHandler(w http.ResponseWriter, r *http.Request) {
//     title := r.URL.Path[len("/view/"):]
//     p, _ := loadPage(title)
//     fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)	
// }

// func editHandler(w http.ResponseWriter, r *http.Request) {
// 	title := r.URL.Path[len("/edit/"):]
// 	p, err := loadPage(title)
// 	if err != nil {
// 		p = &Page{Title: title}
// 	}
// 	fmt.Fprintf(w, "<h1>Editing %s</h1>")+
// 	"<form action=\"/save/%s\" method=\"POST\">"+
// 	"<textarea name=\"body\">%s</textarea><br>"+
// 	"<input type=\"submit\" value=\"Save\">"+
// 	"</form>",
// 	p.Title, p.Title, p.Body)
// }

// func viewHandler(w http.ResponseWriter, r *http.Request) {
//     title := r.URL.Path[len("/view/"):]
//     p, _ := loadPage(title)
//     t, _ := template.ParseFiles("view.html")
//     t.Execute(w,p)
// }

// func editHandler(w http.ResponseWriter, r *http.Request) {
// 	title := r.URL.Path[len("/edit/"):]
// 	p, err := loadPage(title)
// 	if err != nil {
// 		p = &Page{Title: title}
// 	}
// 	t, _ := template.ParseFiles("edit.html")
// 	t.Execute(w, p)
// }



