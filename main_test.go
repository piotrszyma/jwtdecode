package main

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/andreyvit/diff"
)

func removeColorCodes(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(s, "")
}

func TestPrintsJwtData(t *testing.T) {
	// Arrange.
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"

	timeTest, _ := time.Parse(time.RFC3339, "2025-04-19T18:25:00+02:00")

	// Ugly way to fake time now for tests by using global variable.
	timeNow = timeTest

	s := new(strings.Builder)

	// Act.
	err := writeClaimsFromJwt(s, jwt)

	if err != nil {
		t.Errorf("failed = %s", err)
	}

	actual := removeColorCodes(s.String())
	expected := `Header:
{
  "alg": "HS256",
  "typ": "JWT"
}

Payload:
{
  "admin": true,
  "iat": 1516239022, # 2018-01-18T01:30:22Z [2648 day(s) ago]
  "name": "John Doe",
  "sub": "1234567890"
}

`

	if actual != expected {
		t.Errorf("Result not as expected:\n%v", diff.CharacterDiff(expected, actual))
	}
}
