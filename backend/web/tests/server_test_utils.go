package tests

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
)

const (
	dummyGetOkUrl = "/dummygetok"
	dummyLogin    = "/dummyLogin"
	removeUser    = "/removeUser"
)

type TestServer struct {
	*httptest.Server
}

func NewTestServer(t *testing.T, h http.Handler) *TestServer {
	ts := httptest.NewServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &TestServer{ts}

}

func (ts *TestServer) ExecuteGet(t *testing.T, urlPath string) (*http.Response, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs, string(body)
}

func (ts *TestServer) ExecutePost(t *testing.T, urlPath string, values url.Values) (*http.Response, string) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, values)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs, string(body)
}

func DummyGetOkResponse(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
