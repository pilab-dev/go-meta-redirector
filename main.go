package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed landing.html
var landingTemplate string

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
	Path        string `yaml:"path"`
	Description string `yaml:"description,omitempty"`
	Private     bool   `yaml:"private,omitempty"`
	GitURL      string `yaml:"git_url"`
	PkgsiteURL  string `yaml:"pkgsite_url,omitempty"`
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

// matchFallback maps reqPath onto target. The wildcard is the repository name
// only: anything after its first path segment (a Go major-version suffix such
// as /v3, or a package directory) stays in the import path but never becomes
// part of the git URL. It returns the git URL and the module root path
// (pattern prefix + repository), which is what go-import must advertise.
func matchFallback(pattern, target, reqPath string) (gitURL, rootPath string, ok bool) {
	parts := strings.Split(pattern, "*")
	if len(parts) != 2 {
		return "", "", false
	}
	prefix := parts[0]
	suffix := parts[1]

	if !strings.HasPrefix(reqPath, prefix) || !strings.HasSuffix(reqPath, suffix) {
		return "", "", false
	}

	wildcard := strings.TrimPrefix(reqPath, prefix)
	wildcard = strings.TrimSuffix(wildcard, suffix)
	repo, _, _ := strings.Cut(wildcard, "/")
	if repo == "" {
		return "", "", false
	}

	gitURL = strings.Replace(target, "*", repo, 1)
	if !strings.HasSuffix(gitURL, ".git") {
		gitURL += ".git"
	}
	return gitURL, prefix + repo, true
}

func lookup(host, reqPath string) (gitURL, pkgsiteURL, rootPath string, ok bool) {
	host = strings.Split(host, ":")[0]
	domain, exists := config.Domains[host]
	if !exists {
		return "", "", "", false
	}

	for _, repo := range domain.Repos {
		if repo.Path == reqPath {
			return repo.GitURL, repo.PkgsiteURL, repo.Path, true
		}
	}

	if domain.Fallback != nil {
		gitURL, rootPath, ok := matchFallback(domain.Fallback.Pattern, domain.Fallback.Target, reqPath)
		if ok {
			return gitURL, "", rootPath, true
		}
	}

	return "", "", "", false
}

type landingRepo struct {
	Domain      string
	FullPath    string
	ShortPath   string
	Description string
	Private     bool
	GitHubURL   string
	PkgsiteURL  string
}

type landingData struct {
	Repos []landingRepo
}

func renderLanding(w http.ResponseWriter, r *http.Request) {
	var repos []landingRepo

	for domain, cfg := range config.Domains {
		for _, repo := range cfg.Repos {
			fullPath := domain + "/" + repo.Path
			repos = append(repos, landingRepo{
				Domain:      domain,
				FullPath:    fullPath,
				ShortPath:   repo.Path,
				Description: repo.Description,
				Private:     repo.Private,
				GitHubURL:   strings.TrimSuffix(repo.GitURL, ".git"),
				PkgsiteURL:  repo.PkgsiteURL,
			})
		}
	}

	tmpl, err := template.New("landing").Parse(landingTemplate)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, landingData{Repos: repos}); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		renderLanding(w, r)
		return
	}

	host := r.Host
	reqPath := strings.TrimPrefix(r.URL.Path, "/")
	gitURL, pkgsiteURL, rootPath, ok := lookup(host, reqPath)

	if !ok {
		http.NotFound(w, r)
		return
	}

	if r.URL.Query().Get("go-get") == "1" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<html><head><meta name="go-import" content="%s/%s git %s"></head></html>`, host, rootPath, gitURL)
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
