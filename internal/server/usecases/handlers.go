package usecases

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	common "github.com/matthiasBT/monitoring/internal/infra/entities"
)

func (c *BaseController) updateMetric(w http.ResponseWriter, r *http.Request) {
	asJSON := r.Header.Get("Content-Type") == "application/json"
	var metrics *common.Metrics
	if metrics = parseMetric(r, asJSON, true); metrics == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse metric"))
		return
	}

	if err := metrics.Validate(true); err != nil {
		handleInvalidMetric(w, err)
		return
	}

	result, err := UpdateMetric(r.Context(), c, metrics)
	if err != nil {
		if c.handleDuplicate(w, err) {
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := writeMetric(w, asJSON, result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func (c *BaseController) getMetric(w http.ResponseWriter, r *http.Request) {
	asJSON := r.Header.Get("Content-Type") == "application/json"
	var metrics *common.Metrics
	if metrics = parseMetric(r, asJSON, false); metrics == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse metric"))
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
		w.Write([]byte(err.Error()))
	}
}

func (c *BaseController) getAllMetrics(w http.ResponseWriter, r *http.Request) {
	result, err := GetAllMetrics(r.Context(), c, "all_metrics.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(result.Bytes())
}

func (c *BaseController) massUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Supply data as JSON"))
		return
	}

	var batch []*common.Metrics
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err := json.Unmarshal(body, &batch); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	for _, metrics := range batch {
		if err := metrics.Validate(true); err != nil {
			handleInvalidMetric(w, err)
			return
		}
	}

	if err := MassUpdate(r.Context(), c, batch); err != nil {
		if c.handleDuplicate(w, err) {
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *BaseController) ping(w http.ResponseWriter, r *http.Request) {
	if c.DBManager == nil {
		c.Logger.Errorf("Failed to ping the databases: no DB manager\n")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No database was configured"))
	}
	if err := c.DBManager.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
}

func (c *BaseController) handleDuplicate(w http.ResponseWriter, err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		w.WriteHeader(http.StatusBadRequest) // duplicate
		w.Write([]byte("A metric with the same name and another type already exists"))
		return true
	}
	return false
}
