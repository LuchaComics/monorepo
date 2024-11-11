package cpsrn

import (
	"log"
	"testing"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/cpsrn"
)

// go test provider/cpsrn/cpsrn_test.go

type GenerateTestCase struct {
	SpecialCollection int8
	ServiceType       int8
	UserRoleID        int8
	Count             int64
	Expected          string
}

func TestGenerateForSpecialCollections(t *testing.T) {
	// Create a test instance of cpsrnProvider
	p := cpsrn.NewProvider()

	data := []*GenerateTestCase{
		{1, 0, 0, 0, "788346-26649-0-0001"},
		{1, 0, 0, 999, "788346-26649-0-1000"},
		{1, 0, 0, 9998, "788346-26649-0-9999"},
		{2, 0, 0, 9998, "788346-26649-1-9999"},
		{3, 0, 0, 0, "788346-26649-2-0001"},
		{4, 0, 0, 0, "788346-26649-3-0001"},
		{5, 0, 0, 0, "788346-26649-4-0001"},
	}

	for _, d := range data {
		r, err := p.Generate(d.SpecialCollection, d.ServiceType, d.UserRoleID, d.Count)
		if err != nil {
			t.Errorf("Test case special collection error: expected %s, got error %v", d.Expected, err)
		}
		if r != d.Expected {
			t.Errorf("Test case special collection failed: expected %s, got %s", d.Expected, r)
		}
	}
}

func TestGenerateForUGrade(t *testing.T) {
	// Create a test instance of cpsrnProvider
	p := cpsrn.NewProvider()

	data := []*GenerateTestCase{
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleYouGrade, 1, 0, "788346-26649-5-0001"},
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleYouGrade, 1, 999, "788346-26649-5-1000"},
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleYouGrade, 1, 9998, "788346-26649-5-9999"},
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleYouGrade, 1, 9999, "788346-26649-6-0001"},
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleYouGrade, 1, 9999 + 1000, "788346-26649-6-1001"},
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleYouGrade, 1, 9999 + 9998, "788346-26649-6-9999"},
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleYouGrade, 1, 9999 + 9999, "788346-26649-7-0001"},
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleYouGrade, 1, 9999 + 9999 + 9998, "788346-26649-7-9999"},
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleYouGrade, 1, 9999 + 9999 + 9999, "788346-26649-8-0001"},
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleYouGrade, 1, 9999 + 9999 + 9999 + 9998, "788346-26649-8-9999"},
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleYouGrade, 1, 9999 + 9999 + 9999 + 9999, "788346-26649-9-0001"},
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleYouGrade, 1, 9999 + 9999 + 9999 + 9999 + 9998, "788346-26649-9-9999"},
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleYouGrade, 1, 9999 + 9999 + 9999 + 9999 + 9999, "788346-26649-10-0001"},
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleSignatureCollection, 1, 0, "788346-26649-5-0001"},
	}

	for _, d := range data {
		r, err := p.Generate(d.SpecialCollection, d.ServiceType, d.UserRoleID, d.Count)
		if err != nil {
			t.Errorf("Test case special collection error: expected %s, got error %v", d.Expected, err)
			return
		}
		if r != d.Expected {
			log.Println("d.SpecialCollection --->", d.SpecialCollection, "\nd.ServiceType --->", d.ServiceType, "\nd.UserRoleID --->", d.UserRoleID, "\nd.Count --->", d.Count)
			t.Errorf("Test case special collection failed: expected %s, got %s", d.Expected, r)
			return
		}
	}
}

func TestGenerateForIndieMintGem(t *testing.T) {
	// Create a test instance of cpsrnProvider
	p := cpsrn.NewProvider()

	// Success cases.
	data := []*GenerateTestCase{
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleIndieMintGem, cpsrn.UserRoleRoot, 0, "788346-26649-5-0001"},
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleIndieMintGem, cpsrn.UserRoleRoot, 9998, "788346-26649-5-9999"},
	}

	for _, d := range data {
		r, err := p.Generate(d.SpecialCollection, d.ServiceType, d.UserRoleID, d.Count)
		if err != nil {
			t.Errorf("Test indie mint gem error: expected %s, got error %v", d.Expected, err)
			return
		}
		if r != d.Expected {
			t.Errorf("Test indie mint gem failed: expected %s, got %s", d.Expected, r)
			return
		}
	}

	// Confirm non-admin cannot create.
	errDatum := &GenerateTestCase{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleIndieMintGem, 0, 0, "788346-26649-5-0001"}
	_, err := p.Generate(errDatum.SpecialCollection, errDatum.ServiceType, errDatum.UserRoleID, errDatum.Count)
	if err == nil {
		t.Errorf("Test indie mint gem did not get error: expected error for %s", errDatum.Expected)
		return
	}

	// Out of bounds case.
	errDatum2 := &GenerateTestCase{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleYouGrade, 1, 9999 + 9999 + 9999 + 9999 + 9999 + 9999, "788346-26649-10-0001"}
	_, err2 := p.Generate(errDatum2.SpecialCollection, errDatum2.ServiceType, errDatum2.UserRoleID, errDatum2.Count)
	if err2 == nil {
		t.Errorf("Test indie mint gem did not get error: expected error for %s", errDatum2.Expected)
		return
	}
}

func TestGenerateForPedigree(t *testing.T) {
	// Create a test instance of cpsrnProvider
	p := cpsrn.NewProvider()

	// Success cases.
	data := []*GenerateTestCase{
		{0, cpsrn.ServiceTypePedigree, cpsrn.UserRoleRoot, 0, "788346-26649-11-0001"},
		{0, cpsrn.ServiceTypePedigree, cpsrn.UserRoleRoot, 9998, "788346-26649-11-9999"},
		{0, cpsrn.ServiceTypePedigree, cpsrn.UserRoleRoot, 9999, "788346-26649-12-0001"},
		{0, cpsrn.ServiceTypePedigree, cpsrn.UserRoleRoot, 9999 + 9998, "788346-26649-12-9999"},
		{0, cpsrn.ServiceTypePedigree, cpsrn.UserRoleRoot, 9999 + 9999, "788346-26649-13-0001"},
		{0, cpsrn.ServiceTypePedigree, cpsrn.UserRoleRoot, 9999 + 9999 + 9998, "788346-26649-13-9999"},
	}

	for _, d := range data {
		r, err := p.Generate(d.SpecialCollection, d.ServiceType, d.UserRoleID, d.Count)
		if err != nil {
			t.Errorf("Test pedigree error: expected %s, got error %v", d.Expected, err)
			return
		}
		if r != d.Expected {
			t.Errorf("Test pedigree failed: expected %s, got %s", d.Expected, r)
			return
		}
	}

	// Out of bounds case.
	errDatum := &GenerateTestCase{0, cpsrn.ServiceTypePedigree, cpsrn.UserRoleRoot, 9999 + 9999 + 9999, "788346-26649-13-9999"}
	_, err := p.Generate(errDatum.SpecialCollection, errDatum.ServiceType, errDatum.UserRoleID, errDatum.Count)
	if err == nil {
		t.Errorf("Test pedigree did not get error: expected error for %s", errDatum.Expected)
		return
	}
}

func TestGenerateForPreScreening(t *testing.T) {
	// Create a test instance of cpsrnProvider
	p := cpsrn.NewProvider()

	// Success cases.
	data := []*GenerateTestCase{
		{0, cpsrn.ServiceTypePreScreening, cpsrn.UserRoleRoot, 0, "788346-26649-14-0001"},
		{0, cpsrn.ServiceTypePreScreening, cpsrn.UserRoleRoot, 9998, "788346-26649-14-9999"},
	}

	for _, d := range data {
		r, err := p.Generate(d.SpecialCollection, d.ServiceType, d.UserRoleID, d.Count)
		if err != nil {
			t.Errorf("Test pedigree error: expected %s, got error %v", d.Expected, err)
			return
		}
		if r != d.Expected {
			t.Errorf("Test pedigree failed: expected %s, got %s", d.Expected, r)
			return
		}
	}

	// Out of bounds case.
	errDatum := &GenerateTestCase{0, cpsrn.ServiceTypePedigree, cpsrn.UserRoleRoot, 9999 + 9999 + 9999, "788346-26649-14-9999"}
	_, err := p.Generate(errDatum.SpecialCollection, errDatum.ServiceType, errDatum.UserRoleID, errDatum.Count)
	if err == nil {
		t.Errorf("Test pedigree did not get error: expected error for %s", errDatum.Expected)
		return
	}
}

type ClassifyTestCase struct {
	SpecialCollection int8
	ServiceType       int8
	UserRoleID        int8
	Expected          string
}

func TestClassify(t *testing.T) {
	// Create a test instance of cpsrnProvider
	p := cpsrn.NewProvider()

	data := []*ClassifyTestCase{
		{1, 0, 0, "0-0001"},
		{2, 0, 0, "1-0001"},
		{3, 0, 0, "2-0001"},
		{4, 0, 0, "3-0001"},
		{5, 0, 0, "4-0001"},
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleYouGrade, cpsrn.UserRoleRoot, "5-0001"},
		{0, cpsrn.ServiceTypeCPS_NFTSTORECapsuleSignatureCollection, cpsrn.UserRoleRoot, "5-0001"},
		{0, cpsrn.ServiceTypePedigree, cpsrn.UserRoleRoot, "11-0001"},
		{0, cpsrn.ServiceTypePreScreening, cpsrn.UserRoleRoot, "14-0001"},
	}

	for _, d := range data {
		r, err := p.Classify(d.SpecialCollection, d.ServiceType, d.UserRoleID)
		if err != nil {
			t.Errorf("Test case special collection error: expected %v, got error %v", d.Expected, err)
		}
		if r != d.Expected {
			t.Errorf("Test case special collection failed: expected %v, got %v", d.Expected, r)
		}
	}
}
