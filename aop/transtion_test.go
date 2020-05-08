package aop

import (
	"context"
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/xgxw/foundation-go/tests"
)

func TestSetTransactional(t *testing.T) {
	Convey("trans commit", t, func() {
		db, mock, teardown := tests.MockDB(t)
		defer teardown()

		mock.ExpectBegin()
		mock.ExpectCommit()
		_, transHandler, _, err := SetTransactional(context.Background(), db)
		defer transHandler(&err)
		So(err, ShouldBeNil)
		return
	})
	Convey("trans rollback", t, func() {
		db, mock, teardown := tests.MockDB(t)
		defer teardown()

		mock.ExpectBegin()
		mock.ExpectRollback()
		_, transHandler, _, err := SetTransactional(context.Background(), db)
		defer transHandler(&err)
		So(err, ShouldBeNil)
		tempF := func() (int, error) {
			return 1, errors.New("test error")
		}
		a, err := tempF()
		if a == 1 {
		}
		return
	})
}
