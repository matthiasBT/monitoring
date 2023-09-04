package usecases

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	common "github.com/matthiasBT/monitoring/internal/infra/entities"
)

func (c *BaseController) updateMetric(w http.ResponseWriter, r *http.Request) {
	asJSON := r.Header.Get("Content-Type") == "application/json"
	var metrics *common.Metrics
	if metrics = parseMetric(r, asJSON, true); metrics == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := metrics.Validate(true)
	if err != nil {
		handleInvalidMetric(w, err)
		return
	}

	result, err := UpdateMetric(r.Context(), c, metrics)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			w.WriteHeader(http.StatusBadRequest) // duplicate
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := writeMetric(w, asJSON, result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *BaseController) getMetric(w http.ResponseWriter, r *http.Request) {
	asJSON := r.Header.Get("Content-Type") == "application/json"
	var metrics *common.Metrics
	if metrics = parseMetric(r, asJSON, false); metrics == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := metrics.Validate(false)
	if err != nil {
		handleInvalidMetric(w, err)
		return
	}

	result, err := GetMetric(r.Context(), c, metrics)
	if err != nil {
		var status int
		if errors.Is(err, common.ErrUnknownMetric) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		w.WriteHeader(status)
		w.Write([]byte(err.Error()))
		return
	}

	if err := writeMetric(w, asJSON, result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *BaseController) getAllMetrics(w http.ResponseWriter, r *http.Request) {
	result, err := GetAllMetrics(r.Context(), c, "all_metrics.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(result.Bytes())
}

func (c *BaseController) ping(w http.ResponseWriter, r *http.Request) {
	if c.DBManager == nil {
		c.Logger.Errorf("Failed to ping the databases: no DB manager\n")
		w.WriteHeader(http.StatusBadRequest)
	}
	if err := c.DBManager.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
