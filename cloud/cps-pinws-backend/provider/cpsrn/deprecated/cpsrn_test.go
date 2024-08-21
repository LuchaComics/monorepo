package cpsrn

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/cpsrn"
)

// go test provider/cpsrn/cpsrn_test.go

func TestGenerateNumberForCPS_PINWSAdministrationCSVExport(t *testing.T) {
	// Create a test instance of cpsrnProvider
	p := cpsrn.NewProvider()

	// Create a CSV file to save the actual and expected values
	file, err := os.Create("results.csv")
	if err != nil {
		t.Fatalf("Failed to create CSV file: %v", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the CSV header
	header := []string{"Count", "Output"}
	if err := writer.Write(header); err != nil {
		t.Fatalf("Failed to write CSV header: %v", err)
	}

	for i := 0; i <= 8_998*3; i++ {
		actual, err := p.GenerateNumber(cpsrn.UserRoleRoot, 0, 1, int64(i))
		if err != nil {
			t.Errorf("Test case csc error: for count %d, got error %v", i, err)
			return
		}

		// Write the actual and expected values to the CSV file
		row := []string{fmt.Sprintf("%d", i), actual}
		if err := writer.Write(row); err != nil {
			t.Errorf("Failed to write CSV row: %v", err)
		}
	}
}

func TestGenerateNumberForCPS_PINWSAdministrationFor19999(t *testing.T) {
	// Create a test instance of cpsrnProvider
	p := cpsrn.NewProvider()

	for i := 0; i <= 8_998; i++ {
		expected := "788346-26649-1-xxxx"
		expected = strings.Replace(expected, "xxxx", fmt.Sprintf("%d", (1001+i)), 1)

		actual, err := p.GenerateNumber(cpsrn.UserRoleRoot, 0, 1, int64(i))
		if err != nil {
			t.Errorf("Test case 19999 error: expected %s, got error %v", expected, err)
			return
		}

		if expected != actual {
			t.Errorf("Test case 19999 failed: expected %s, got %s", expected, actual)
			return
		}
	}
}

func TestGenerateNumberForCPS_PINWSAdministrationFor29999(t *testing.T) {
	// Create a test instance of cpsrnProvider
	p := cpsrn.NewProvider()

	for i := 8_999; i <= 8_998*2; i++ {
		expected := "788346-26649-2-xxxx"
		expected = strings.Replace(expected, "xxxx", fmt.Sprintf("%d", (1000+i-8_998)), 1)

		actual, err := p.GenerateNumber(cpsrn.UserRoleRoot, 0, 1, int64(i))
		if err != nil {
			t.Errorf("Test case 29999 error: expected %s, got error %v", expected, err)
			return
		}

		if expected != actual {
			t.Errorf("Test case 29999 failed: expected %s, got %s for i of %v", expected, actual, i)
			return
		}
	}
}

func TestGenerateNumberForCPS_PINWSAdministrationFor39999(t *testing.T) {
	// Create a test instance of cpsrnProvider
	p := cpsrn.NewProvider()

	for i := 8_998*2 + 1; i <= 8_998*3; i++ {
		expected := "788346-26649-3-xxxx"
		expected = strings.Replace(expected, "xxxx", fmt.Sprintf("%d", (1000+i-8_998*2)), 1)

		actual, err := p.GenerateNumber(cpsrn.UserRoleRoot, 0, 1, int64(i))
		if err != nil {
			t.Errorf("Test case 39999 error: expected %s, got error %v", expected, err)
			return
		}

		if expected != actual {
			t.Errorf("Test case 39999 failed: expected %s, got %s for i of %v", expected, actual, i)
			return
		}
	}
}

func TestGenerateNumberForCPS_PINWSAdministration(t *testing.T) {
	// Create a test instance of cpsrnProvider
	p := cpsrn.NewProvider()

	// CASE 1: 1-1001

	e1 := "788346-26649-1-1001"
	r1, err := p.GenerateNumber(cpsrn.UserRoleRoot, 0, 1, 0)
	if err != nil {
		t.Errorf("Test case 1 error: expected %s, got error %v", e1, err)
	}
	if r1 != e1 {
		t.Errorf("Test case 1 failed: expected %s, got %s", e1, r1)
	}

	// CASE 2: 1-1002

	e2 := "788346-26649-1-1002"
	r2, err := p.GenerateNumber(cpsrn.UserRoleRoot, 0, 1, 1)
	if err != nil {
		t.Errorf("Test case 2 error: expected %s, got error %v", e2, err)
	}
	if r2 != e2 {
		t.Errorf("Test case 2 failed: expected %s, got %s", e2, r2)
	}

	// CASE 3: 1-9999

	e3 := "788346-26649-1-9999"
	r3, err := p.GenerateNumber(cpsrn.UserRoleRoot, 0, 1, 8_998)
	if err != nil {
		t.Errorf("Test case 3 error: expected %s, got error %v", e3, err)
	}
	if r3 != e3 {
		t.Errorf("Test case 3 failed: expected %s, got %s", e3, r3)
	}

	// CASE 4: 2-1001

	e4 := "788346-26649-2-1001"
	r4, err := p.GenerateNumber(cpsrn.UserRoleRoot, 0, 1, 8_998+1)
	if err != nil {
		t.Errorf("Test case 4 error: expected %s, got error %v", e4, err)
	}
	if r4 != e4 {
		t.Errorf("Test case 4 failed: expected %s, got %s", e4, r4)
	}

	// CASE 5: 2-1002

	e5 := "788346-26649-2-1002"
	r5, err := p.GenerateNumber(cpsrn.UserRoleRoot, 0, 1, 8_998+2)
	if err != nil {
		t.Errorf("Test case 5 error: expected %s, got error %v", e5, err)
	}
	if r5 != e5 {
		t.Errorf("Test case 5 failed: expected %s, got %s", e5, r5)
	}

	// // CASE 6: 2-9999
	//
	// e6 := "788346-26649-2-9999"
	// r6, err := p.GenerateNumber(cpsrn.UserRoleRoot, 0, 1, 8_998+8_998)
	// if err != nil {
	// 	t.Errorf("Test case 6 error: expected %s, got error %v", e6, err)
	// }
	// if r6 != e6 {
	// 	t.Errorf("Test case 6 failed: expected %s, got %s", e6, r6)
	// }

	// CASE 7: 3-1001

	e7 := "788346-26649-3-1001"
	r7, err := p.GenerateNumber(cpsrn.UserRoleRoot, 0, 1, 8_998+8_998+1)
	if err != nil {
		t.Errorf("Test case 7 error: expected %s, got error %v", e7, err)
	}
	if r7 != e7 {
		t.Errorf("Test case 7 failed: expected %s, got %s", e7, r7)
	}

	// CASE 8: 3-1002

	e8 := "788346-26649-3-1002"
	r8, err := p.GenerateNumber(cpsrn.UserRoleRoot, 0, 1, 8_998+8_998+2)
	if err != nil {
		t.Errorf("Test case 8 error: expected %s, got error %v", e8, err)
	}
	if r8 != e8 {
		t.Errorf("Test case 8 failed: expected %s, got %s", e8, r8)
	}

	// // CASE 9: 3-9999
	//
	// e9 := "788346-26649-3-9999"
	// r9, err := p.GenerateNumber(cpsrn.UserRoleRoot, 0, 1, 8_998+8_998+8_998)
	// if err != nil {
	// 	t.Errorf("Test case 9 error: expected %s, got error %v", e9, err)
	// }
	// if r9 != e9 {
	// 	t.Errorf("Test case 9 failed: expected %s, got %s", e9, r9)
	// }

	// CASE 10: 4-1001

	r10, err := p.GenerateNumber(cpsrn.UserRoleRoot, 0, 1, 8_998+8_998+8_998+8_998)
	if err == nil {
		t.Errorf("Test case 10: Expected error but got none and value of %s", r10)
	}

}

func TestGenerateNumberForSpecialCollection(t *testing.T) {
	// Create a test instance of cpsrnProvider
	p := cpsrn.NewProvider()

	// CASE 1: 1-1001

	e1 := "788346-26649-4-1001"
	r1, err := p.GenerateNumber(0, 1, 1, 0)
	if err != nil {
		t.Errorf("Test case special collection error: expected %s, got error %v", e1, err)
	}
	if r1 != e1 {
		t.Errorf("Test case special collection failed: expected %s, got %s", e1, r1)
	}
}
