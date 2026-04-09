package gateway

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type service struct {
	registry   map[ServiceName]string
	httpClient *http.Client
}

func (e *ProxyError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Service, e.Message)
}

func NewGatewayService(configs []ServiceConfig) *service {
	registry := make(map[ServiceName]string, len(configs))

	for _, c := range configs {
		registry[c.Name] = c.BaseURL
	}

	return &service{
		registry: registry,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func DefaultConfigs() []ServiceConfig {
	return []ServiceConfig{
		{Name: ServiceAuth, BaseURL: "http://localhost:8081"},
		{Name: ServiceUser, BaseURL: "http://localhost:8082"},
	}
}

func (s *service) Proxy(pr ProxyRequest) (http.HandlerFunc, error) {
	targetBase, ok := s.registry[pr.TargetService]

	if !ok {
		return nil, &ProxyError{
			Code:    "service_not_configured",
			Message: "upstream service not registered",
			Service: pr.TargetService,
		}
	}

	target, err := url.Parse(targetBase)
	if err != nil {
		return nil, &ProxyError{
			Code:    "invalid_service_url",
			Message: "upstream URL is malformed",
			Service: pr.TargetService,
		}
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(req *httputil.ProxyRequest) {
			req.SetURL(target)
			if pr.StripPrefix != "" {
				req.Out.URL.Path = strings.TrimPrefix(req.Out.URL.Path, pr.StripPrefix)
				if !strings.HasPrefix(req.Out.URL.Path, "/") {
					req.Out.URL.Path = "/" + req.Out.URL.Path
				}
			}
			req.SetXForwarded()

			if userID := resolveUserID(req.In); userID != "" {
				req.Out.Header.Set("X-User-ID", userID)
			}

			if reqID := resolveRequestID(req.In); reqID != "" {
				req.Out.Header.Set("X-Request-ID", reqID)
			}

			log.Printf("[gateway:proxy] → %s %s%s", req.Out.Method, target.Host, req.Out.URL.Path)

		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("[gateway:service] upstream %q unreachable: %v", pr.TargetService, err)
			w.WriteHeader(http.StatusBadGateway)
		},

		ModifyResponse: func(res *http.Response) error {
			res.Header.Set("X-Served-By", string(pr.TargetService))
			res.Header.Del("X-Internal-Token")
			return nil
		},
		Transport: s.httpClient.Transport,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}, nil
}

func (s *service) HealthCheck() (*HealthResponse, error) {
	serviceStatus := make(map[string]string, len(s.registry))

	for name, baseUrl := range s.registry {
		resp, err := s.httpClient.Head(baseUrl + "/health")
		if err != nil || resp.StatusCode >= 500 {
			serviceStatus[string(name)] = "unreachable"
		} else {
			serviceStatus[string(name)] = "ok"
		}
	}

	overall := "ok"
	for _, status := range serviceStatus {
		if status != "ok" {
			return nil, errors.New("one or more upstream services are unavailable")
		}
	}
	return &HealthResponse{
		Status:   overall,
		Services: serviceStatus,
	}, nil
}

func resolveUserID(req *http.Request) string {
	if v := req.Context().Value(ContextKeyRequestID); v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func resolveRequestID(req *http.Request) string {
	if v := req.Context().Value(ContextKeyRequestID); v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}
