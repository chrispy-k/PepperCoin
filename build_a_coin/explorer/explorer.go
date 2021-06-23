package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/chrispy-k/build_a_coin/blockchain"
	"github.com/gorilla/mux"
)

const (
	templateDir string = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	// to share things across classes we need upper cases
	Blocks []*blockchain.Block
}

// rw => writer, where we write the data that we want to send to user
// r => pointer to request, we dont want to copy request cause it can be big
func home(rw http.ResponseWriter, r *http.Request) {
	// this prints the data to the writer
	// fmt.Fprint(rw, "Hello from home!")

	// tmpl, err := template.ParseFiles("templates/home.html")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// "Must" does what we have done above for us!

	data := homeData{"Home", blockchain.Blockchain().Blocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.Blockchain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start(portNum int) {
	// we are loading pages and partials
	// we use newservemux because when we pass "nil" for handler,
	// listenandserve method uses default multiplexer and that clashes with
	// rest.go's multiplexer and causes error with handlefunc's "/"
	// need to change http.handlefunc to handler.handlefunc
	handler := mux.NewRouter()
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	fmt.Printf("Listening on http://localhost:%d\n", portNum)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portNum), handler))
}
