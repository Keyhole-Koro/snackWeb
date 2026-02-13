package adapter

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// LambdaAdapter converts Lambda Proxy Request -> http.Request -> http.Handler -> Lambda Proxy Response
type LambdaAdapter struct {
	handler http.Handler
}

func New(h http.Handler) *LambdaAdapter {
	return &LambdaAdapter{handler: h}
}

func (a *LambdaAdapter) Proxy(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 1. Create http.Request
	body := req.Body
	if req.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(body)
		if err == nil {
			body = string(decoded)
		}
	}

	url := req.Path
	if len(req.QueryStringParameters) > 0 {
		url += "?"
		params := []string{}
		for k, v := range req.QueryStringParameters {
			params = append(params, k+"="+v)
		}
		url += strings.Join(params, "&")
	}

	httpReq, err := http.NewRequest(req.HTTPMethod, url, strings.NewReader(body))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, nil
	}

	// Add Headers
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// 2. Record Response
	rec := httptest.NewRecorder()
	a.handler.ServeHTTP(rec, httpReq)

	// 3. Convert to Lambda Response
	resp := rec.Result()
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	
	proxyResp := events.APIGatewayProxyResponse{
		StatusCode: resp.StatusCode,
		Headers:    make(map[string]string),
		Body:       string(respBody),
	}

	for k, v := range resp.Header {
		proxyResp.Headers[k] = strings.Join(v, ",")
	}
	
	// Add CORS headers by default for Lambda execution if not handled by middleware
	// But our middleware does it, so it should be in Headers map already.

	return proxyResp, nil
}
