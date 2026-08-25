package console

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"coldstore/internal/service"
)

func healthHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, s.Health())
	}
}

func statusHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, s.Status())
	}
}

func addRoomHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			ID          string  `json:"id"`
			Name        string  `json:"name"`
			Temperature float64 `json:"temperature"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := s.AddRoom(payload.ID, payload.Name, payload.Temperature); err != nil {
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"id": payload.ID})
	}
}

func setTemperatureHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var payload struct {
			Temperature float64 `json:"temperature"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := s.SetRoomTemperature(id, payload.Temperature); err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"id": id})
	}
}

func stockInHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var payload struct {
			Added float64 `json:"added"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		baseline, err := s.StockIn(id, payload.Added)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]float64{"baseline": baseline})
	}
}

func startHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.StartPlant(); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "started"})
	}
}

func stageHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.StageLoad(); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "staged"})
	}
}

func defrostHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.Defrost(); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "defrosting"})
	}
}

func runDefrostHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.RunDefrost(); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "running"})
	}
}

func reconcileHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		threshold, _ := strconv.ParseFloat(r.URL.Query().Get("threshold"), 64)
		stages := s.Reconcile(threshold)
		writeJSON(w, http.StatusOK, map[string]int{"stages": stages})
	}
}

func stageByLoadHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		stages, err := s.StageByLoad(id)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]int{"stages": stages})
	}
}

func roomForMachineHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		roomID, ok := s.RoomForMachine(id)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "machine not mapped"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"room": roomID})
	}
}

func ammoniaHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Reading float64 `json:"reading"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"tripped": s.CheckInterlock(payload.Reading)})
	}
}

func calibrateAmmoniaHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Threshold float64 `json:"threshold"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		s.CalibrateAmmonia(payload.Threshold)
		writeJSON(w, http.StatusOK, map[string]string{"status": "calibrated"})
	}
}

func fanVerdictHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Temperature float64 `json:"temperature"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"over_temperature": s.FanVerdict(payload.Temperature)})
	}
}

func quotaHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Amount int `json:"amount"`
		}
		if err := readJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"allowed": s.QuotaAllow(payload.Amount)})
	}
}

func resetQuotaHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.ResetQuota()
		writeJSON(w, http.StatusOK, map[string]string{"status": "reset"})
	}
}

func saveHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.Save(); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "saved"})
	}
}

func auditHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, s.ListAudit())
	}
}
