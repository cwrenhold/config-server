package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		// Ignore any error here, a .env file is not required
	}

	// Build a list of keys to ignore, these are the variables which are put in by default and won't be user-specific most likely
	ignoredKeys := []string{
		"HOSTNAME",
		"SHLVL",
		"HOME",
		"PATH",
		"PWD",
	}

	// Build up a map of the variables supplied
	variables := make(map[string]string)

	for _, pair := range os.Environ() {
		parts := strings.Split(pair, "=")
		key := parts[0]

		// Ignore any keys listed in the ignoredKeys slice
		skip := false
		for _, ignoredKey := range ignoredKeys {
			if key == ignoredKey {
				skip = true
			}
		}

		if !skip {
			values := strings.Join(parts[1:], "")
			variables[key] = values
		}
	}

	// Serve up an env format
	envHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")

		for key, values := range variables {
			w.Write([]byte(fmt.Sprintf("%s=%s", key, values) + "\n"))
		}
	}

	// Serve up a json format
	jsonHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		// Use a map[string]interface{} as a generic interface for the whole object
		data := make(map[string]interface{})

		for key, values := range variables {
			// Assume any double underscore indicates a nested object (like in languages like C#)
			parts := strings.Split(key, "__")

			var findParent func(parent map[string]interface{}, parts []string) map[string]interface{}

			// Recursive function to walk through the parts determining how deeply nested the key actually is
			findParent = func(parent map[string]interface{}, parts []string) map[string]interface{} {
				// If this is the last part, the current parent is the final parent for the key
				if len(parts) <= 1 {
					return parent
				}

				currentPart := parts[0]

				// If this part doesn't already exist, created it as a new nested map[string]interface{}
				val, exists := parent[currentPart]
				if !exists {
					parent[currentPart] = make(map[string]interface{})
					val = parent[currentPart]
				}

				currentMap := val.(map[string]interface{})

				return findParent(currentMap, parts[1:])
			}

			parent := findParent(data, parts)

			// The last part is the part of the key which is needed in the last parent
			// For example, NESTED__VALUE relates to the VALUE key within the NESTED object
			lastPart := parts[len(parts)-1]
			parent[lastPart] = values
		}

		// Use json.Marshal to convert this object to a json string to do some sanity checks
		jsonStr, err := json.Marshal(data)

		if err != nil {
			log.Fatal("Unable to process json data")
		}

		w.Write([]byte(jsonStr))
	}

	http.HandleFunc("/env", envHandler)
	http.HandleFunc("/json", jsonHandler)

	port := 80

	fmt.Println("Server running on port: ", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
