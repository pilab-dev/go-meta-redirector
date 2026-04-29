# Go Module Redirect Server
Tiny multi-domain Go module redirector for `go get` requests. Supports `go.pilab.hu`, `go.paalgyula.com`, `go.pira.hu` with per-domain config in `repos.yaml`.

## Usage
```bash
# Build
go build -o go-redirector .

# Run (default :8080, specify port as arg)
./go-redirector :8080
```

## How it works
- Responds to `?go-get=1` requests with `go-import` meta tags
- Browser requests redirect to pkgsite or GitHub repo
- Uses `Host` header to route to correct domain config
- Falls back to pattern matching if no explicit repo defined

## Config
Edit `repos.yaml` to add/remove repos or domains. Format:
```yaml
domains:
  go.example.com:
    fallback:
      pattern: "prefix/*"
      target: "https://github.com/org/*"
    repos:
      - path: prefix/repo
        git_url: https://github.com/org/repo.git
        pkgsite_url: https://pkg.go.dev/go.example.com/prefix/repo
```
