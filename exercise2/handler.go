package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"net/http"
)

type Redirect struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestedPath := r.URL.Path

		if destination, ok := pathsToUrls[requestedPath]; ok {
			http.Redirect(w, r, destination, http.StatusFound)
			return
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	redirects, err := parseYAML(yml)

	if err != nil {
		fmt.Println("Error unmarshalling yaml")
		return nil, err
	}

	pathMap := buildMap(redirects)

	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yml []byte) ([]Redirect, error) {
	var redirects []Redirect

	err := yaml.Unmarshal([]byte(yml), &redirects)
	return redirects, err
}

func buildMap(redirects []Redirect) map[string]string {
	redirectMap := make(map[string]string)

	for _, redirect := range redirects {
		redirectMap[redirect.Path] = redirect.Url
	}

	return redirectMap
}
