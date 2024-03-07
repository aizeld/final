package models

import (
	"database/sql"
	"errors"
	"time"
)

type FighterModelInterface interface {
	Insert(name string, wrestling int, striking int, stamina int) (int, error)
	Get(id int) (*Fighter, error)
	Latest() ([]*Fighter, error)
	Delete(id int) error
	Update(id int, name string, wrestling int, striking int, stamina int) (int, error)
}

type Fighter struct {
	ID        int
	Name      string
	Wrestling int
	Striking  int
	Stamina   int
	Created   time.Time
}

type FighterModel struct {
	DB *sql.DB
}

func (m *FighterModel) Delete(id int) error {
	stmt := `delete FROM fighters WHERE id = ?`
	result, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	Affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if Affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
func (m *FighterModel) Update(id int, name string, wrestling int, striking int, stamina int) (int, error) {
	stmt := `Update fighters set name = ?, wrestling = ?, striking= ?, stamina=? where id = ?`
	result, err := m.DB.Exec(stmt, name, wrestling, striking, stamina, id)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	if affected == 0 {
		return 0, sql.ErrNoRows
	}

	return int(id), nil

}
func (m *FighterModel) Insert(name string, wrestling int, striking int, stamina int) (int, error) {
	stmt := `INSERT INTO fighters (name, wrestling, striking, stamina, created) VALUES(?, ?, ?, ?, UTC_TIMESTAMP())`
	result, err := m.DB.Exec(stmt, name, wrestling, striking, stamina)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *FighterModel) Get(id int) (*Fighter, error) {
	stmt := `SELECT id, name, wrestling, striking, stamina, created FROM fighters WHERE id = ?`
	row := m.DB.QueryRow(stmt, id)

	f := &Fighter{}
	err := row.Scan(&f.ID, &f.Name, &f.Wrestling, &f.Striking, &f.Stamina, &f.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return f, nil
}

func (m *FighterModel) Latest() ([]*Fighter, error) {
	stmt := `SELECT id, name, wrestling, striking, stamina, created FROM fighters ORDER BY created DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fighters := []*Fighter{}
	for rows.Next() {
		f := &Fighter{}
		err := rows.Scan(&f.ID, &f.Name, &f.Wrestling, &f.Striking, &f.Stamina, &f.Created)
		if err != nil {
			return nil, err
		}
		fighters = append(fighters, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return fighters, nil
}
