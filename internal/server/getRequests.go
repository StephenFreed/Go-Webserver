package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
    "encoding/json"
)


func (s *Server) handleGetUserByID(w http.ResponseWriter, r *http.Request) {
    fmt.Println("handleGetUserByID")
    // parse int from url path
    var urlPath []string = strings.Split(r.URL.Path, "/")[1:]
    idFromPath, err := strconv.Atoi(urlPath[1])
    if err != nil {
        fmt.Println(err)
        fmt.Printf("Error got %v of type %T \n", urlPath[1], urlPath[1])
    } else {
        // user := s.store.Get(idFromPath, r.URL.Query().Get("name"))
        type test struct {
            ID int
            Name string
        }
        user := &test {ID: idFromPath, Name: "somethign"}
        json.NewEncoder(w).Encode(user)
    }
}
