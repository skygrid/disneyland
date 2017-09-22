package optimus

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type OptimusStorageConfig struct {
	DatabaseURI string `json:"db_uri"`
}

type OptimusStorage struct {
	db     *sql.DB
	Config OptimusStorageConfig
}

func NewOptimusStorage(db_uri string) (*OptimusStorage, error) {
	ret := &OptimusStorage{
		Config: OptimusStorageConfig{DatabaseURI: db_uri},
	}

	err := ret.Connect()

	return ret, err
}

func (storage *OptimusStorage) Connect() error {
	db, err := sql.Open("postgres", storage.Config.DatabaseURI)
	storage.db = db
	return err
}

func (storage *OptimusStorage) CreatePoint(point *Point) (*Point, error) {
	tx, err := storage.db.Begin()
	if err != nil {
		return nil, err
	}

	created_point := &Point{}
	err = tx.QueryRow(`
		INSERT INTO points (project, status, coordinate, metric_value, metadata)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, project, status, coordinate, metric_value, metadata;`,
		point.Project, point.Status, point.Coordinate, point.MetricValue, point.Metadata,
	).Scan(
		&created_point.Id,
		&created_point.Project,
		&created_point.Status,
		&created_point.Coordinate,
		&created_point.MetricValue,
		&created_point.Metadata,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}
	return created_point, err
}