package clientdb

import (
	"database/sql"
	"github.com/tsawler/goblender/pkg/models"
)

// PageModel wraps database
type PageModel struct {
	DB *sql.DB
}

// AllPages returns slice of pages from goBlender's database
func (m *PageModel) AllPages() ([]*models.Page, error) {
	stmt := "SELECT id, page_title, active, slug, created_at, updated_at FROM pages ORDER BY page_title"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []*models.Page

	for rows.Next() {
		var s models.Page
		err = rows.Scan(&s.ID, &s.PageTitle, &s.Active, &s.Slug, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		pages = append(pages, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pages, nil
}
