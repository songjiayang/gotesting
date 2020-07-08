package post

import (
	"database/sql"
	"encoding/json"
	"net/http"

	. "gotesting/testinghttp"
)

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
