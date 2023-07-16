package helpers

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"packettrackingnet/config"
	"packettrackingnet/dto"
	"strconv"
	"strings"
)

var logger *log.Logger

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

func ReadUploadedCSV(file io.Reader, header bool) ([][]string, error) {
	reader := csv.NewReader(file)
	if !header {
		if _, err := reader.Read(); err != nil {
			return [][]string{}, errors.New(err.Error())
		}
	}
	records, err := reader.ReadAll()
	if err != nil {
		return [][]string{}, errors.New(err.Error())
	}

	return records, nil
}

func WriteCSV(header []string, data [][]string) (*os.File, error) {
	file, err := os.CreateTemp("", "tmp_")

	if err != nil {
		return nil, errors.New(err.Error())
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	completeData := make([][]string, 0)
	completeData = append(completeData, header)
	completeData = append(completeData, data[:]...)
	for _, record := range completeData {
		if err := writer.Write(record); err != nil {
			return nil, errors.New(err.Error())
		}
	}
	file.Seek(0, 0)
	return file, nil
}

func InitErrorLogger() {
	file, err := os.OpenFile("./error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer file.Close()

	logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
}

func LogError(err error) {
	logger.Fatal(err.Error())
}
