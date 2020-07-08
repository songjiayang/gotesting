package sqlmocktesting

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"

	post "gotesting/post"
)

func TestPostDaoList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	Convey("PostDao.Fetch", t, func() {
		dao := post.NewPostDao(db)

		Convey("should be successful", func() {
			rows := sqlmock.NewRows([]string{"id", "title", "body"}).
				AddRow(1, "post 1", "hello").
				AddRow(2, "post 2", "world")
			mock.ExpectQuery("^SELECT (.+) FROM posts$").
				WithArgs().WillReturnRows(rows)

			items, err := dao.List()
			So(items, ShouldHaveLength, 2)
			So(err, ShouldBeNil)

		})

		Convey("should be failed", func() {
			mock.ExpectQuery("^SELECT (.+) FROM posts$").
				WillReturnError(fmt.Errorf("list post error"))

			items, err := dao.List()
			So(items, ShouldBeNil)
			So(err.Error(), ShouldContainSubstring, "list post error")
		})
	})
}
