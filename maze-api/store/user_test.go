package store_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/egurnov/maze-api/maze-api/model"
	storepkg "github.com/egurnov/maze-api/maze-api/store"
)

const testDBConnString = "root:root@localhost:3306/maze_api"

var _ = Describe("UserStore", func() {
	var s *storepkg.UserStore

	BeforeEach(func() {
		store, err := storepkg.NewMySQLStore(testDBConnString)
		Expect(err).ToNot(HaveOccurred())

		s = &storepkg.UserStore{Store: store}
	})

	AfterEach(func() {
		Expect(s.Close()).To(Succeed())
		Expect(os.Remove(testDBConnString)).To(Succeed())
	})

	Specify("full flow", func() {

		user := &model.User{
			Username:     "me@example.com",
			PasswordHash: "passw0rd",
		}

		id, err := s.Create(user)
		Expect(err).ToNot(HaveOccurred())
		user.ID = id

		_, err = s.Create(&model.User{
			Username:     "me@example.com",
			PasswordHash: "passw0rd",
		})
		Expect(err).To(MatchError(model.ErrUsernameAlreadyUsed))

		found, err := s.GetByID(user.ID)
		Expect(err).ToNot(HaveOccurred())
		Expect(found).To(Equal(user))

		found, err = s.GetByUsername(user.Username)
		Expect(err).ToNot(HaveOccurred())
		Expect(found).To(Equal(user))
	})
})
