package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sync"

	ep "github.com/wrossmorrow/envoy-extproc-sdk-go"
)

const digestKey = "x-extproc-request-digest"

type cache struct {
	mu    sync.Mutex
	cache map[string][]byte
}

func newCache() *cache {
	return &cache{
		cache: map[string][]byte{},
	}
}

func (c *cache) set(k string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[k] = data
}

func (c *cache) get(k string) (data []byte, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	data, ok = c.cache[k]
	return
}

type cachingRequestProcessor struct {
	opts *ep.ProcessingOptions
	c    *cache
}

func (s *cachingRequestProcessor) GetName() string {
	return "caching"
}

func (s *cachingRequestProcessor) GetOptions() *ep.ProcessingOptions {
	return s.opts
}

func calcAndSetPathDigest(ctx *ep.RequestContext) string {
	hasher := sha256.New()
	hasher.Write([]byte(ctx.FullPath))
	digest := hex.EncodeToString(hasher.Sum(nil))
	ctx.SetValue(digestKey, digest)
	return digest

}
func (s *cachingRequestProcessor) ProcessRequestHeaders(ctx *ep.RequestContext, headers ep.AllHeaders) error {
	if ctx.Method != http.MethodGet {
		return ctx.ContinueRequest()
	}

	digest := calcAndSetPathDigest(ctx)
	body, ok := s.c.get(digest)
	if !ok {
		return ctx.ContinueRequest()
	}
	return ctx.CancelRequest(200, map[string]ep.HeaderValue{}, string(body))

}

func (s *cachingRequestProcessor) ProcessRequestBody(ctx *ep.RequestContext, body []byte) error {
	return ctx.ContinueRequest()
}

func (s *cachingRequestProcessor) ProcessRequestTrailers(ctx *ep.RequestContext, trailers ep.AllHeaders) error {
	return ctx.ContinueRequest()
}

func (s *cachingRequestProcessor) ProcessResponseHeaders(ctx *ep.RequestContext, headers ep.AllHeaders) error {
	return ctx.ContinueRequest()
}

func (s *cachingRequestProcessor) ProcessResponseBody(ctx *ep.RequestContext, body []byte) error {
	digest, err := ctx.GetValue(digestKey)
	if err == nil {
		s.c.set(digest.(string), body)
	}
	return ctx.ContinueRequest()
}

func (s *cachingRequestProcessor) ProcessResponseTrailers(ctx *ep.RequestContext, trailers ep.AllHeaders) error {
	return ctx.ContinueRequest()
}

func (s *cachingRequestProcessor) Init(opts *ep.ProcessingOptions, nonFlagArgs []string) error {
	s.opts = opts
	s.c = newCache()
	return nil
}

func (s *cachingRequestProcessor) Finish() {}
