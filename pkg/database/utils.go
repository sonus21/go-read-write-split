package database

import "net/http"

// determineDbName this method is used to determine the actual database name from request method and path
// here more complex logic can be placed like for this path, path regex, header parameters etc
func determineDbName(method, path string, header http.Header) string {
	if method == http.MethodGet {
		return Secondary
	}
	// add any other logic
	return Primary
}

// dbNameFromRequest tries to identify the database name from request
func dbNameFromRequest(r *http.Request) string {
	method := r.Method
	path := r.URL.Path
	return determineDbName(method, path, r.Header)
}
