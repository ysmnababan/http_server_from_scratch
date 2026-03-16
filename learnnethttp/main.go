package main

import "net/http"

func main() {
	http.HandleFunc("PUT /test", basehandler)
	err := http.ListenAndServe(":1324", nil)
	if err != nil {
		panic(err)
	}
}

func basehandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world \n"))
}
