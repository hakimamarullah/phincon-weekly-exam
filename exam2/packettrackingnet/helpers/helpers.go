package helpers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"os"
	"packettrackingnet/config"
	"packettrackingnet/dto"
	"strconv"
	"strings"
)

func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateIdLocation specialized function to generate incremental Location ID.
func GenerateIdLocation(lastIndex int) string {
	stringId := strconv.Itoa(lastIndex)

	var sb strings.Builder
	sb.WriteString("GDNG-")
	for i := 0; i < (5 - len(stringId)); i++ {
		sb.WriteString("0")
	}
	sb.WriteString(stringId)
	return sb.String()
}

func ResponseJSON(w http.ResponseWriter, payload dto.ResponseBody) {
	if payload.Message == "" {
		payload.Message = "Success"
	}

	if payload.Code == 0 {
		payload.Code = http.StatusOK
	}

	response, _ := json.Marshal(payload)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(payload.Code)
	w.Write(response)
}

func TruncateDatabase() error {
	var filePaths = []string{config.LOCATION,
		config.PACKET,
		config.RECEIVER,
		config.SENDER,
		config.SHIPMENT,
		config.SERVICE}

	for _, item := range filePaths {
		if err := os.Truncate(item, 0); err != nil {
			return err
		}
		writer, _ := os.Create(item)
		json.NewEncoder(writer).Encode(make([]any, 0))
	}
	return nil
}
