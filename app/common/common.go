package common

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

func GetStr(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

func GetArr[T any](ptr *[]T) []T {
	if ptr != nil {
		return *ptr
	}
	return []T{}
}

func GetInt(ptr *int64) int64 {
	if ptr != nil {
		return *ptr
	}
	return 0
}

func GetInt32(ptr *int32) int32 {
	if ptr != nil {
		return *ptr
	}
	return 0
}

func GetBool(ptr *bool) bool {
	if ptr != nil {
		return *ptr
	}
	return false
}

func GetDecimal(ptr *decimal.Decimal) decimal.Decimal {
	if ptr != nil {
		return *ptr
	}
	return decimal.NewFromInt(0)
}

func GetTime(ptr *time.Time) time.Time {
	if ptr != nil {
		return *ptr
	}
	return time.Unix(0, 0)
}

func ClearQuery(input string) string {
	// Экранируем одиночные кавычки
	input = strings.ReplaceAll(input, ";", " ")

	// Можно добавить другие экранирования по необходимости
	return input
}

type ErrorMessage struct {
	Error string `json:"error"`
}

func ServerError(c *fiber.Ctx, text string) error {
	println(text)
	msg := ErrorMessage{
		Error: text,
	}
	return c.Status(fiber.StatusInternalServerError).JSON(msg)
}

func GetVendorId() int8 {
	id_str := os.Getenv("VENDOR_VENDOR_ID")
	num, err := strconv.Atoi(id_str)
	if err != nil {
		return 3
	}

	return int8(num)
}

func GetPackageFromGlobalId(globalId string) (int64, error) { // packageId, err
	reg := regexp.MustCompile(`vendor-(\d+)-\w+-\w+`)
	match := reg.FindStringSubmatch(globalId)

	if len(match) < 2 {
		return 0, fmt.Errorf("invalid globalId format")
	}

	packageId, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, err
	}

	return int64(packageId), nil
}

func IsValidICCID(iccid string) bool {
	match, _ := regexp.MatchString(`^\d{19,20}$`, iccid)
	return match
}

// --------------------------------------------------------------------------------
func ExtractLpa(iosSetup string) string {
	parts := strings.Split(iosSetup, "carddata=")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

func ParseTime(timeStr string) *time.Time {
	if timeStr == "" {
		return nil
	}
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil
	}
	return &parsedTime
}
