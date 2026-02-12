package proxy

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ReverseProxy handles proxying requests to backend services
type ReverseProxy struct {
	targetURL  string
	httpClient *http.Client
}

// NewReverseProxy creates a new reverse proxy
func NewReverseProxy(targetURL string) *ReverseProxy {
	return &ReverseProxy{
		targetURL: targetURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// Handler returns the gin handler for proxying requests
func (p *ReverseProxy) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		p.proxyRequest(c)
	}
}

func (p *ReverseProxy) proxyRequest(c *gin.Context) {
	// Build target URL
	targetURL, err := url.Parse(p.targetURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    nil,
			"message": "Internal server error",
		})
		return
	}

	// Preserve the original path
	targetURL.Path = c.Request.URL.Path
	targetURL.RawQuery = c.Request.URL.RawQuery

	// Create new request
	proxyReq, err := http.NewRequest(c.Request.Method, targetURL.String(), c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    nil,
			"message": "Failed to create proxy request",
		})
		return
	}

	// Copy headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Set forwarding headers
	proxyReq.Header.Set("X-Forwarded-For", c.ClientIP())
	proxyReq.Header.Set("X-Forwarded-Host", c.Request.Host)
	proxyReq.Header.Set("X-Forwarded-Proto", getScheme(c.Request))

	// Add request ID if present
	if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
		proxyReq.Header.Set("X-Request-ID", requestID)
	}

	// Execute request
	resp, err := p.httpClient.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"success": false,
			"data":    nil,
			"message": "Service unavailable",
		})
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Copy response body
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if scheme := r.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	if strings.HasPrefix(r.Host, "localhost") {
		return "http"
	}
	return "http"
}
