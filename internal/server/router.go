package server

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
    "strings"
)


type route struct {
    method string
    regexp *regexp.Regexp
    handler http.HandlerFunc
}

func (s *Server) newRoute(method string, pattern string, handler http.HandlerFunc) *route {
    return &route{
        method: method,
        regexp: regexp.MustCompile("^" + pattern + "$"),
        handler: handler,
    }
}

type ctxKey struct {}

func (s *Server) router(w http.ResponseWriter, r *http.Request) {

    fmt.Println("==============================")

    // print http method and path that will be parsed
    fmt.Printf("Method: %v | Path: %v \n", r.Method, r.URL.Path[1:]) // 0 is empty

    //print query parameters
    queryParameters := r.URL.Query()
    for k,v := range queryParameters {
        fmt.Println("key: ", k, " => value: ", v)
    }

    var routes = []*route {
        s.newRoute("GET", "/user/([0-9]+)", s.handleGetUserByID),
    }

    var methodNotAllowed []string
    for _, route := range routes {
        fmt.Printf("route %T: %+v | ", *route, *route)
        pathMatch := (*route).regexp.FindStringSubmatch(r.URL.Path)
        fmt.Printf("Matched Path: %v \n", pathMatch)
        if len(pathMatch) > 0 {
            if r.Method != (*route).method {
                methodNotAllowed = append(methodNotAllowed, (*route).method)
            }
            var ctx context.Context = context.WithValue(r.Context(), ctxKey{}, pathMatch[1:])
            (*route).handler(w, r.WithContext(ctx))
            return
        }
    }
    if len(methodNotAllowed) > 0 {
        w.Header().Set("Not Allowed", strings.Join(methodNotAllowed, ", "))
        http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
    http.NotFound(w, r)
}
