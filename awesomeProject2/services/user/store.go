package user

import (
	"awesomeProject2/types"
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE email = ?", email)
	u := &types.User{}
	err := scanRowIntoUser(row, u)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return u, nil
}
func (s *Store) FindBySearchName(lastname string) ([]*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE lastname LIKE ?", "%"+lastname+"%")
	if err != nil {
		return nil, err
	}
	users := make([]*types.User, 0)
	for rows.Next() {
		p, errs := scanRowsIntoUser(rows)

		if errs != nil {
			return nil, errs
		}
		users = append(users, p)
	}
	return users, nil
}
func (s *Store) GetAllUserId() ([]*types.User, error) {
	rows, err := s.db.Query("SELECT  * FROM users")
	if err != nil {
		return nil, err
	}
	users := make([]*types.User, 0)
	for rows.Next() {
		p, errs := scanRowsIntoUser(rows)

		if errs != nil {
			return nil, errs
		}

		users = append(users, p)
	}
	return users, nil

}
func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}
func (s *Store) UpdateUser(user types.UserUpdate, userId int) (string, error) {
	_, err := s.db.Exec("update users set firstname = ?, lastname = ? , email = ? , password = ? where  id = ?", user.FirstName, user.LastName, user.Email, user.Password, userId)
	if err != nil {
		return "update fail ", err
	}
	return "update successfully", nil
}
func (s *Store) GetUserByID(id int) (*types.User, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	u := &types.User{}
	err := scanRowIntoUser(row, u)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return u, nil
}
func (s *Store) DeleteUserByID(id int) error {
	var userID int
	err := s.db.QueryRow("SELECT id FROM users WHERE id = ?", id).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user with ID %d not found", id)
		}
		return err
	}

	result, err := s.db.Exec("DELETE  FROM users where id = ?", id)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
func (s *Store) GetAllUserIdPage(page int, limit int) ([]*types.User, error) {
	offset := (page - 1) * limit
	query := fmt.Sprintf("SELECT * FROM users LIMIT %d OFFSET %d", limit, offset)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	users := make([]*types.User, 0)
	for rows.Next() {
		p, errs := scanRowsIntoUser(rows)

		if errs != nil {
			return nil, errs
		}

		users = append(users, p)
	}
	return users, nil
}

func scanRowIntoUser(row *sql.Row, u *types.User) error {
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
