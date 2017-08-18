package handlers

import "net/http"

// Handler handles the '/' route
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func Worker(w http.ResponseWriter, r *http.Request) {
	// ctx := appengine.NewContext(r)
	// name := r.FormValue("name")
	// key := datastore.NewKey(ctx, "Counter", name, 0, nil)
	// if err := datastore.Get(ctx, key, &counter); err == datastore.ErrNoSuchEntity {
	// } else if err != nil {
	// 	log.Printf("%v", err)
	// 	return
	// }
	// counter.Count++
	// if _, err := datastore.Put(ctx, key, &counter); err != nil {
	// 	log.Printf("%v", err)
	// }
}
