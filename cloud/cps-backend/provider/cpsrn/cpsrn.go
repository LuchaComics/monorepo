package cpsrn

import (
	"fmt"
)

// Provider provides an interface for abstracting `CPS Registry Number`.
type Provider interface {
	Classify(specialCollection int8, serviceType int8, userRoleID int8) (string, error)
	Generate(specialCollection int8, serviceType int8, userRoleID int8, count int64) (string, error)
}

// cpsrnProvider is the structure to hold the base values to start the numbers.
// The most simple example is as follows: `788346-26649-1-1001` where `788346`
// is section A, `26649` is section B, `1` is section C and `1001` is section D.
type cpsrnProvider struct{}

// NewProvider function is a contructor that returns the default `CPS Registry Number` provider.
func NewProvider() Provider {
	return cpsrnProvider{}
}

const (
	UserRoleRoot = 1
)

// The following were copied from `app/comicsub/datastore/datastore.go`.
const (
	ServiceTypePreScreening                  = 1
	ServiceTypePedigree                      = 2
	ServiceTypeCPSCapsule                    = 3
	ServiceTypeCPSCapsuleIndieMintGem        = 4
	ServiceTypeCPSCapsuleSignatureCollection = 5
	ServiceTypeCPSCapsuleYouGrade            = 6
)

const (
	SectionA                          = 788346
	SectionB                          = 26649
	SectionD                          = 1
	MaxCountPerSectionD               = 9_999
	MaxUserSubmissionsCountPerSeciton = 9_998
)

// Generates the unique `CPS Registry Number` required for tracking submissions.
func (p cpsrnProvider) Generate(specialCollection int8, serviceType int8, userRoleID int8, count int64) (string, error) {
	// CASE 1 Special collections.
	switch specialCollection {
	case 1:
		return fmt.Sprintf("%d-%d-%d-%04d", SectionA, SectionB, 0, SectionD+count), nil
	case 2:
		return fmt.Sprintf("%d-%d-%d-%04d", SectionA, SectionB, 1, SectionD+count), nil
	case 3:
		return fmt.Sprintf("%d-%d-%d-%04d", SectionA, SectionB, 2, SectionD+count), nil
	case 4:
		return fmt.Sprintf("%d-%d-%d-%04d", SectionA, SectionB, 3, SectionD+count), nil
	case 5:
		return fmt.Sprintf("%d-%d-%d-%04d", SectionA, SectionB, 4, SectionD+count), nil
	}

	switch serviceType {
	case ServiceTypeCPSCapsuleYouGrade, ServiceTypeCPSCapsuleSignatureCollection, ServiceTypeCPSCapsuleIndieMintGem, ServiceTypeCPSCapsule:
		// STEP 1: Ensure the correct role was applied.
		switch serviceType {
		case ServiceTypeCPSCapsuleIndieMintGem, ServiceTypeCPSCapsule:
			if userRoleID != UserRoleRoot {
				return "", fmt.Errorf("You do not have sufficient access with role ID %v", userRoleID)
			}
		}

		// STEP 2: Calculate.

		if count <= MaxUserSubmissionsCountPerSeciton {
			return fmt.Sprintf("%d-%d-%d-%04d", SectionA, SectionB, 5, SectionD+count), nil
		}
		if count < MaxCountPerSectionD*2 {
			return fmt.Sprintf("%d-%d-%d-%04d", SectionA, SectionB, 6, SectionD+count-MaxCountPerSectionD), nil
		}
		if count < MaxCountPerSectionD*3 {
			return fmt.Sprintf("%d-%d-%d-%04d", SectionA, SectionB, 7, SectionD+count-(MaxCountPerSectionD*2)), nil
		}
		if count < MaxCountPerSectionD*4 {
			return fmt.Sprintf("%d-%d-%d-%04d", SectionA, SectionB, 8, SectionD+count-(MaxCountPerSectionD*3)), nil
		}
		if count < MaxCountPerSectionD*5 {
			return fmt.Sprintf("%d-%d-%d-%04d", SectionA, SectionB, 9, SectionD+count-(MaxCountPerSectionD*4)), nil
		}
		if count < MaxCountPerSectionD*6 {
			return fmt.Sprintf("%d-%d-%d-%04d", SectionA, SectionB, 10, SectionD+count-(MaxCountPerSectionD*5)), nil
		}
		return "", fmt.Errorf("Out of numbers with current count at %v", count)

	case ServiceTypePedigree:
		if count <= MaxUserSubmissionsCountPerSeciton {
			return fmt.Sprintf("%d-%d-%d-%04d", SectionA, SectionB, 11, SectionD+count), nil
		}
		if count < MaxCountPerSectionD*2 {
			return fmt.Sprintf("%d-%d-%d-%04d", SectionA, SectionB, 12, SectionD+count-MaxCountPerSectionD), nil
		}
		if count < MaxCountPerSectionD*3 {
			return fmt.Sprintf("%d-%d-%d-%04d", SectionA, SectionB, 13, SectionD+count-(MaxCountPerSectionD*2)), nil
		}
		return "", fmt.Errorf("Out of numbers with current count at %v", count)

	case ServiceTypePreScreening:
		if count <= MaxUserSubmissionsCountPerSeciton {
			return fmt.Sprintf("%d-%d-%d-%04d", SectionA, SectionB, 14, SectionD+count), nil
		}
		return "", fmt.Errorf("Out of numbers with current count at %v", count)

	default:
		return "", fmt.Errorf("unsupported service type %v", serviceType)
	}
}

func (p cpsrnProvider) Classify(specialCollection int8, serviceType int8, userRoleID int8) (string, error) {
	// CASE 1 Special collections.
	switch specialCollection {
	case 1:
		return "0-0001", nil
	case 2:
		return "1-0001", nil
	case 3:
		return "2-0001", nil
	case 4:
		return "3-0001", nil
	case 5:
		return "4-0001", nil
	}

	switch serviceType {
	case ServiceTypeCPSCapsuleYouGrade, ServiceTypeCPSCapsuleSignatureCollection, ServiceTypeCPSCapsuleIndieMintGem, ServiceTypeCPSCapsule:
		// Ensure the correct role was applied.
		switch serviceType {
		case ServiceTypeCPSCapsuleIndieMintGem, ServiceTypeCPSCapsule:
			if userRoleID != UserRoleRoot {
				return "", fmt.Errorf("You do not have sufficient access with role ID %v", userRoleID)
			}
		}
		return "5-0001", nil

	case ServiceTypePedigree:
		return "11-0001", nil

	case ServiceTypePreScreening:
		return "14-0001", nil

	default:
		return "", fmt.Errorf("unsupported service type %v", serviceType)
	}
}
