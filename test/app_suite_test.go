package test_test

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/egurnov/maze-api/maze-api/app"
	"github.com/egurnov/maze-api/maze-api/jwtservice"
	"github.com/egurnov/maze-api/maze-api/service"
	storepkg "github.com/egurnov/maze-api/maze-api/store"
)

const testDBConnString = "root:root@(localhost:3306)/maze_api_test"

func TestApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}

var (
	store   *storepkg.Store
	mazeAPI *app.App
	server  *httptest.Server

	_ = BeforeSuite(func() {
		var err error
		store, err = storepkg.NewMySQLStore(testDBConnString)
		Expect(err).ToNot(HaveOccurred())

		logger := logrus.New()
		logger.SetLevel(logrus.DebugLevel)
		logger.SetOutput(GinkgoWriter)

		jwtService := jwtservice.New([]byte("test"), time.Hour)

		gin.SetMode(gin.TestMode)
		engine := gin.New()
		engine.Use(gin.LoggerWithWriter(GinkgoWriter), gin.Recovery())

		mazeAPI = &app.App{
			Log:        logger,
			JWTService: jwtService,

			UserService: &service.UserService{Store: &storepkg.UserStore{Store: store}},
			MazeService: &service.MazeService{Store: &storepkg.MazeStore{Store: store}},
		}

		mazeAPI.SetRoutes(engine)

		server = httptest.NewServer(engine)
	})

	_ = AfterSuite(func() {
		server.Close()

		Expect(store.Close()).To(Succeed())
	})

	_ = BeforeEach(func() {
		Expect(store.Wipe()).To(Succeed())
	})

	_ = AfterEach(func() {
		Expect(store.Wipe()).To(Succeed())
	})
)

func printResponse(v io.Reader) string {
	b, err := io.ReadAll(v)
	Expect(err).ToNot(HaveOccurred())
	return fmt.Sprintf("Body: %s\n", b)
}
