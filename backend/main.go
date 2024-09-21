package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
	db    *gorm.DB
	err   error
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Email    string
	Password string
}

type Profile struct {
	gorm.Model
	UserID   uint
	FullName string
	Age      int
}

func initDB() {
	db, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	db.AutoMigrate(&User{}, &Profile{})
}

func main() {
	initDB()

	seedUsers()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/logout", logoutHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user User
		db.Where("username = ? AND password = ?", username, password).First(&user)

		if user.ID != 0 {
			session, _ := store.Get(r, "session-id")
			session.Values["authenticated"] = true
			session.Values["user_id"] = user.ID
			session.Save(r, w)

			http.Redirect(w, r, "/profile", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/login?error=1", http.StatusFound)
		return
	}
	tmpl, _ := template.ParseFiles("templates/login.html")
	tmpl.Execute(w, nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		user := User{Username: username, Email: email, Password: password}
		result := db.Create(&user)

		if result.Error != nil {
			http.Redirect(w, r, "/register?error=1", http.StatusFound)
			return
		}

		profile := Profile{UserID: user.ID, FullName: "", Age: 0}
		db.Create(&profile)

		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	tmpl, _ := template.ParseFiles("templates/register.html")
	tmpl.Execute(w, nil)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-id")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	userID := session.Values["user_id"].(uint)
	var user User
	db.First(&user, userID)

	var profile Profile
	db.Where("user_id = ?", userID).First(&profile)

	tmpl, _ := template.ParseFiles("templates/profile.html")
	data := struct {
		User    User
		Profile Profile
	}{
		User:    user,
		Profile: profile,
	}
	tmpl.Execute(w, data)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-id")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func seedUsers() {
	for i := 1; i <= 100; i++ {
		username := "user" + strconv.Itoa(i)
		email := "user" + strconv.Itoa(i) + "@example.com"
		password := "password"

		user := User{Username: username, Email: email, Password: password}
		db.Create(&user)

		profile := Profile{UserID: user.ID, FullName: "User " + strconv.Itoa(i), Age: 20 + i%10}
		db.Create(&profile)
	}
}
