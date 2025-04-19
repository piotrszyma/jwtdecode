package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Ugly & fast way to expose ability for mocking time in tests.
var timeNow time.Time

const (
	colorReset  = "\033[0m"
	colorBlue   = "\033[34m" // Blue
	colorGreen  = "\033[32m" // Green
	colorYellow = "\033[33m" // Yellow
	colorCyan   = "\033[36m" // Cyan
	colorRed    = "\033[31m" // Red
	colorGray   = "\033[90m" // Dark Gray
)

// Helper function to base64 decode a string
func base64Decode(s string) ([]byte, error) {
	// Add padding if necessary.  JWT doesn't always include it.
	missing := len(s) % 4
	if missing != 0 {
		s += strings.Repeat("=", 4-missing)
	}
	return base64.StdEncoding.DecodeString(s)
}

func sortedMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys) // Sort the keys lexicographically
	return keys
}

func timeToISOString(value time.Time) string {
	// Format time.Time to ISO 8601 string
	return value.Format(time.RFC3339Nano) //Use RFC3339Nano for most precision
}

func float64TimestampToTimeUtc(timestamp float64) time.Time {
	// Convert float64 timestamp (seconds since epoch) to time.Time
	t := time.Unix(int64(timestamp), int64((timestamp-float64(int64(timestamp)))*float64(time.Second)))
	return t.UTC()
}

func humanReadableDelta(t time.Time) string {
	now := timeNow.UTC()
	diff := now.Sub(t)
	absDiff := diff
	if diff < 0 {
		absDiff = -diff
	}

	if absDiff < time.Minute {
		if diff < 0 {
			return "in less than a minute"
		}
		return "less than a minute ago"
	}

	if absDiff < time.Hour {
		minutes := int(absDiff.Minutes())
		if diff < 0 {
			return fmt.Sprintf("in %d minute(s)", minutes)
		}

		return fmt.Sprintf("%d minute(s) ago", minutes)
	}

	if absDiff < 24*time.Hour {
		hours := int(absDiff.Hours())
		if diff < 0 {
			return fmt.Sprintf("in %d hour(s)", hours)
		}

		if diff == 1 {
			return fmt.Sprintf("%d hour ago", hours)
		}

		return fmt.Sprintf("%d hour(s) ago", hours)
	}

	days := int(absDiff.Hours() / 24)
	if diff < 0 {
		return fmt.Sprintf("in %d day(s)", days)
	}
	return fmt.Sprintf("%d day(s) ago", days)
}

func printStructAsColoredJson(writer io.Writer, v interface{}) error {
	valueAsMap, ok := v.(map[string]any)
	if !ok {
		return fmt.Errorf("payload must be an object, got %s", v)
	}

	keysSorted := sortedMapKeys(valueAsMap)

	fmt.Fprintln(writer, "{")

	size := len(valueAsMap)
	item := 0

	for _, key := range keysSorted {
		value := valueAsMap[key]

		fmt.Fprint(writer, "  ")           // indent
		fmt.Fprint(writer, colorBlue)      // color of key
		fmt.Fprintf(writer, "\"%s\"", key) // key escaped
		fmt.Fprint(writer, colorReset)     // reset
		fmt.Fprint(writer, ": ")

		switch reflect.TypeOf(value).String() {
		case "bool":
			vAsBool := value.(bool)
			fmt.Fprint(writer, colorCyan) // color of string
			fmt.Fprintf(writer, "%s", strconv.FormatBool(vAsBool))
			fmt.Fprint(writer, colorReset) // reset

			if item != size-1 {
				fmt.Fprint(writer, ",")
			}

		case "float64":
			vAsFloat64 := value.(float64)
			fmt.Fprint(writer, colorYellow) // color of string
			fmt.Fprintf(writer, "%s", strconv.FormatFloat(vAsFloat64, 'f', -1, 64))
			fmt.Fprint(writer, colorReset) // reset

			if item != size-1 {
				fmt.Fprint(writer, ",")
			}

			if vAsFloat64 > 1_000_000_000 && vAsFloat64 < 10_000_000_000 {
				vAsTime := float64TimestampToTimeUtc(vAsFloat64)
				timeDeltaUntilNow := humanReadableDelta(vAsTime)

				fmt.Fprint(writer, " ")
				fmt.Fprint(writer, colorGray)
				fmt.Fprint(writer, "# ")
				fmt.Fprint(writer, timeToISOString(vAsTime))
				fmt.Fprint(writer, colorRed)
				fmt.Fprint(writer, " ")
				fmt.Fprint(writer, "[")
				fmt.Fprint(writer, timeDeltaUntilNow)
				fmt.Fprint(writer, "]")
				fmt.Fprint(writer, colorReset)
			}

		default:
			fmt.Fprint(writer, colorGreen) // color of string
			fmt.Fprintf(writer, "\"%s\"", value)
			fmt.Fprint(writer, colorReset) // reset

			if item != size-1 {
				fmt.Fprint(writer, ",")
			}
		}

		fmt.Fprint(writer, "\n")

		item += 1
	}

	fmt.Fprintln(writer, "}")

	return nil
}

func writeClaimsFromJwt(writer io.Writer, tokenString string) error {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return errors.New("error: invalid JWT format, expected 3 parts")
	}

	headerBytes, err := base64Decode(parts[0])
	if err != nil {
		return fmt.Errorf("error decoding header: %s", err)
	}
	var headerData interface{}
	err = json.Unmarshal(headerBytes, &headerData)
	if err != nil {
		return fmt.Errorf("error unmarshaling header: %s", err)
	}

	fmt.Fprintln(writer, "Header:")

	if err := printStructAsColoredJson(writer, headerData); err != nil {
		return fmt.Errorf("failed to print struct = %s", err)
	}

	payloadBytes, err := base64Decode(parts[1])
	if err != nil {
		return fmt.Errorf("error decoding payload: %s", err)
	}

	var payloadData interface{} // Use interface{} to handle arbitrary JSON structure.
	err = json.Unmarshal(payloadBytes, &payloadData)
	if err != nil {
		return fmt.Errorf("error unmarshaling payload: %s", err)
	}

	fmt.Fprintln(writer)

	fmt.Fprintln(writer, "Payload:")

	if err := printStructAsColoredJson(writer, payloadData); err != nil {
		return fmt.Errorf("failed to print payload = %s", err)
	}

	fmt.Fprintln(writer)

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: jwtdecode <jwt_token>")
		return
	}

	tokenString := os.Args[1]

	if err := writeClaimsFromJwt(os.Stdout, tokenString); err != nil {
		fmt.Printf("failed to write claims = %d", err)
	}
}

func init() {
	timeNow = time.Now()
}
