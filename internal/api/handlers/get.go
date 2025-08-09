package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"go.uber.org/zap"
)

const (
	minValue = 1.0
	maxValue = 10000.0
)

func Get(rtp float64, logger *zap.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var multiplier float64

		prob := rand.Float64()
		boundary := minValue + (maxValue-minValue)*(1-rtp)

		if prob < rtp {
			multiplier = randFloat(boundary, maxValue)

		} else {
			multiplier = randFloat(minValue, boundary)
		}

		writeMultiplier(w, logger, multiplier)
		logger.Info("Get: successful get multiplier", zap.Float64("multiplier", multiplier))
	}

}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func writeMultiplier(w http.ResponseWriter, logger *zap.Logger, multiplier float64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := response{
		Multiplier: multiplier,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Error("writeMultiplier: failed to encoding response", zap.Error(err))
	}
}

type response struct {
	Result float64 `json:"result"`
}
