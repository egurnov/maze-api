package test_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("App", func() {
	var (
		c *client
	)

	BeforeEach(func() {
		c = &client{server: server}
	})

	Specify("Method not allowed and page not found", func() {
		resp := c.sendReq(http.MethodPost, "/user", `{"username": "alex", "password": "passw0rd"}`)
		Expect(resp.StatusCode).To(Equal(http.StatusCreated))

		resp = c.sendReq(http.MethodGet, "/nosuchpage", "")
		Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		Expect(ioutil.ReadAll(resp.Body)).To(MatchJSON(`{"message":"not found"}`))
	})

	Specify("New user registration", func() {
		resp := c.sendReq(http.MethodPost, "/user", `{"username": "alex", "password": "passw0rd"}`)
		Expect(resp.StatusCode).To(Equal(http.StatusCreated))

		resp = c.sendReq(http.MethodPost, "/user", `{"username": "alex", "password": "passw0rd"}`)
		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		Expect(ioutil.ReadAll(resp.Body)).To(ContainSubstring("username already used"))

		c.login("alex", "passw0rd")
	})

	Specify("Full flow", func() {
		resp := c.sendReq(http.MethodPost, "/user", `{"username": "alex", "password": "passw0rd"}`)
		Expect(resp.StatusCode).To(Equal(http.StatusCreated))

		c.login("alex", "passw0rd")

		var mazeId int64
		By("create maze")
		{
			resp = c.sendReq(http.MethodPost, "/maze", `{"gridSize": "8x10", "entrance": "A1", "walls": ["C1", "G1", "A2"]}`)
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))
			var maze IDResp
			Expect(json.NewDecoder(resp.Body).Decode(&maze)).To(Succeed())
			mazeId = maze.ID
		}

		By("get one")
		{
			resp = c.sendReq(http.MethodGet, fmt.Sprintf("/maze/%d", mazeId), "")
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			var maze Maze
			Expect(json.NewDecoder(resp.Body).Decode(&maze)).To(Succeed())
			Expect(maze).To(Equal(Maze{GridSize: "8x10", Entrance: "A1", Walls: []string{"C1", "G1", "A2"}}))
		}

		By("get all")
		{
			resp = c.sendReq(http.MethodGet, "/maze", "")
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			var mazes Mazes
			Expect(json.NewDecoder(resp.Body).Decode(&mazes)).To(Succeed())

			Expect(mazes.Mazes).To(HaveLen(1))
			Expect(mazes.Mazes[0]).To(Equal(Maze{GridSize: "8x10", Entrance: "A1", Walls: []string{"C1", "G1", "A2"}}))
		}

		By("other user")
		{
			c2 := &client{server: server}

			resp := c2.sendReq(http.MethodPost, "/user", `{"username": "alex2", "password": "passw0rd"}`)
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			c2.login("alex2", "passw0rd")

			resp = c2.sendReq(http.MethodGet, "/maze", "")
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			var mazes Mazes
			Expect(json.NewDecoder(resp.Body).Decode(&mazes)).To(Succeed())
			Expect(mazes.Mazes).To(BeEmpty())

			resp = c2.sendReq(http.MethodGet, fmt.Sprintf("/maze/%d", mazeId), "")
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			Expect(ioutil.ReadAll(resp.Body)).To(BeEquivalentTo(`{"message":"not found"}`))
		}
	})
})

type Maze struct {
	GridSize string
	Entrance string
	Walls    []string
}

type Mazes struct {
	Mazes []Maze
}

type IDResp struct {
	ID int64
}
