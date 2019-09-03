package sessions

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRepository_Save(t *testing.T) {

	Convey("should save expected session", t, func() {
		factory := NewFactory()

		sess := factory.NewSession("test@test.com")
		repo := NewRepository()

		err := repo.Save(sess)

		So(err, ShouldBeNil)

		actual := repo.Store[sess.ID]
		So(actual, ShouldResemble, sess)
	})
}

func TestRepository_GetByID(t *testing.T) {
	Convey("getByID should return the expected session", t, func() {
		repo := NewRepository()

		sess, err := repo.GetByID("")
		So(sess, ShouldBeNil)
		So(err, ShouldResemble, IDEmptyError)
	})

	Convey("getByID should return the expected session", t, func() {
		repo := NewRepository()

		factory := NewFactory()
		sess := factory.NewSession("test@test.com")

		err := repo.Save(sess)
		So(err, ShouldBeNil)

		actual, err := repo.GetByID(sess.ID)
		So(err, ShouldBeNil)
		So(actual, ShouldResemble, sess)
	})
}
