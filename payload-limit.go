package main

import (
	"log"
	"strconv"

	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	ep "github.com/wrossmorrow/envoy-extproc-sdk-go"
)

type payloadLimitRequestProcessor struct {
	opts         *ep.ProcessingOptions
	payloadLimit int64
}

func (s *payloadLimitRequestProcessor) GetName() string {
	return "payload-limit"
}

func (s *payloadLimitRequestProcessor) GetOptions() *ep.ProcessingOptions {
	return s.opts
}

const kContentLen = "content-length"

func (s *payloadLimitRequestProcessor) ProcessRequestHeaders(ctx *ep.RequestContext, headers map[string][]string, rawValues map[string][]byte) error {
	cancel := func(code int32) error {
		return ctx.CancelRequest(code, map[string]string{}, typev3.StatusCode_name[code])
	}
	raw, ok := rawValues[kContentLen]
	if !ok {
		return cancel(413)
	}

	size, _ := strconv.ParseInt(string(raw), 10, 64)
	if size > s.payloadLimit {
		log.Printf("the body size: %d exceeded the maximum size: %d\n", size, s.payloadLimit)
		return cancel(413)
	}

	return ctx.ContinueRequest()
}

func (s *payloadLimitRequestProcessor) ProcessRequestBody(ctx *ep.RequestContext, body []byte) error {
	return ctx.ContinueRequest()
}

func (s *payloadLimitRequestProcessor) ProcessRequestTrailers(ctx *ep.RequestContext, trailers map[string][]string, rawValues map[string][]byte) error {
	return ctx.ContinueRequest()
}

func (s *payloadLimitRequestProcessor) ProcessResponseHeaders(ctx *ep.RequestContext, headers map[string][]string, rawValues map[string][]byte) error {
	return ctx.ContinueRequest()
}

func (s *payloadLimitRequestProcessor) ProcessResponseBody(ctx *ep.RequestContext, body []byte) error {
	return ctx.ContinueRequest()
}

func (s *payloadLimitRequestProcessor) ProcessResponseTrailers(ctx *ep.RequestContext, trailers map[string][]string, rawValues map[string][]byte) error {
	return ctx.ContinueRequest()
}

const kPayloadLimit = "payload-limit"

func (s *payloadLimitRequestProcessor) Init(opts *ep.ProcessingOptions, nonFlagArgs []string) error {
	s.opts = opts
	s.payloadLimit = 16

	var (
		i            int
		err          error
		payloadLimit int64
	)

	nArgs := len(nonFlagArgs)
	for ; i < nArgs-1; i++ {
		if nonFlagArgs[i] == kPayloadLimit {
			break
		}
	}

	if i == nArgs {
		log.Printf("the argument: 'payload-limit' is missing, use the default.\n")
		return nil
	}

	payloadLimit, err = strconv.ParseInt(nonFlagArgs[i+1], 10, 64)
	if err != nil {
		log.Printf("parse the value for parameter: 'payload-limit' is failed: %v,use the default.\n", err.Error())
		return nil
	}

	if payloadLimit > 0 {
		s.payloadLimit = payloadLimit
		log.Printf("the payload limit is: %d.\n", s.payloadLimit)
	}

	return nil
}

func (s *payloadLimitRequestProcessor) Finish() {}
