package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorKey    = "\033[34m" // Blue
	colorString = "\033[32m" // Green
	colorNumber = "\033[33m" // Yellow
	colorBool   = "\033[36m" // Cyan
	colorNull   = "\033[31m" // Red
	colorGray   = "\033[90m" // Dark Gray
	colorError  = "\033[31m" // Red
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

func float64ToISOString(timestamp float64) string {
	// Convert float64 timestamp (seconds since epoch) to time.Time
	t := time.Unix(int64(timestamp), int64((timestamp-float64(int64(timestamp)))*float64(time.Second)))

	// Format time.Time to ISO 8601 string
	isoString := t.UTC().Format(time.RFC3339Nano) //Use RFC3339Nano for most precision

	return isoString
}

func printStructAsColoredJson(v interface{}) error {
	valueAsMap := v.(map[string]any)
	keysSorted := sortedMapKeys(valueAsMap)

	fmt.Println("{")

	size := len(valueAsMap)
	item := 0

	for _, key := range keysSorted {
		value := valueAsMap[key]

		fmt.Print("  ")           // indent
		fmt.Print(colorKey)       // color of key
		fmt.Printf("\"%s\"", key) // key escaped
		fmt.Print(colorReset)     // reset
		fmt.Print(": ")

		switch reflect.TypeOf(value).String() {
		case "bool":
			vAsBool := value.(bool)
			fmt.Print(colorBool) // color of string
			fmt.Printf("%s", strconv.FormatBool(vAsBool))
			fmt.Print(colorReset) // reset
		case "float64":
			vAsFloat64 := value.(float64)
			fmt.Print(colorNumber) // color of string
			fmt.Printf("%s", strconv.FormatFloat(vAsFloat64, 'f', -1, 64))
			fmt.Print(colorReset) // reset

			if vAsFloat64 > 1_000_000_000 && vAsFloat64 < 10_000_000_000 {
				fmt.Print("  ")
				fmt.Print(colorGray)
				fmt.Print("# ")
				fmt.Print(float64ToISOString(vAsFloat64))
				fmt.Print(colorReset)
			}
		default:
			fmt.Print(colorString) // color of string
			fmt.Printf("\"%s\"", value)
			fmt.Print(colorReset) // reset
		}

		if item != size-1 {
			fmt.Print(",")
		}

		fmt.Print("\n")

		item += 1
	}

	fmt.Println("}")

	return nil
}

func main() {
	// 1. Check for the JWT token as a command-line argument.
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <jwt_token>")
		return
	}
	tokenString := os.Args[1]

	// 2. Split the JWT token by the '.' character.
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		fmt.Println("Error: Invalid JWT format. Expected 3 parts.")
		return
	}

	// 3. Base64 decode the header (part 1).
	headerBytes, err := base64Decode(parts[0])
	if err != nil {
		fmt.Println("Error decoding header:", err)
		return
	}

	// 4. Base64 decode the payload (part 2).
	payloadBytes, err := base64Decode(parts[1])
	if err != nil {
		fmt.Println("Error decoding payload:", err)
		return
	}

	// 5. JSON unmarshal the header.
	var headerData interface{} // Use interface{} to handle arbitrary JSON structure.
	err = json.Unmarshal(headerBytes, &headerData)
	if err != nil {
		fmt.Println("Error unmarshaling header:", err)
		return
	}

	// 6. JSON unmarshal the payload.
	var payloadData interface{} // Use interface{} to handle arbitrary JSON structure.
	err = json.Unmarshal(payloadBytes, &payloadData)
	if err != nil {
		fmt.Println("Error unmarshaling payload:", err)
		return
	}

	// 7. Print the header and payload.
	fmt.Println("Header:")

	printStructAsColoredJson(headerData)

	fmt.Println()

	fmt.Println("Payload:")

	printStructAsColoredJson(payloadData)

	fmt.Println()
}
