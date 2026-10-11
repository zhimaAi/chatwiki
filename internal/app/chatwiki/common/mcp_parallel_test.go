package common

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"chatwiki/internal/app/chatwiki/define"
	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"sync/atomic"
)

func TestParallelMCPTransport(t *testing.T) {
	s := server.NewMCPServer("search-fixture", "1.0.0")
	s.AddTool(mcp.NewTool("web_search", mcp.WithString("objective", mcp.Required())), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("https://example.org: useful search excerpt"), nil
	})
	s.AddTool(mcp.NewTool("web_fetch", mcp.WithString("url", mcp.Required())), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("Example page content"), nil
	})
	handler := server.NewStreamableHTTPServer(s)
	var requests atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		if r.URL.Path != "/mcp" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		// mcp-go's session cleanup DELETE does not use its configured headers.
		if r.Method != http.MethodDelete && r.UserAgent() != "ChatWiki/0.1.0" {
			t.Errorf("unexpected User-Agent %q", r.UserAgent())
		}
		if r.Header.Get("Authorization") != "" {
			t.Error("anonymous request has Authorization")
		}
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handler.ServeHTTP(w, r)
	}))
	defer ts.Close()
	endpoint := ts.URL + "/mcp"
	clientType, err := DetectMCPTransportType(endpoint)
	if err != nil || clientType != define.McpClientTypeHttp {
		t.Fatalf("detect: %d %v", clientType, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	c, err := NewMcpClient(ctx, clientType, endpoint, `{"User-Agent":"ChatWiki/0.1.0"}`)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	tools, err := ListTools(ctx, c)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools) != 2 {
		t.Fatalf("tools: %v", tools)
	}
	for _, tool := range tools {
		args := map[string]any{"objective": "Find a useful page"}
		want := "useful search excerpt"
		if tool.Name == "web_fetch" {
			args = map[string]any{"url": "https://example.org"}
			want = "Example page content"
		}
		caller, err := NewMcpClient(ctx, clientType, endpoint, `{"User-Agent":"ChatWiki/0.1.0"}`)
		if err != nil {
			t.Fatal(err)
		}
		output, err := CallTool(ctx, caller, tool, args)
		caller.Close()
		if err != nil || !strings.Contains(output, want) {
			t.Fatalf("%s: %q %v", tool.Name, output, err)
		}
	}
	if requests.Load() < 7 {
		t.Fatalf("expected discovery and calls, saw %d requests", requests.Load())
	}
}

type parallelObservedTransport struct {
	base     http.RoundTripper
	requests atomic.Int32
}

func (r *parallelObservedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Session cleanup in mcp-go v0.42.0 omits configured headers.
	if req.Method == http.MethodDelete {
		return r.base.RoundTrip(req)
	}
	if req.URL.String() != "https://search.parallel.ai/mcp" || req.UserAgent() != "ChatWiki/0.1.0" || req.Header.Get("Authorization") != "" {
		return nil, fmt.Errorf("unexpected anonymous MCP request: %s User-Agent=%q", req.URL, req.UserAgent())
	}
	r.requests.Add(1)
	return r.base.RoundTrip(req)
}

func TestParallelMCPLive(t *testing.T) {
	if os.Getenv("CHATWIKI_PARALLEL_LIVE") != "1" {
		t.Skip("set CHATWIKI_PARALLEL_LIVE=1 for anonymous endpoint validation")
	}
	original := http.DefaultTransport
	observed := &parallelObservedTransport{base: original}
	http.DefaultTransport = observed
	defer func() { http.DefaultTransport = original }()
	endpoint := "https://search.parallel.ai/mcp"
	kind, err := DetectMCPTransportType(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	if kind != define.McpClientTypeHttp {
		t.Fatalf("transport=%d", kind)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	c, err := NewMcpClient(ctx, kind, endpoint, `{"User-Agent":"ChatWiki/0.1.0"}`)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	tools, err := ListTools(ctx, c)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	sessionID := uuid.NewString()
	for _, tool := range tools {
		var args map[string]any
		switch tool.Name {
		case "web_search":
			args = map[string]any{"objective": "Find the official Go programming language website", "search_queries": []any{"Go programming language official"}, "session_id": sessionID}
		case "web_fetch":
			args = map[string]any{"urls": []any{"https://go.dev"}, "session_id": sessionID}
		default:
			continue
		}
		caller, err := NewMcpClient(ctx, kind, endpoint, `{"User-Agent":"ChatWiki/0.1.0"}`)
		if err != nil {
			t.Fatal(err)
		}
		output, err := CallTool(ctx, caller, tool, args)
		caller.Close()
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(output, "go.dev") || len(output) < 100 {
			t.Fatalf("%s did not return useful Go website output: %s", tool.Name, output)
		}
		seen[tool.Name] = true
		t.Logf("%s returned %d characters: %.200s", tool.Name, len(output), output)
	}
	if !seen["web_search"] || !seen["web_fetch"] {
		t.Fatalf("missing tools: %v", seen)
	}
	t.Logf("Observed %d anonymous requests to exact /mcp endpoint with ChatWiki User-Agent", observed.requests.Load())
}
