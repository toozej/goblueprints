package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	// uncomment the following two imports when enabling tracing
	// "os"
	//"goblueprints/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// serveHTTP handles the HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	// setup port address for application web server
	var addr = flag.String("addr", ":3000", "The address of the application")
	flag.Parse()

	// set up gomniauth
	gomniauth.SetSecurityKey("iuehfwiufh3297t4g39ios")
	gomniauth.WithProviders(
		facebook.New("key",
			"secret",
			"http://localhost:3000/auth/callback/facebook"),
		github.New("key",
			"secret",
			"http://localhost:3000/auth/callback/github"),
		google.New("930300705643-lph2h3dngnbehq56bonkvcv8c1hq9cm1.apps.googleusercontent.com",
			"NxrcNrk9xDLbHsrRg6jutAMN",
			"http://localhost:3000/auth/callback/google"),
	)

	// create a new room
	r := newRoom(UseGravatar)

	// uncomment the following line to enable tracing
	// r.tracer = trace.New(os.Stdout)

	// setup handles
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})

	// get the room going
	go r.run()

	//start the web server
	log.Println("Starting the web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
