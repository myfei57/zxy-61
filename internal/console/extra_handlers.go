package console

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"coldstore/internal/room"
	"coldstore/internal/service"
)

func statusDetailedHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, s.DetailedStatus())
	}
}

func roomTargetHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		target, ok := s.RoomTarget(id)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "room not found"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]float64{"target": target})
	}
}

func setRoomTargetHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var payload struct {
			Target float64 `json:"target"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if !s.SetRoomTarget(id, payload.Target) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "room not found"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

func roomDeviationHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		deviation, ok := s.RoomDeviation(id)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "room not found"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]float64{"deviation": deviation})
	}
}

func removeOccupancyHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var payload struct {
			Amount float64 `json:"amount"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		stock, ok := s.RemoveOccupancy(id, payload.Amount)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "room not found"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]float64{"stock": stock})
	}
}

func roomStateHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		writeJSON(w, http.StatusOK, map[string]string{"state": string(s.RoomState(id))})
	}
}

func setRoomStateHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var payload struct {
			State string `json:"state"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		s.SetRoomState(id, mapRoomState(payload.State))
		writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

func zoneOfHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		zone, ok := s.ZoneOf(id)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "room not zoned"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"zone": zone})
	}
}

func auditExportHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := s.ExportAudit()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}
}

func mapRoomState(state string) room.State {
	switch state {
	case "cooling", "defrost", "normal":
		return room.State(state)
	default:
		return room.StateNormal
	}
}

func coolingStartHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.CoolingStart()
		writeJSON(w, http.StatusOK, map[string]string{"status": "started"})
	}
}

func coolingStopHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.CoolingStop()
		writeJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
	}
}

func coolingHoursHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Hours float64 `json:"hours"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		s.CoolingAddHours(payload.Hours)
		writeJSON(w, http.StatusOK, map[string]string{"status": "recorded"})
	}
}

func updateCapacityHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Tons     float64 `json:"tons"`
			Setpoint float64 `json:"setpoint"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		s.UpdateCapacity(payload.Tons, payload.Setpoint)
		writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

func fanDutyHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Active float64 `json:"active"`
			Total  float64 `json:"total"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		s.AddFanDuty(payload.Active, payload.Total)
		writeJSON(w, http.StatusOK, map[string]string{"status": "recorded"})
	}
}

func fanRuntimeHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Hours float64 `json:"hours"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		s.AddFanRuntime(payload.Hours)
		writeJSON(w, http.StatusOK, map[string]string{"status": "recorded"})
	}
}

func fanRuntimeDueHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		interval, _ := strconv.ParseFloat(r.URL.Query().Get("interval"), 64)
		writeJSON(w, http.StatusOK, map[string]bool{"due": s.FanRuntimeDue(interval)})
	}
}

func ammoniaLevelHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Value float64 `json:"value"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		s.AddAmmoniaLevel(payload.Value)
		writeJSON(w, http.StatusOK, map[string]string{"status": "recorded"})
	}
}

func ammoniaLevelCountHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]int{"count": s.AmmoniaLevelCount()})
	}
}

func ammoniaMarginHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]float64{"margin": s.AmmoniaMargin()})
	}
}

func quotaWindowRecordHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Amount int `json:"amount"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		s.RecordQuotaWindow(payload.Amount)
		writeJSON(w, http.StatusOK, map[string]string{"status": "recorded"})
	}
}

func quotaWindowAllowedHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Amount int `json:"amount"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"allowed": s.QuotaWindowAllowed(payload.Amount)})
	}
}

func quotaWindowSizeHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]int{"size": s.QuotaWindowSize()})
	}
}

func quotaCycleConsumeHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Amount int `json:"amount"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"allowed": s.ConsumeQuotaCycle(payload.Amount)})
	}
}

func quotaCycleResetHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.ResetQuotaCycle()
		writeJSON(w, http.StatusOK, map[string]string{"status": "reset"})
	}
}

func registryGetHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		value, ok := s.RegistryGet(key)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "key not found"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"value": value})
	}
}

func registryKeysHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string][]string{"keys": s.RegistryKeys()})
	}
}

func zoneRoomsHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		zone := chi.URLParam(r, "zone")
		writeJSON(w, http.StatusOK, map[string][]string{"rooms": s.ZoneRooms(zone)})
	}
}

func removeRoomZoneHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		s.RemoveRoomFromZone(id)
		writeJSON(w, http.StatusOK, map[string]string{"status": "removed"})
	}
}

func valveTargetHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Target float64 `json:"target"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		s.SetValveTarget(payload.Target)
		writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

func bandStepHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Amount int `json:"amount"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]int{"band": s.StepBand(payload.Amount)})
	}
}
