package utils

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/manifoldco/promptui"
	"strconv"
	"strings"
)

// GenerateId function to generate random id for the struct.
func GenerateId() string {
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

// ScanStringNonBlank function to take input from the terminal using promptui.
// This function validates the string input to has length greater than 0.
func ScanStringNonBlank(label string) string {
	prompt := promptui.Prompt{Label: label, Validate: StringNotBlank}
	val, _ := prompt.Run()
	return val
}

// ScanInt function to take string input from the terminal using promptui
// returning value of type int.
func ScanInt(label string) int {
	prompt := promptui.Prompt{Label: label}
	in, _ := prompt.Run()
	val, _ := strconv.Atoi(in)
	return val
}

// ConfirmInput function to ask confirmation from the terminal.
// The answer is only yes or no.
// It returns boolean value.
func ConfirmInput(label string) bool {
	prompt := promptui.Prompt{Label: label, IsConfirm: true}
	in, _ := prompt.Run()
	return strings.EqualFold(in, "y")
}

// StringNotBlank function to validate that the input string has a length greater than 0.
// This function will be used in promptui for validation.
func StringNotBlank(input string) error {
	if len(input) == 0 {
		return errors.New("input can't be blank")
	}

	return nil
}

// GetStatus function to convert flag IsReceived to string "RECEIVED" or "ON PROCESS"
// Returns "RECEIVED" if the argument is true.
func GetStatus(isReceived bool) string {
	if isReceived {
		return "RECEIVED"
	}
	return "ON PROCESS"
}

// StandardTable function to create a new table Writer.
// Returns a table Writer with basic styling.
func StandardTable(title string) table.Writer {
	t := table.NewWriter()
	t.SetTitle(strings.ToTitle(title))
	t.SetAutoIndex(
		true)
	t.Style().Format.Header = text.FormatTitle
	t.Style().Format.Footer = text.FormatTitle
	t.SetStyle(table.StyleBold)

	return t
}
