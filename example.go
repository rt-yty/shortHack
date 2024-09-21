package main

import ("fmt"
		"net/http"
		"html/template")
		
func home_page(w http.ResponseWriter, r *http.Request){
	person := User{"Alex", 17, 5.0, 0.9}

	tmpl, _ := template.ParseFiles("templates/home_page.html")
	tmpl.Execute(w, person)

}

func contacts_page(w http.ResponseWriter, r *http.Request){
	person := User{"Alex", 17, 5.0, 0.9}
	person.setNewName("XXX")
	fmt.Fprintf(w, person.getAllInfo())
}

func handleRequest(){
    http.HandleFunc("/", home_page)
	http.HandleFunc("/contacts/", contacts_page)
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleRequest()
}

func (u User) getAllInfo() string {
	return fmt.Sprint("User name is ", u.Name, " and he is ", u.Age)
}

func (u *User) setNewName(newName string) {
	u.Name = newName
}

type User struct {
	Name string
	Age uint
	Money, Happines float64
}
