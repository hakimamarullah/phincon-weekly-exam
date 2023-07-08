package utils

import (
	"errors"
	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
	"strconv"
	"strings"
)

func GenerateId() string {
	return uuid.New().String()
}

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

func ScanString(label string) string {
	prompt := promptui.Prompt{Label: label}
	val, _ := prompt.Run()
	return val
}

func ScanStringNonBlank(label string) string {
	prompt := promptui.Prompt{Label: label, Validate: StringNotBlank}
	val, _ := prompt.Run()
	return val
}

func ScanInt(label string) int {
	prompt := promptui.Prompt{Label: label}
	in, _ := prompt.Run()
	val, _ := strconv.Atoi(in)
	return val
}

func ConfirmInput(label string) bool {
	prompt := promptui.Prompt{Label: label, IsConfirm: true}
	in, _ := prompt.Run()
	return strings.EqualFold(in, "y")
}

func StringNotBlank(input string) error {
	if len(input) == 0 {
		return errors.New("input can't be blank")
	}

	return nil
}

func GetStatus(isReceived bool) string {
	if isReceived {
		return "RECEIVED"
	}
	return "ON PROCESS"
}
