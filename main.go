package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"time"
)

type CacheEntry struct {
	URL        string
	StatusCode int
	Header     http.Header
	Body       []byte
	StoredAt   time.Time
	TTL        time.Duration
}
type Cache struct {
	mu    sync.RWMutex
	store map[string]CacheEntry
}

var cache = Cache{
	store: make(map[string]CacheEntry),
}

func main() {
	/* args := os.Args */
	/* name := flag.String("name", "Gopher", "Your name")
	age := flag.Int("age", 25, "Your age")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()
	fmt.Println(*name, *age, *verbose)
	/* println(args[1]) */
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello there gang")
	})
	http.HandleFunc("/proxy", ProxyHandler(&cache))
	http.ListenAndServe(":8080", nil)
}
func ProxyHandler(cache *Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "missing url parameter", http.StatusBadRequest)
			return
		}

		cache.mu.RLock()
		entry, ok := cache.store[url]
		cache.mu.RUnlock()

		if ok && time.Since(entry.StoredAt) < entry.TTL {
			for k, v := range entry.Header {
				w.Header()[k] = v
			}
			w.Header().Set("X-Cache", "HIT")
			w.WriteHeader(entry.StatusCode)
			w.Write(entry.Body)
			return
		}

		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, "upstream error", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "read error", http.StatusInternalServerError)
			return
		}

		cache.mu.Lock()
		cache.store[url] = CacheEntry{
			URL:        url,
			StatusCode: resp.StatusCode,
			Header:     resp.Header.Clone(),
			Body:       body,
			StoredAt:   time.Now(),
			TTL:        time.Hour,
		}
		cache.mu.Unlock()

		for k, v := range resp.Header {
			if k == "Content-Length" {
				continue
			}
			w.Header()[k] = v
		}

		w.Header().Set("X-Cache", "MISS")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	}
}
