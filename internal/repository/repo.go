package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/elgntt/effective-mobile/internal/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

var allowedParams = map[string]string{
	"male":   " AND gender = 'male'",
	"female": " AND gender = 'female'",
}

func createFilterString(params model.Params) (string, []interface{}) {
	var (
		builder        strings.Builder
		queryParams    []interface{}
		parameterCount = 1
	)

	builder.WriteString(allowedParams[params.Gender])

	if params.MinAge > 0 {
		builder.WriteString(fmt.Sprintf(" AND age >= $%d", parameterCount))
		queryParams = append(queryParams, params.MinAge)
		parameterCount++
	}
	if params.MaxAge > 0 {
		builder.WriteString(fmt.Sprintf(" AND age <= $%d", parameterCount))
		queryParams = append(queryParams, params.MaxAge)
		parameterCount++
	}

	if params.Name != "" {
		builder.WriteString(fmt.Sprintf(" AND LOWER(name) = $%d", parameterCount))
		queryParams = append(queryParams, params.Name)
		parameterCount++
	}

	return builder.String(), queryParams
}

func (r *Repo) AddPeople(ctx context.Context, data model.NewPeople) error {
	query := `INSERT INTO people (name, surname, patronymic, age, gender, country_id)
			  VALUES (:name, :surname, :patronymic, :age, :gender, :country_id)`

	_, err := r.db.NamedExecContext(ctx, query, data)
	if err != nil {
		return fmt.Errorf("SQL: AddPeople(): Exec: %w", err)
	}

	return nil
}

func (r *Repo) GetPeople(ctx context.Context, params model.Params) ([]model.People, int, error) {
	filterStr, args := createFilterString(params)
	countQuery := `SELECT COUNT(*)
				   FROM people
				   WHERE 1 = 1` + filterStr

	var totalCount int
	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("SQL: GetPeople(): Count: %w", err)
	}

	peopleQuery := `SELECT people_id, name, surname, patronymic, age, gender, country_id
					FROM people
					WHERE 1 = 1` + filterStr

	peopleQuery += fmt.Sprintf(` LIMIT %d OFFSET %d`, params.Limit, params.Offset)

	var people []model.People

	err = r.db.SelectContext(ctx, &people, peopleQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("SQL: GetPeople(): Select: %w", err)
	}

	return people, totalCount, nil
}

func (r *Repo) DeletePeople(ctx context.Context, id int) error {
	query := `DELETE FROM people
			  WHERE people_id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("SQL: DeletePeople(): Exec: %w", err)
	}

	return nil
}

func (r *Repo) GetPeopleById(ctx context.Context, id int) (*model.People, error) {
	query := `SELECT people_id, name, surname, patronymic, age, gender, country_id 
			  FROM people
			  WHERE people_id = $1`

	people := model.People{}
	err := r.db.GetContext(ctx, &people, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return &model.People{}, err
	}

	return &people, nil
}

func (r *Repo) UpdatePeople(ctx context.Context, data model.People) error {
	query := `UPDATE people SET name = :name, surname = :surname, patronymic = :patronymic,
								age = :age, gender = :gender, country_id = :country_id
			  WHERE people_id = :people_id`

	_, err := r.db.NamedExecContext(ctx, query, data)
	if err != nil {
		return fmt.Errorf("SQL: UpdatePeople(): NamedExec: %w", err)
	}

	return nil
}
