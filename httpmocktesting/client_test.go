package httpmocktesting

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPostClientFetch(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	postFetchURL := "https://api.mybiz.com/posts"

	client := &PostClient{
		Client: &http.Client{
			Transport: httpmock.DefaultTransport,
		},
	}

	Convey("PostClient.Fetch", t, func() {
		Convey("without error", func() {
			httpmock.RegisterResponder("GET", postFetchURL,
				httpmock.NewStringResponder(200, `[{"id": 1, "title": "title", "body": "body"}]`))

			items, err := client.Fetch(postFetchURL, 1)
			So(len(items), ShouldEqual, 1)
			So(err, ShouldEqual, nil)
		})

		Convey("with error", func() {
			Convey("response data invalid", func() {
				httpmock.RegisterResponder("GET", postFetchURL,
					httpmock.NewStringResponder(200, `[{"id": "213"}]`))

				items, err := client.Fetch(postFetchURL, 1)
				So(items, ShouldBeEmpty)
				So(err, ShouldNotBeNil)
			})

			Convey("without error", func() {
				httpmock.RegisterResponder("GET", postFetchURL,
					httpmock.NewStringResponder(500, `some error`))

				items, err := client.Fetch(postFetchURL, 1)
				So(items, ShouldBeEmpty)
				So(err.Error(), ShouldContainSubstring, "some error")
			})
		})
	})
}
