package store

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/russross/meddler"
	sq "github.com/Masterminds/squirrel"
)

type FeatureStatusHistory struct {
	Id              int64      `json:"id"                            meddler:"id,pk"`
	CreatedAt       *time.Time `json:"created_at,omitempty"          meddler:"created_at"`
	Enabled         *bool      `json:"enabled,omitempty"             meddler:"enabled"`
	FeatureId       *int64     `json:"feature_id,omitempty"          meddler:"feature_id"`
	EnvironmentId   *int64     `json:"environment_id,omitempty"      meddler:"environment_id"`
	FeatureStatusId *int64     `json:"feature_status_id,omitempty"   meddler:"feature_status_id"`
}

func (s *store) ListFeatureStatusHistory(featureId *int64, environmentId *int64, featureStatusId *int64) ([]*FeatureStatusHistory, error) {
	query := sq.Select("*").From("feature_status_history")

	if featureId != nil {
		query = query.Where(sq.Eq{"feature_id": featureId})
	}

	if environmentId != nil {
		query = query.Where(sq.Eq{"environment_id": environmentId})
	}

	if featureStatusId != nil {
		query = query.Where(sq.Eq{"environment_id": environmentId})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	log.Debug(sql)

	featuresStatusHistory := []*FeatureStatusHistory{}
	err = meddler.QueryAll(s.db, &featuresStatusHistory, sql, args...)

	return featuresStatusHistory, err
}
