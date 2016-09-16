package myauthboss

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/authboss.v0"
	_ "gopkg.in/authboss.v0/auth"
	_ "gopkg.in/authboss.v0/confirm"
	_ "gopkg.in/authboss.v0/lock"
	aboauth "gopkg.in/authboss.v0/oauth2"
	_ "gopkg.in/authboss.v0/recover"
	_ "gopkg.in/authboss.v0/register"
	_ "gopkg.in/authboss.v0/remember"

	"github.com/aarondl/tpl"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
)

var funcs = template.FuncMap{
	"formatDate": func(date time.Time) string {
		return date.Format("2006/01/02 03:04pm")
	},
	"yield": func() string { return "" },
}

var (
	ab        = authboss.New()
	database  = NewMemStorer()
	templates = tpl.Must(tpl.Load("myauthboss/views", "myauthboss/views/partials", "layout.html.tpl", funcs))
)

func GetAuthboss() *authboss.Authboss {
	return ab
}

func SetupStorer() {
	cookieStoreKey, _ := base64.StdEncoding.DecodeString(`NpEPi8pEjKVjLGJ6kYCS+VTCzi6BUuDzU0wrwXyf5uDPArtlofn2AG6aTMiPmN3C909rsEWMNqJqhIVPGP3Exg==`)
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(`AbfYwmmt8UCwUuhd9qvfNA9UCuN1cVcKJN1ofbiky6xCyyBj20whe40rJa3Su0WOWLWcPpO1taqJdsEI/65+JA==`)
	cookieStore = securecookie.New(cookieStoreKey, nil)
	sessionStore = sessions.NewCookieStore(sessionStoreKey)
}

func SetupAuthboss() {
	ab.Storer = database
	ab.OAuth2Storer = database
	ab.MountPath = "/auth"
	ab.ViewsPath = "myauthboss/ab_views"
	ab.RootURL = `http://localhost:9000`

	ab.LayoutDataMaker = layoutData

	ab.OAuth2Providers = map[string]authboss.OAuth2Provider{
		"google": authboss.OAuth2Provider{
			OAuth2Config: &oauth2.Config{
				ClientID:     ``,
				ClientSecret: ``,
				Scopes:       []string{`profile`, `email`},
				Endpoint:     google.Endpoint,
			},
			Callback: aboauth.Google,
		},
	}

	b, err := ioutil.ReadFile(filepath.Join("myauthboss/views", "layout.html.tpl"))
	if err != nil {
		panic(err)
	}
	ab.Layout = template.Must(template.New("layout").Funcs(funcs).Parse(string(b)))

	ab.XSRFName = "csrf_token"
	ab.XSRFMaker = func(_ http.ResponseWriter, r *http.Request) string {
		return nosurf.Token(r)
	}

	ab.CookieStoreMaker = NewCookieStorer
	ab.SessionStoreMaker = NewSessionStorer

	//ab.Mailer = authboss.LogMailer(os.Stdout)
	ab.Mailer = authboss.SMTPMailer("smtp.ym.163.com:25", smtp.PlainAuth("", "noreply@xuntonglab.com", "7X11MLtiF6", "smtp.ym.163.com"))

	ab.Policies = []authboss.Validator{
		authboss.Rules{
			FieldName:       "email",
			Required:        true,
			AllowWhitespace: false,
		},
		authboss.Rules{
			FieldName:       "password",
			Required:        true,
			MinLength:       4,
			MaxLength:       8,
			AllowWhitespace: false,
		},
	}

	ab.RegisterOKPath = "/auth/login"
	ab.AuthLoginOKPath = "/minio/"
	ab.AuthLoginFailPath = "/auth/login"

	if err := ab.Init(); err != nil {
		log.Fatal(err)
	}
}

type Blog struct {
	ID       int
	Title    string
	AuthorID string
	Date     time.Time
	Content  string
}

type Blogs []Blog

var blogs = Blogs{
	{1, "My first Hook", "Stitches", time.Now().AddDate(0, 0, -3),
		`When I was young I was weak. But then I grew up and now I'm ridiculous` +
			`because I can hook a whole team and destroy them with gorge lololo. ` +
			`Look at me mom. I'm a big kid now!`,
	},
	{2, "Halp the nerfed!11", "Murky", time.Now().AddDate(0, 0, -1),
		`I used to be really amazing, then they nerfed me, then I was really` +
			`good in the right hands, now... I have no idea, why Blizzard? Why??!!`,
	},
}

func main2() {
	// Initialize Sessions and Cookies
	// Typically gorilla securecookie and sessions packages require
	// highly random secret keys that are not divulged to the public.
	//
	// In this example we use keys generated one time (if these keys ever become
	// compromised the gorilla libraries allow for key rotation, see gorilla docs)
	// The keys are 64-bytes as recommended for HMAC keys as per the gorilla docs.
	//
	// These values MUST be changed for any new project as these keys are already "compromised"
	// as they're in the public domain, if you do not change these your application will have a fairly
	// wide-opened security hole. You can generate your own with the code below, or using whatever method
	// you prefer:
	//
	//    func main() {
	//        fmt.Println(base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(64)))
	//    }
	//
	// We store them in base64 in the example to make it easy if we wanted to move them later to
	// a configuration environment var or file.
	cookieStoreKey, _ := base64.StdEncoding.DecodeString(`NpEPi8pEjKVjLGJ6kYCS+VTCzi6BUuDzU0wrwXyf5uDPArtlofn2AG6aTMiPmN3C909rsEWMNqJqhIVPGP3Exg==`)
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(`AbfYwmmt8UCwUuhd9qvfNA9UCuN1cVcKJN1ofbiky6xCyyBj20whe40rJa3Su0WOWLWcPpO1taqJdsEI/65+JA==`)
	cookieStore = securecookie.New(cookieStoreKey, nil)
	sessionStore = sessions.NewCookieStore(sessionStoreKey)

	// Initialize ab.
	SetupAuthboss()

	mux := mux.NewRouter()

	mux.PathPrefix("/auth").Handler(ab.NewRouter())

	// Set up our middleware chain
	stack := alice.New(logger, nosurfing, ab.ExpireMiddleware).Then(mux)

	// Start the server
	/*port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}*/

	port := "8888"
	//http.HandleFunc("/", _defaultHandler)
	http.ListenAndServe(":"+port, stack)
}

func layoutData(w http.ResponseWriter, r *http.Request) authboss.HTMLData {
	currentUserName := ""
	userInter, err := ab.CurrentUser(w, r)
	if userInter != nil && err == nil {
		currentUserName = userInter.(*User).Name
	}

	return authboss.HTMLData{
		"loggedin":               userInter != nil,
		"username":               "",
		authboss.FlashSuccessKey: ab.FlashSuccess(w, r),
		authboss.FlashErrorKey:   ab.FlashError(w, r),
		"current_user_name":      currentUserName,
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	data := layoutData(w, r).MergeKV("posts", blogs)
	mustRender(w, r, "index", data)
}

func listBuckets(w http.ResponseWriter, r *http.Request) {
	data := layoutData(w, r).MergeKV("posts", blogs)
	mustRender(w, r, "index", data)
}

func _defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello from Cisco Shipped testing!</h1>\n")
}

func mustRender(w http.ResponseWriter, r *http.Request, name string, data authboss.HTMLData) {
	data.MergeKV("csrf_token", nosurf.Token(r))
	err := templates.Render(w, name, data)
	if err == nil {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "Error occurred rendering template:", err)
}

func badRequest(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "Bad request:", err)

	return true
}
