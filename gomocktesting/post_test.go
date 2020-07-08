package gomocktesting

import (
	"errors"
	"net/http/httptest"
	"testing"

	post "gotesting/post"

	gomock "github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPostIndexWithGoMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Convey("PostController.Index", t, func() {
		Convey("should be 200", func() {
			posts := []*post.PostModel{
				{1, "title", "body"},
				{2, "title2", "body2"},
			}

			m := NewMockPostService(ctrl)
			m.
				EXPECT().
				List().
				Return(posts, nil)

			handler := post.PostController{
				PostService: m,
			}

			req := httptest.NewRequest("GET", "http://example.com/foo", nil)
			w := httptest.NewRecorder()

			handler.Index(w, req)

			So(w.Result().StatusCode, ShouldEqual, 200)
		})

		Convey("should be 500", func() {
			m := NewMockPostService(ctrl)
			m.
				EXPECT().
				List().
				Return(nil, errors.New("list post with error"))

			handler := post.PostController{
				PostService: m,
			}

			req := httptest.NewRequest("GET", "http://example.com/foo", nil)
			w := httptest.NewRecorder()
			handler.Index(w, req)
			So(w.Result().StatusCode, ShouldEqual, 500)
		})
	})
}
