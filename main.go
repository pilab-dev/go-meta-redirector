package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Domains map[string]DomainConfig `yaml:"domains"`
}

type DomainConfig struct {
	Fallback *FallbackConfig `yaml:"fallback"`
	Repos    []Repo          `yaml:"repos"`
}

type FallbackConfig struct {
	Pattern string `yaml:"pattern"`
	Target  string `yaml:"target"`
}

type Repo struct {
	Path       string `yaml:"path"`
	GitURL     string `yaml:"git_url"`
	PkgsiteURL string `yaml:"pkgsite_url,omitempty"`
}

var config Config

func loadConfig() error {
	data, err := ioutil.ReadFile("repos.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			config = Config{
				Domains: map[string]DomainConfig{
					"go.pilab.hu": {
						Fallback: &FallbackConfig{
							Pattern: "cloud/*",
							Target:  "https://github.com/pilab-dev/*",
						},
						Repos: []Repo{},
					},
				},
			}
			return nil
		}
		return err
	}

	return yaml.Unmarshal(data, &config)
}

func matchFallback(pattern, target, reqPath string) (string, bool) {
	parts := strings.Split(pattern, "*")
	if len(parts) != 2 {
		return "", false
	}
	prefix := parts[0]
	suffix := parts[1]

	if !strings.HasPrefix(reqPath, prefix) || !strings.HasSuffix(reqPath, suffix) {
		return "", false
	}

	wildcard := strings.TrimPrefix(reqPath, prefix)
	wildcard = strings.TrimSuffix(wildcard, suffix)
	if wildcard == "" {
		return "", false
	}

	gitURL := strings.Replace(target, "*", wildcard, 1)
	if !strings.HasSuffix(gitURL, ".git") {
		gitURL += ".git"
	}
	return gitURL, true
}

func lookup(host, reqPath string) (gitURL, pkgsiteURL string, ok bool) {
	host = strings.Split(host, ":")[0]
	domain, exists := config.Domains[host]
	if !exists {
		return "", "", false
	}

	for _, repo := range domain.Repos {
		if repo.Path == reqPath {
			return repo.GitURL, repo.PkgsiteURL, true
		}
	}

	if domain.Fallback != nil {
		gitURL, ok := matchFallback(domain.Fallback.Pattern, domain.Fallback.Target, reqPath)
		if ok {
			return gitURL, "", true
		}
	}

	return "", "", false
}

func handler(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	reqPath := strings.TrimPrefix(r.URL.Path, "/")
	gitURL, pkgsiteURL, ok := lookup(host, reqPath)

	if !ok {
		http.NotFound(w, r)
		return
	}

	if r.URL.Query().Get("go-get") == "1" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<html><head><meta name="go-import" content="%s/%s git %s"></head></html>`, host, reqPath, gitURL)
		return
	}

	if pkgsiteURL != "" {
		http.Redirect(w, r, pkgsiteURL, http.StatusFound)
		return
	}

	repoURL := strings.TrimSuffix(gitURL, ".git")
	http.Redirect(w, r, repoURL, http.StatusFound)
}

func main() {
	if err := loadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Loaded config with %d domains", len(config.Domains))

	addr := ":8080"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	http.HandleFunc("/", handler)
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
