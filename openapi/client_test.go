package openapi

import (
	"github.com/jack0829/letsgo/config"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/url"
	"os"
	"testing"
)

var (
	testClient *Client
)

type testTokenStorage struct {
	path string
}

func (s *testTokenStorage) Get() *AccessToken {
	b, err := os.ReadFile(s.path)
	if err != nil {
		return nil
	}
	var tk AccessToken
	if err = jsoniter.Unmarshal(b, &tk); err != nil {
		return nil
	}
	return &tk
}

func (s *testTokenStorage) Set(tk *AccessToken) error {
	b, _ := jsoniter.MarshalIndent(tk, "", "  ")
	return os.WriteFile(s.path, b, 0755)
}

func json(v any, pretty ...bool) (s string) {
	if len(pretty) < 1 {
		s, _ = jsoniter.MarshalToString(v)
		return
	}
	b, _ := jsoniter.MarshalIndent(v, "", "  ")
	s = string(b)
	return
}

func TestMain(m *testing.M) {

	testClient = NewClient(config.OpenAPI{
		Addr: "http://127.0.0.1:30084",
		Client: config.OpenAPIClient{
			ID:     "bKkjR84JDePBaOQe1Qdr5Voyln3wYMWN",
			Secret: "klqIIAtxerTZq7y7MyiXG3GaPhQlQgZ8",
		},
		Debug: true,
	}, WithAccessTokenStorager(&testTokenStorage{
		path: "./access.token",
	}))

	m.Run()
}

func TestClient_Token(t *testing.T) {

	tk, err := testClient.getAccessToken()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(json(tk))
}

func TestClient_Get(t *testing.T) {

	qs := url.Values{}
	qs.Set("state", "0,1")
	qs.Set("type", "NBP.ArticleList")

	resp, err := testClient.Get("/v1/nbp/admin/tasks", qs)
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	t.Logf("body: %s", string(b))
}

func TestClient_Post(t *testing.T) {

	resp, err := testClient.Post("/v1/nbp/admin/task", map[string]any{
		"id":    1,
		"state": 2,
		"meta": json(map[string]any{
			"sitemap": map[string]any{
				"gsdata": map[string]any{
					"accounts": []int{35, 36},
				},
				"download_url": "https://download-1.tar.gz",
				"data_total":   21,
			},
		}),
	})
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	t.Logf("body: %s", string(b))
}
