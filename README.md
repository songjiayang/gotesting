# Testing In Go

The examples about my slide share - [Testing In Go](https://www.slideshare.net/songjiayang/testing-in-go-236709707).


## Examples

### Basic

A exmaple about sum function:

```golang
import "testing"

func TestSum(t *testing.T) {
	sum := Sum(1, 1)

	if sum != 2 {
		t.Errorf("Sum(1, 1) = %d; want 2", sum)
	}
}
```

Summary:

- Tests are written on files ending with "_test.go" 
- Test function starts with Test* and has only the parameter *testing.T
- Use go test command to run your tests

Usage of `go test`:

```bash
go test        // testing the local package
go test some/pkg       // testing a specific package
go test some/pkg/…    // testing a specific package in recursive
go test -v some/pkg –run ^TestSum$    // runs a specified tests
go test -cover   // code coverage
go test -count=1   // testing without cache
```

Reference:

- https://golang.org/pkg/testing
- https://golang.org/pkg/cmd/go/internal/test

### Table testing

A example about IPV4 regexp:

```golang
package tabletesting

import "regexp"

var (
	ipRegex = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
)

func IsIPV4(ip string) bool {
	return ipRegex.MatchString(ip)
}
```

Without table testing:

```golang
func TestIsIPV4WithoutTable(t *testing.T) {
	if IsIPV4("") {
		t.Errorf("IsIPV4(%s) should be false", "")
	}

	if IsIPV4("192.168.0") {
		t.Errorf("IsIPV4(%s) should be false", "192.168.0")
	}

	if IsIPV4("192.168.x.1") {
		t.Errorf("IsIPV4(%s) should be false", "192.168.x.1")
	}

	if IsIPV4("192.168.0.1.1") {
		t.Errorf("IsIPV4(%s) should be false", "192.168.0.1.1")
	}

	if !IsIPV4("127.0.0.1") {
		t.Errorf("IsIPV4(%s) should be true", "127.0.0.1")
	}

	if !IsIPV4("192.168.0.1") {
		t.Errorf("IsIPV4(%s) should be true", "192.168.0.1")
	}

	if !IsIPV4("255.255.255.255") {
		t.Errorf("IsIPV4(%s) should be true", "255.255.255.255")
	}

	if !IsIPV4("120.52.148.118") {
		t.Errorf("IsIPV4(%s) should be true", "120.52.148.118")
	}
}
```

With table testing:

```golang
func TestIsIPV4WithTable(t *testing.T) {
	testCases := []struct {
		IP    string
		valid bool
	}{
		{"", false},
		{"192.168.0", false},
		{"192.168.x.1", false},
		{"192.168.0.1.1", false},
		{"127.0.0.1", true},
		{"192.168.0.1", true},
		{"255.255.255.255", true},
		{"120.52.148.118", true},
	}

	for _, tc := range testCases {
		t.Run(tc.IP, func(t *testing.T) {
			if IsIPV4(tc.IP) != tc.valid {
				t.Errorf("IsIPV4(%s) should be %v", tc.IP, tc.valid)
			}
		})
	}
}
```

Summary:

- Using anonymous structs to represent test cases
- Using Subtests with t.Run

Reference:

- https://golang.org/pkg/testing/#hdr-Subtests_and_Sub_benchmarks

### Testing HTTP

A example about user login:

```golang
type LoginForm struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HandleResponse(w, 500, "read post body failed")
		return
	}
	defer r.Body.Close()

	var input LoginForm
	if err = json.Unmarshal(data, &input); err != nil {
		HandleResponse(w, 400, "input invalid format")
		return
	}

	if input.Code != "a@example.com" || input.Password != "password" {
		HandleResponse(w, 400, "invalid code or password")
		return
	}

	HandleResponse(w, 200, "ok")
}

func HandleResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}
```

The testing code:

```golang
type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("mock body error")
}

func TestLoginHandler(t *testing.T) {

	testCases := []struct {
		Name string
		Code int
		Body interface{}
	}{
		{"ok", 200, `{"code":"a@example.com", "password":"password"}`},
		{"read body error", 500, new(errorReader)},
		{"invalid format", 400, `{"code":1, "password":"password"}`},
		{"invalid code", 400, `{"code":"a@example.com1", "password":"password"}`},
		{"invalid password", 400, `{"code":"a@example.com", "password":"password1"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			var body io.Reader
			if stringBody, ok := tc.Body.(string); ok {
				body = strings.NewReader(stringBody)
			} else {
				body = tc.Body.(io.Reader)
			}

			req := httptest.NewRequest("POST", "http://example.com/foo", body)
			w := httptest.NewRecorder()

			LoginHandler(w, req)

			resp := w.Result()
			if resp.StatusCode != tc.Code {
				t.Errorf("response code is invalid, expect=%d but got=%d",
					tc.Code, resp.StatusCode)
			}
		})
	}
}
```

Summary:

- Use `httptest.NewRecorder` don't need http listen and speed testing
- Use errorReader to improve coverage

Reference:

- https://golang.org/pkg/net/http/httptest

### Testify 

A example about refactoring our `IPV4` [unit test](https://github.com/songjiayang/gotesting/blob/master/tabletesting/ip_test.go#L7) with testify:

```golang
func TestIsIPV4WithTestify(t *testing.T) {
	assertion := assert.New(t)

	assertion.False(IsIPV4(""))
	assertion.False(IsIPV4("192.168.0"))
	assertion.False(IsIPV4("192.168.x.1"))
	assertion.False(IsIPV4("192.168.0.1.1"))
	assertion.True(IsIPV4("127.0.0.1"))
	assertion.True(IsIPV4("192.168.0.1"))
	assertion.True(IsIPV4("255.255.255.255"))
	assertion.True(IsIPV4("120.52.148.118"))
}
```

Summary:

- Testify run within go test
- Assertions, mostly shortcuts 
- Testify can do mocking

Reference: 

- https://github.com/stretchr/testify  

### GinkGo

Refactoring our `IPV4` [unit test](https://github.com/songjiayang/gotesting/blob/master/tabletesting/ip_test.go#L7) with GinkGo:

```golang
import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "gotesting/tabletesting"
)

var _ = Describe("Ip", func() {
	Describe("IsIPV4()", func() {
		// fore content level prepare
		BeforeEach(func() {
			// prepare data before every case
		})

		AfterEach(func() {
			// clear data after every case
		})

		Context("should be invalid", func() {
			It("empty string", func() {
				Expect(IsIPV4("")).To(Equal(false))
			})

			It("with less length", func() {
				Expect(IsIPV4("192.0.1")).To(Equal(false))
			})

			It("with more length", func() {
				Expect(IsIPV4("192.168.1.0.1")).To(Equal(false))
			})

			It("with invalid character", func() {
				Expect(IsIPV4("192.168.x.1")).To(Equal(false))
			})
		})

		Context("should be valid", func() {
			It("loopback address", func() {
				Expect(IsIPV4("127.0.0.1")).To(Equal(true))
			})

			It("extranet address", func() {
				Expect(IsIPV4("120.52.148.118")).To(Equal(true))
			})
		})
	})
})

func TestGinkgotesting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ginkgotesting Suite")
}
```

Run with command `ginkgo` or `go test`, you also can watch the changes with `ginkgo watch`.

Summary:

- BDD 
- Work with go test, but has it's own structure
- Uses a custom lib for assertion (Gomega) 
- Rerun testes on change

Reference:

- https://github.com/onsi/ginkgo 
- https://en.wikipedia.org/wiki/Behavior-driven_development 
- https://github.com/onsi/gomega

### GoConvey

Refactoring our `IPV4` [unit test](https://github.com/songjiayang/gotesting/blob/master/tabletesting/ip_test.go#L7) with GoConvey:

```golang

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "gotesting/tabletesting"
)

func TestIsIPV4WithGoconvey(t *testing.T) {
	Convey("ip.IsIPV4()", t, func() {
		Convey("should be invalid", func() {
			Convey("empty string", func() {
				So(IsIPV4(""), ShouldEqual, false)
			})

			Convey("with less length", func() {
				So(IsIPV4("192.0.1"), ShouldEqual, false)
			})

			Convey("with more length", func() {
				So(IsIPV4("192.168.1.0.1"), ShouldEqual, false)
			})

			Convey("with invalid character", func() {
				So(IsIPV4("192.168.x.1"), ShouldEqual, false)
			})
		})

		Convey("should be valid", func() {
			Convey("loopback address", func() {
				So(IsIPV4("127.0.0.1"), ShouldEqual, true)
			})

			Convey("extranet address", func() {
				So(IsIPV4("120.52.148.118"), ShouldEqual, true)
			})
		})
	})
}
```

Use `goconvey` || `go test` to run testing.

Summary:

- BDD alike
- Work with go test
- Pretty browser interface
- Reload tests on file changes
- Custom DSL 

Reference:

- https://github.com/smartystreets/goconvey 

### GoMock

A example about PostService:

```golang
type PostController struct {
	PostService PostService
}

func (c *PostController) Index(w http.ResponseWriter, r *http.Request) {
	posts, err := c.PostService.List()
	if err != nil {
		HandleResponse(w, 500, "list posts with error")
		return
	}

	data, _ := json.Marshal(posts)
	w.WriteHeader(200)
	w.Write(data)
}

type PostService interface {
	List() ([]*PostModel, error)
	Find(int64) (*PostModel, error)
	Create(PostModel) error
	Update(PostModel) error
	Destroy(int64) error
}

type PostModel struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
```

Using `mockgen` to generate `NewMockPostService` with command like `mockgen -source=../post/post.go -destination=mock_post_test.go -package=gomocktesting`.

```golang
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

```

Summary:

- Official library
- An interface mock
- Mock and stub support
- Generate code with command `mockgen`

Reference:

- https://github.com/golang/mock

### HTTPMock

A example about post client to fetch items and unmarshal:

```golang
type PostClient struct {
	*http.Client
}

func (c *PostClient) Fetch(url string, page int) ([]*post.PostModel, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return nil, errors.New(string(data))
	}

	var items []*post.PostModel
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	return items, nil
}
```

The testing code:

```golang
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
```

Summary: 

- Mock http request
- Custom any response
- Base URL regular matching  

Reference:

- https://github.com/jarcoal/httpmock

### SQLMock

A example about PostDAO implement:

```golang
type PostDao struct {
	db *sql.DB
}

func NewPostDao(db *sql.DB) *PostDao {
	return &PostDao{
		db: db,
	}
}

func (dao *PostDao) List() ([]*PostModel, error) {
	rows, err := dao.db.Query("SELECT id, title, body FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*PostModel
	for rows.Next() {
		p := &PostModel{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Body); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return posts, nil
}

func (*PostDao) Find(int64) (*PostModel, error) {
	return nil, nil
}

func (*PostDao) Create(PostModel) error {
	return nil
}

func (*PostDao) Update(PostModel) error {
	return nil
}

func (*PostDao) Destroy(int64) error {
	return nil
}

type PostModel struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

```

The testing code:

```golang
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
```

Summary:

- Mock for database/sql
- Base regular matching
- Support query, exec, transaction 

Reference:

- https://github.com/DATA-DOG/go-sqlmock

### Testing With Docker

You can use official library for your testing, but you also can custom it for complex situation. 

A example about testing with MongoDB:

```bash
FROM ubuntu:16.04
RUN apt-get update && apt-get install -y libssl1.0.0 libssl-dev gcc

RUN mkdir -p /data/db /opt/go/ /opt/gopath
COPY mongodb/bin/* /usr/local/bin/

ADD go /opt/go
RUN cp /opt/go/bin/* /usr/local/bin/
ENV GOROOT=/opt/go GOPATH=/opt/gopath

WORKDIR /ws
CMD mongod --fork --logpath /var/log/mongodb.log && GOPROXY=off go test -mod=vendor ./...
```

## Conclusion

- Unit test should be a consensus
- Unit test in go is easy
- Mock make our testing efficient
- Service should be an interface for mock friendly
- Standard lib is good enough  (table testing, testing  HTTP)
- Other pacakges make testing better and documentation
- Docker can be used for complex situation

 
