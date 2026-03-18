package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("GET /a/b/{id}/c", basehandler)
	err := http.ListenAndServe(":1324", nil)
	if err != nil {
		panic(err)
	}
}

func basehandler(w http.ResponseWriter, r *http.Request) {
	if r.URL != nil {
		u := r.URL

		type field struct {
			label string
			value string
		}

		fields := []field{
			{"Scheme", u.Scheme},
			{"Opaque", u.Opaque},
			{"Host", u.Host},
			{"Path", u.Path},
			{"RawPath", u.RawPath},           // only differs from Path if encoded chars exist
			{"EscapedPath", u.EscapedPath()}, // use this over RawPath directly
			{"RawQuery", u.RawQuery},
			{"Fragment", u.Fragment},
			{"RawFragment", u.RawFragment},
		}

		if u.User != nil {
			fields = append(fields, field{"User", u.User.String()})
		}
		if u.ForceQuery {
			fields = append(fields, field{"ForceQuery", "true"})
		}
		if u.OmitHost {
			fields = append(fields, field{"OmitHost", "true"})
		}

		fmt.Println("--- URL Breakdown ---")
		for _, f := range fields {
			if f.value != "" {
				fmt.Printf("%-14s %s\n", f.label+":", f.value)
			}
		}
		fmt.Printf("%-14s %s\n", "Full URL:", u.String())
		fmt.Println("pathvalue", r.PathValue("id"))
		fmt.Println("queryname", u.Query())
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world \n"))
}
