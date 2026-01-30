package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxyCacheHitAndMiss(t *testing.T) {
	// Fake  server
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("hello world"))
	}))
	defer upstream.Close()

	cache := &Cache{
		store: make(map[string]CacheEntry),
	}

	handler := ProxyHandler(cache)
	req1 := httptest.NewRequest(
		http.MethodGet,
		"/proxy?url="+upstream.URL,
		nil,
	)
	rec1 := httptest.NewRecorder()
	//case 1
	handler.ServeHTTP(rec1, req1)

	if rec1.Header().Get("X-Cache") != "MISS" {
		t.Fatalf("expected MISS, got %s", rec1.Header().Get("X-Cache"))
	}
	req2 := httptest.NewRequest(
		http.MethodGet,
		"/proxy?url="+upstream.URL,
		nil,
	)
	//case 2
	rec2 := httptest.NewRecorder()

	handler.ServeHTTP(rec2, req2)

	if rec2.Header().Get("X-Cache") != "HIT" {
		t.Fatalf("expected HIT, got %s", rec2.Header().Get("X-Cache"))
	}
}

//to do
//
