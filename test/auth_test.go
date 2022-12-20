package test_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
	var (
		c *client
	)

	BeforeEach(func() {
		c = &client{server: server}
	})

	Specify("No token", func() {
		testCases := []struct {
			method string
			path   string
		}{}
		for _, tc := range testCases {
			By(tc.method + " " + tc.path)
			req, err := http.NewRequest(tc.method, server.URL+tc.path, nil)
			Expect(err).ToNot(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
			Expect(ioutil.ReadAll(resp.Body)).To(BeEmpty())
		}
	})

	Specify("User token", func() {
		By("Registration")
		resp := c.sendReq(http.MethodPost, "/user", `{"username": "alex", "password": "passw0rd"}`)
		Expect(resp.StatusCode).To(Equal(http.StatusCreated))

		By("Bad login")
		resp = c.sendReq(http.MethodPost, "/login", `{"username": "alex", "password": "wrongPassw0rd"}`)
		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))

		By("Good login")
		c.login("alex", "passw0rd")

		testCases := []struct {
			method    string
			path      string
			expStatus int
		}{}
		for _, tc := range testCases {
			By(tc.method + " " + tc.path)
			resp := c.sendReq(tc.method, tc.path, "")
			Expect(resp.StatusCode).To(Equal(tc.expStatus), printResponse(resp.Body))
		}
	})
})

type loginResponse struct {
	Token string
}

type client struct {
	server *httptest.Server
	token  string
}

func (c *client) login(email, password string) {
	resp, err := http.Post(c.server.URL+"/login", "application/json", strings.NewReader(
		fmt.Sprintf(`{"username": "%s", "password": "%s"}`, email, password),
	))
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	var respObj loginResponse
	Expect(json.NewDecoder(resp.Body).Decode(&respObj)).To(Succeed())
	Expect(respObj.Token).ToNot(BeEmpty())
	c.token = respObj.Token
}

func (c *client) sendReq(method, path string, body string) *http.Response {
	req, err := http.NewRequest(method, c.server.URL+path, strings.NewReader(body))
	Expect(err).ToNot(HaveOccurred())
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	Expect(err).ToNot(HaveOccurred())
	return resp
}
