package checkout

import (
	"database/sql"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}
func (s *Store) UpdateStatusAdmin(orderItemId int, status string) (string, error) {
	_, err := s.db.Exec("update orders set status = ? where id = ?", status, orderItemId)
	if err != nil {
		return "update fail", err
	}
	return "update successfully", nil
}
