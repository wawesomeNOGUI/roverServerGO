package main

import (
    "fmt"
    "log"
    "net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}



func formHandler(w http.ResponseWriter, r *http.Request) {
  if err := r.ParseForm(); err != nil {
      fmt.Fprintf(w, "ParseForm() err: %v", err)
      return
  }

    //Send User A File
    http.ServeFile(w, r, "./public/index.html")
  //fmt.Fprintf(w, "POST request successful")
  name := r.FormValue("name")
  address := r.FormValue("address")

  //Print to client
  //fmt.Fprintf(w, "Name = %s\n", name)
  //fmt.Fprintf(w, "Address = %s\n", address)

  //Print to server console
  //fmt.Printf("Name = %s\n", name)         //the %s\n is just to format the string writing?
  //fmt.Printf("Address = %s\n", address)   //look it up in the fmt documentation
  fmt.Printf("Name = " + name + "\n")
  fmt.Printf("Address = " + address + "\n")
}

func main() {
    // Register handler for default route
    http.HandleFunc("/hello", HelloHandler)
    http.HandleFunc("/formPOST", formHandler)

    //For serving clients files
    fileServer := http.FileServer(http.Dir("./public"))
    http.Handle("/", fileServer)



    // Start server listening
    fmt.Println("Listening on port 80...")
    err := http.ListenAndServe(":80", nil)
    if err != nil {
        //panic(err)
        log.Fatal(err)
    }

}
