package expect_test

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/crumbandbase/expect"
)

type starship struct {
	Name string `json:"name"`
}

var (
	enterprise = starship{"enterprise"}
	voyager    = starship{"voyager"}
)

func TestEqual(t *testing.T) {
	t.Run("succeeds when the expected value equals the actual value", func(t *testing.T) {
		test := &testing.T{}

		captureOutput(func() {
			expect.Equal(test, enterprise, enterprise)
		})

		if test.Failed() {
			t.Error("values were not equal")
		}
	})

	t.Run("fails when the expected value does not equal the actual value", func(t *testing.T) {
		test := &testing.T{}

		captureOutput(func() {
			expect.Equal(test, enterprise, voyager)
		})

		if !test.Failed() {
			t.Error("values were equal")
		}
	})
}

func TestNotEqual(t *testing.T) {
	t.Run("succeeds when the expected value does not equal the actual value", func(t *testing.T) {
		test := &testing.T{}

		captureOutput(func() {
			expect.NotEqual(test, enterprise, voyager)
		})

		if test.Failed() {
			t.Error("values were equal")
		}
	})

	t.Run("fails when the expected value equals the actual value", func(t *testing.T) {
		test := &testing.T{}

		captureOutput(func() {
			expect.NotEqual(test, enterprise, enterprise)
		})

		if !test.Failed() {
			t.Error("values were not equal")
		}
	})
}

func TestStreamEqual(t *testing.T) {
	t.Run("succeeds when the expected response equals the actual response", func(t *testing.T) {
		test := &testing.T{}
		w := httptest.NewRecorder()

		b, err := json.Marshal(enterprise)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}

		w.Write(b)

		captureOutput(func() {
			expect.StreamEqual(test, json.NewDecoder(w.Body), &enterprise)
		})

		if test.Failed() {
			t.Error("values were not equal")
		}
	})

	t.Run("fails when the expected response does not equal the actual response", func(t *testing.T) {
		test := &testing.T{}
		w := httptest.NewRecorder()

		b, err := json.Marshal(enterprise)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}

		w.Write(b)

		captureOutput(func() {
			expect.StreamEqual(test, json.NewDecoder(w.Body), &voyager)
		})

		if !test.Failed() {
			t.Error("values were equal")
		}
	})
}

func TestStreamNotEqual(t *testing.T) {
	t.Run("succeeds when the expected response does not equal the actual response", func(t *testing.T) {
		test := &testing.T{}
		w := httptest.NewRecorder()

		b, err := json.Marshal(enterprise)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}

		w.Write(b)

		captureOutput(func() {
			expect.StreamNotEqual(test, json.NewDecoder(w.Body), &voyager)
		})

		if test.Failed() {
			t.Error("values were equal")
		}
	})

	t.Run("fails when the expected response equals the actual response", func(t *testing.T) {
		test := &testing.T{}
		w := httptest.NewRecorder()

		b, err := json.Marshal(enterprise)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}

		w.Write(b)

		captureOutput(func() {
			expect.StreamNotEqual(test, json.NewDecoder(w.Body), &enterprise)
		})

		if !test.Failed() {
			t.Error("values were not equal")
		}
	})
}

// captureOutput redirects output from stdout to a buffer.
func captureOutput(f func()) {
	old := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	f() // Output is generated in this function.
	w.Close()
	os.Stdout = old
}
