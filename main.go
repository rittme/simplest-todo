package main


import (
  "encoding/json"
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "strconv"
  "time"
  "github.com/gorilla/mux"
  "rittme.com/rittme/simple-list/model" 
)

// Serve the favicon file
func serveFavicon(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "static/favicon.ico")
}

// Return the todo list
func returnList(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, model.FilePath)
}

// Return random entry from the list
func returnRandom(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  includeDone := vars["includeDone"]

  entry, err := model.GetRandom(includeDone)
  if err != nil {
    fmt.Fprintf(w, "Error getting random entry: %v", err)
    return
  }
  if json.NewEncoder(w).Encode(entry) != nil {
    fmt.Fprintf(w, "Error parsing entries: %v", err)
    return
  }
}

// Add a new entry to the list
func addNewEntry(w http.ResponseWriter, r *http.Request) {
  label, _ := ioutil.ReadAll(r.Body)
  
  newEntries, err := model.CreateNew(string(label))
  if err != nil {
    fmt.Fprintf(w, "Error adding new entry: %v", err)
    return
  }
  if json.NewEncoder(w).Encode(newEntries) != nil {
    fmt.Fprintf(w, "Error parsing entries: %v", err)
    return
  }
}

// Toggle if an entry is done or not
func toggleEntryDone(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id, err := strconv.Atoi(vars["id"])

  if err != nil {
      fmt.Fprintf(w, "Index is not an integer: %d", id)
  }

  newEntries, err := model.ToggleDone(id)
  if err != nil {
    fmt.Fprintf(w, "Error changing entry status: %v", err)
    return
  }
  if json.NewEncoder(w).Encode(newEntries) != nil {
    fmt.Fprintf(w, "Error parsing entries: %v", err)
    return
  }
}

// Delete an entry from the list
func deleteEntry(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id, err := strconv.Atoi(vars["id"])

  if err != nil {
      fmt.Fprintf(w, "Index is not an integer: %d", id)
  }

  newEntries, err := model.DeleteEntry(id)
  if err != nil {
    fmt.Fprintf(w, "Error deleting entry: %v", err)
    return
  }
  if json.NewEncoder(w).Encode(newEntries) != nil {
    fmt.Fprintf(w, "Error parsing entries: %v", err)
    return
  }
}

// Serves the static home page
func homePage(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "static/index.html")
}

// Starts the server
func startServer(port, certFile, keyFile string) {
  // New router
  router := mux.NewRouter().StrictSlash(true)

  // Define routes
  router.HandleFunc("/", homePage).Methods("GET")
  router.HandleFunc("/entry", returnList).Methods("GET")
  router.HandleFunc("/entry", addNewEntry).Methods("POST")
  router.HandleFunc("/entry/{id}", deleteEntry).Methods("DELETE")
  router.HandleFunc("/entry/{id}", toggleEntryDone).Methods("PUT")
  router.HandleFunc("/random", returnRandom).Methods("GET")
  router.HandleFunc("/random/{includeDone}", returnRandom).Methods("GET")

  // Define route for favicon and for static files
  router.HandleFunc("/favicon.ico", serveFavicon).Methods("GET")
  router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

  // Define server parameters
  fmt.Printf("Starting server at port %s\n", port)
  srv := &http.Server {
    Handler: router,
    Addr:    "0.0.0.0:" + port,
    WriteTimeout: 15 * time.Second,
    ReadTimeout:  15 * time.Second,
	}

  // Start the server as HTTPS if cert and key were defined, else starts HTTP
  if certFile != "" && keyFile != "" {
    log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
  } else {
    log.Fatal(srv.ListenAndServe())
  }
}

// Launches the main program
func main() {
  // Read flags from command line
  portPtr := flag.String("port", "8080", "the server port")
  certPathPtr := flag.String("cert", "", "the path to the certificate file")
  keyFilePathPtr := flag.String("pk", "", "the path to the private key file")
  flag.Parse()

  // Start the server
  startServer(*portPtr, *certPathPtr, *keyFilePathPtr)
}