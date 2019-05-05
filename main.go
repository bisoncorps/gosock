/*
 Simple Chat application using sockets in GoLang
*/


package main

import (
	"flag"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
	"strings"
	"github.com/deven96/gosock/pkg/custlog"
)

// declare global variables for use throughout the main package
var TemplateDir = filepath.Join("templates")
var AssetsDir = filepath.Join("assets")
// command line arguments and defaults
var LogFile = flag.String("log", "gosock.log", "Name of the log file to save to")
var ServerLocation = flag.String("addr", ":8008", "The addr of the application.")

// templateHandler represents a single template
type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}


// ServeHTTP handles the HTTPRequest
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func(){
		filearr := []string{TemplateDir, t.filename}
		filepath := strings.Join(filearr, "/")
		t.templ = template.Must(template.ParseFiles(filepath))
	})
	t.templ.Execute(w, r)
}

func main() {
	flag.Parse() // parse the flags
	defwriters := custlog.DefaultWriters(*LogFile, false)
	//TRACE will be Discarded, while the rest will be routed accordingly
	custlog.LogInit(defwriters)	
	custlog.Trace.Println("Imported Custom Logging")
	custlog.Info.Println("Log file can be found at ", custlog.Logfile)
	//os.Setenv("GOSOCK_LOG", custlog.Logfile)

	// create a room
	r := newRoom()
	
	/* Routes */
	// Handle function for route "/"
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir(AssetsDir))))
	


	
	//start the room
	custlog.Info.Println("Initializing Room...")
	go r.run()
	//start the webserver
	custlog.Info.Printf("Running server started on %s", *ServerLocation)
	
	if err := http.ListenAndServe(*ServerLocation, nil); err != nil {
		custlog.Error.Println(err)
	}
}
