package main

import (
  "encoding/json"
  "log"
  "net/http"
  "fmt"
  "time"
  "html/template"
  "regexp"
  "database/sql"
  _ "github.com/lib/pq"
  "github.com/gorilla/mux"
  "github.com/speps/go-hashids"
)

//tinyURL Model
type TinyURL struct {
  ID string `json:"id"`
  LongURL string `json:"longURL"`
  ShortURL string `json:"shortURL"`
}

// Init db
var database *sql.DB

//Create a tinyURL and save into db
func CreateURL(tinyURL TinyURL) (TinyURL){
  hd := hashids.NewData()
  h, _ := hashids.NewWithData(hd)
  now := time.Now()
  tinyURL.ID, _ = h.Encode([]int{int(now.Unix())})
  tinyURL.ShortURL = "http://localhost:8080/d" + tinyURL.ID

  statement:= `INSERT INTO tinyURLs (id, shortURL, longURL) VALUES ($1, $2, $3)`
  _, err := database.Exec(statement, tinyURL.ID, tinyURL.ShortURL, tinyURL.LongURL)
  if err != nil {
    panic(err)
    return tinyURL
  }

  return tinyURL
}

func CheckHTTP(LongURL string) (string){
  matched, err := regexp.MatchString("http://*", LongURL)
  matched2, err := regexp.MatchString("https://*", LongURL)
  if err != nil {
    panic(err.Error())
    return LongURL
  }
  if !(matched) && !(matched2) && LongURL != ""{
    LongURL = "http://" + LongURL
  }
  return LongURL
}

// Create new router
func NewRouter() *mux.Router {
	r := mux.NewRouter()
  r.HandleFunc("/", Index).Methods("GET")
  r.HandleFunc("/d{id}", RedirectURL).Methods("GET")
	return r
}

// Redirect to LongURL associated with the ShortURL
func RedirectURL(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")

  // Extract ID from URL
  params := mux.Vars(r)
  var tinyURL TinyURL
  tinyURL.ID = params["id"]
  tinyURL.ShortURL = "http://localhost:8080/d" + tinyURL.ID

  // Query Database for the matching ShortURL row
  rows, err := database.Query(`SELECT shortURL, longURL FROM tinyURLS WHERE shortURL=$1`, tinyURL.ShortURL)
  if err != nil {
    panic(err.Error())
    return
  }
  defer rows.Close()

  // Scan for LongURL and redirect
  for rows.Next(){
    err = rows.Scan(&tinyURL.ShortURL, &tinyURL.LongURL)
    http.Redirect(w, r, tinyURL.LongURL, 301)
  }
  json.NewEncoder(w).Encode(tinyURL)
}

// Take in form data and presents tinyURL
func Index(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("views/index.html")
  var tempt TinyURL
  err := t.Execute(w, tempt)
  if err != nil {
    panic(err.Error())
    return
  }
  r.ParseForm()
  var tinyURL TinyURL
  tinyURL.LongURL=CheckHTTP(r.FormValue("LongURL"))
  tinyURL = CreateURL(tinyURL)
  if tinyURL.LongURL != "" {
    t, _ := template.ParseFiles("views/long.html")
    err := t.Execute(w, tinyURL)
    if err != nil {
      panic(err)
      return
    }
  }
}

// postgress local database constant
const (
  host = "localhost"
  port = 5432
  user = "postgres"
  password = "12345"
  dbname = "url_db"
)

func main() {
  // Initializing Router
  r := NewRouter()

  // Initialize postgres database connection//
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
  var err error
  database, err = sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer database.Close()

  // Force code to create a connection to the database once open
  err = database.Ping()
  if err != nil {
    panic(err)
  }

  // Create table in database
  statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS tinyURLS (id TEXT, longURL TEXT, shortURL TEXT)")
  statement.Exec()

  //Run Server
  log.Fatal(http.ListenAndServe(":8080", r))
}
