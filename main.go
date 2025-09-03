package main

import "share-docs/pkg/db"

func main() {
	// r := routes.SetupRouter()

	// s := &http.Server{
	// 	Addr: ":8080",
	// 	Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	// 		if req.Method == "HEAD" {
	// 			req.Method = "GET"
	// 		}
	// 		r.ServeHTTP(w, req)
	// 	}),
	// 	ReadTimeout:    10 * time.Second,
	// 	WriteTimeout:   10 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }

	// s.ListenAndServe()

	db.Connect()

}
