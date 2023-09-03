package server

import (
	"fmt"
	"net/http"
	"regexp"
    "strings"
	"context"
    "runtime"
    "reflect"
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

func getField(r *http.Request, index int) string {
    var fields []string = r.Context().Value(ctxKey{}).([]string)
    return fields[index]
}

func (s *Server) router(w http.ResponseWriter, r *http.Request) {

    fmt.Println("==============================")

    // print http method and path that will be parsed
    fmt.Printf("Method: %v | Path: %v \n", r.Method, r.URL.Path)

    //print query parameters | example /some/path?name=bob
    queryParameters := r.URL.Query()
    for k,v := range queryParameters {
        fmt.Println("key: ", k, " => value: ", v)
    }

    var routes = []*route {
        s.newRoute("GET", "/user/([0-9]+)", s.handleGetUserByID),
    }

    var methodsAllowed []string
    for _, route := range routes {
        fmt.Printf("Route Check: {method:%v regexp:%v handler:%v}", 
            route.method, route.regexp, runtime.FuncForPC(reflect.ValueOf(routes[0].handler).Pointer()).Name())
        // pathMatch[0] is full path | pathMatch[1:] is array of all regex matches in path
        var pathMatch []string = (*route).regexp.FindStringSubmatch(r.URL.Path)
        if len(pathMatch) > 0 {
            // prints full path matched, then regex parameters that matched in array
            fmt.Printf("\nMatched Path %v on %v \n", pathMatch[0], pathMatch[1:])
            // if path matches, but method does not
            if r.Method != (*route).method {
                methodsAllowed = append(methodsAllowed, (*route).method)
                continue
            }
            // create context to pass matches to handler
            var ctx context.Context = context.WithValue(r.Context(), ctxKey{}, pathMatch[1:])
            (*route).handler(w, r.WithContext(ctx))
            return
        }
    }
    // if path matched with incorrect method, and no correct method for that path was found
    if len(methodsAllowed) > 0 {
        w.Header().Set("Allow", strings.Join(methodsAllowed, ", "))
        http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
    // path not found
    http.NotFound(w, r)
}
