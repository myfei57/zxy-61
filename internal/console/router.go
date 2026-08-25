package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"coldstore/internal/service"
)

func NewRouter(s *service.Service) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Get("/health", healthHandler(s))
	r.Get("/status", statusHandler(s))
	r.Get("/status/detailed", statusDetailedHandler(s))
	r.Post("/rooms", addRoomHandler(s))
	r.Post("/rooms/{id}/temperature", setTemperatureHandler(s))
	r.Post("/rooms/{id}/stock", stockInHandler(s))
	r.Get("/rooms/{id}/target", roomTargetHandler(s))
	r.Post("/rooms/{id}/target", setRoomTargetHandler(s))
	r.Get("/rooms/{id}/deviation", roomDeviationHandler(s))
	r.Get("/rooms/{id}/state", roomStateHandler(s))
	r.Post("/rooms/{id}/state", setRoomStateHandler(s))
	r.Get("/rooms/{id}/zone", zoneOfHandler(s))
	r.Post("/rooms/{id}/occupancy/remove", removeOccupancyHandler(s))
	r.Post("/start", startHandler(s))
	r.Post("/stage", stageHandler(s))
	r.Post("/defrost", defrostHandler(s))
	r.Post("/run-defrost", runDefrostHandler(s))
	r.Post("/reconcile", reconcileHandler(s))
	r.Post("/stage-by-load/{id}", stageByLoadHandler(s))
	r.Get("/room-for-machine/{id}", roomForMachineHandler(s))
	r.Post("/ammonia", ammoniaHandler(s))
	r.Post("/ammonia/calibrate", calibrateAmmoniaHandler(s))
	r.Post("/ammonia/level", ammoniaLevelHandler(s))
	r.Get("/ammonia/level/count", ammoniaLevelCountHandler(s))
	r.Get("/ammonia/margin", ammoniaMarginHandler(s))
	r.Post("/fan-verdict", fanVerdictHandler(s))
	r.Post("/fan/duty", fanDutyHandler(s))
	r.Post("/fan/runtime", fanRuntimeHandler(s))
	r.Get("/fan/runtime/due", fanRuntimeDueHandler(s))
	r.Post("/quota", quotaHandler(s))
	r.Post("/quota/reset", resetQuotaHandler(s))
	r.Post("/quota/window/record", quotaWindowRecordHandler(s))
	r.Post("/quota/window/allowed", quotaWindowAllowedHandler(s))
	r.Get("/quota/window/size", quotaWindowSizeHandler(s))
	r.Post("/quota/cycle/consume", quotaCycleConsumeHandler(s))
	r.Post("/quota/cycle/reset", quotaCycleResetHandler(s))
	r.Post("/cooling/start", coolingStartHandler(s))
	r.Post("/cooling/stop", coolingStopHandler(s))
	r.Post("/cooling/hours", coolingHoursHandler(s))
	r.Post("/cooling/capacity", updateCapacityHandler(s))
	r.Get("/registry/{key}", registryGetHandler(s))
	r.Get("/registry", registryKeysHandler(s))
	r.Get("/zones/{zone}/rooms", zoneRoomsHandler(s))
	r.Delete("/zones/{id}", removeRoomZoneHandler(s))
	r.Post("/valve/target", valveTargetHandler(s))
	r.Post("/band/step", bandStepHandler(s))
	r.Post("/save", saveHandler(s))
	r.Get("/audit", auditHandler(s))
	r.Get("/audit/export", auditExportHandler(s))
	return r
}
