package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"

	"yatter-backend-go/app/config"
	"yatter-backend-go/app/dao"

	"github.com/stretchr/testify/assert"
)

func TestAccountRegistration(t *testing.T) {
	c := setup(t)
	defer c.Close()

	func() {
		resp, err := c.PostJSON("/v1/accounts", `{"username":"john"}`)
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, resp.StatusCode, http.StatusOK) {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var j map[string]interface{}
		if assert.NoError(t, json.Unmarshal(body, &j)) {
			assert.Equal(t, "john", j["username"])
		}
	}()

	func() {
		resp, err := c.Get("/v1/accounts/john")
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, resp.StatusCode, http.StatusOK) {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var j map[string]interface{}
		if assert.NoError(t, json.Unmarshal(body, &j)) {
			assert.Equal(t, "john", j["username"])
		}
	}()
}

func setup(t *testing.T) *C {
	db, err := dao.NewDB(config.MySQLConfig())
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := db.Exec("SET FOREIGN_KEY_CHECKS=0"); err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if _, err := db.Exec("SET FOREIGN_KEY_CHECKS=1"); err != nil {
			log.Fatalln(err)
		}
	}()

	for _, table := range []string{"account"} {
		if _, err := db.Exec("TRUNCATE TABLE " + table); err != nil {
			log.Fatalln(err)
		}
	}
	server := httptest.NewServer(NewRouter(
		dao.NewAccount(db), dao.NewStatus(db),
	))

	return &C{
		Server: server,
	}
}

type C struct {
	Server *httptest.Server
}

func (c *C) Close() {
	c.Server.Close()
}

func (c *C) PostJSON(apiPath string, payload string) (*http.Response, error) {
	return c.Server.Client().Post(c.asURL(apiPath), "application/json", bytes.NewReader([]byte(payload)))
}

func (c *C) Get(apiPath string) (*http.Response, error) {
	return c.Server.Client().Get(c.asURL(apiPath))
}

func (c *C) asURL(apiPath string) string {
	baseURL, _ := url.Parse(c.Server.URL)
	baseURL.Path = path.Join(baseURL.Path, apiPath)
	return baseURL.String()
}
