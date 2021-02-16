package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopher-jasmine/suite"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
type Api struct {
	suites []suite.Suite
}

func NewApi(suites []suite.Suite) *Api {
	return &Api{
		suites: suites,
	}
}
func (api *Api) ListenAndServe(port string) {
	r := mux.NewRouter()
	endpoints := make([]string, 0)
	for _, s := range api.suites {
		name := strings.ToLower(s.GetName())
		name = strings.Join(strings.Split(name, " "), "-")
		endpoints = append(endpoints, name)
		r.HandleFunc(fmt.Sprintf("/%s", name), createSuiteHandler(s))
	}
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		j, _ := json.Marshal(endpoints)
		fmt.Fprintf(w, string(j))
	})
	fmt.Printf("starting server on port%s\n", port)
	http.ListenAndServe(port, r)
}

func createSuiteHandler(suite suite.Suite) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s\n", suite.GetName())
		result := suite.Run()
		j, err := json.Marshal(result)
		if err != nil {
			errorResponse, _ := json.Marshal(ErrorResponse{
				Status:  "500",
				Message: fmt.Sprintf("Failed to get results with error: %s", err.Error()),
			})
			fmt.Fprintf(w, string(errorResponse))
		}
		fmt.Fprintf(w, string(j))
	}
}