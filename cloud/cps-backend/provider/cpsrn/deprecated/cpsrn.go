package cpsrn

import (
	"fmt"
)

// Provider provides an interface for abstracting `CPS Registry Number`.
type Provider interface {
	GenerateNumber(userRoleID int8, specialCollection int8, serviceType int8, count int64) (string, error)
}

// cpsrnProvider is the structure to hold the base values to start the numbers.
// The most simple example is as follows: `788346-26649-1-1001` where `788346`
// is section A, `26649` is section B, `1` is section C and `1001` is section D.
type cpsrnProvider struct {
	sectionA int64
	sectionB int64
	sectionD int64
}

// NewProvider function is a contructor that returns the default `CPS Registry Number` provider.
func NewProvider() Provider {
	return cpsrnProvider{
		sectionA: 788346,
		sectionB: 26649,
		sectionD: 1001,
	}
}

const (
	UserRoleRoot                             = 1
	ServiceTypePreScreening                  = 1
	ServiceTypePedigree                      = 2
	ServiceTypeCPSCapsule                    = 3
	ServiceTypeCPSCapsuleIndieMintGem        = 4
	ServiceTypeCPSCapsuleSignatureCollection = 5
	ServiceTypeCPSCapsuleYouGrade            = 6
)

const (
	MaxCountPerSectionD               = 9_999
	MaxUserSubmissionsCountPerSeciton = 8_998
)

// Generates the unique `CPS Registry Number` required for tracking submissions.
func (p cpsrnProvider) GenerateNumber(userRoleID int8, specialCollection int8, serviceType int8, count int64) (string, error) {
	// CASE 1: CPS Administration.
	if userRoleID == UserRoleRoot {
		if count <= MaxUserSubmissionsCountPerSeciton {
			return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 1, p.sectionD+count), nil
		}
		if count <= MaxUserSubmissionsCountPerSeciton*2 {
			return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 2, p.sectionD+count-MaxUserSubmissionsCountPerSeciton-1), nil
		}
		if count <= MaxUserSubmissionsCountPerSeciton*3 {
			return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 3, p.sectionD+count-MaxUserSubmissionsCountPerSeciton*2-1), nil
		}

		return "", fmt.Errorf("Out of numbers for CPS org with current count at %v", count)
	}

	// CASE 2: Special collections.
	if specialCollection == 1 {
		return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 4, p.sectionD+count), nil
	}

	// CASE 3: Service Type
	switch serviceType {
	case ServiceTypeCPSCapsuleYouGrade: // 788346-26649-5-1001
		if count <= MaxUserSubmissionsCountPerSeciton {
			return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 5, p.sectionD+count), nil
		}
		if count <= MaxUserSubmissionsCountPerSeciton*2 {
			return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 6, p.sectionD+count-MaxUserSubmissionsCountPerSeciton-1), nil
		}
		if count <= MaxUserSubmissionsCountPerSeciton*3 {
			return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 7, p.sectionD+count-MaxUserSubmissionsCountPerSeciton*2-1), nil
		}
		if count <= MaxUserSubmissionsCountPerSeciton*4 {
			return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 8, p.sectionD+count-MaxUserSubmissionsCountPerSeciton*3-2), nil
		}
		if count <= MaxUserSubmissionsCountPerSeciton*5 {
			return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 9, p.sectionD+count-MaxUserSubmissionsCountPerSeciton*4-3), nil
		}
		return "", fmt.Errorf("Out of numbers with current count at %v", count)

	case ServiceTypePedigree: // 788346-26649-10-1001
		if count <= MaxUserSubmissionsCountPerSeciton {
			return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 10, p.sectionD+count), nil
		}
		if count <= MaxUserSubmissionsCountPerSeciton*2 {
			return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 11, p.sectionD+count-MaxUserSubmissionsCountPerSeciton-1), nil
		}
		if count <= MaxUserSubmissionsCountPerSeciton*3 {
			return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 12, p.sectionD+count-MaxUserSubmissionsCountPerSeciton*2-1), nil
		}
		if count <= MaxUserSubmissionsCountPerSeciton*4 {
			return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 13, p.sectionD+count-MaxUserSubmissionsCountPerSeciton*3-2), nil
		}
		return "", fmt.Errorf("Out of numbers with current count at %v", count)

	case ServiceTypePreScreening: // 788346-26649-14-1001
		if count <= MaxUserSubmissionsCountPerSeciton {
			return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 14, p.sectionD+count), nil
		}
		if count <= MaxUserSubmissionsCountPerSeciton*2 {
			return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 15, p.sectionD+count-MaxUserSubmissionsCountPerSeciton-1), nil
		}
		if count <= MaxUserSubmissionsCountPerSeciton*3 {
			return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 16, p.sectionD+count-MaxUserSubmissionsCountPerSeciton*2-1), nil
		}
		if count <= MaxUserSubmissionsCountPerSeciton*4 {
			return fmt.Sprintf("%d-%d-%d-%d", p.sectionA, p.sectionB, 17, p.sectionD+count-MaxUserSubmissionsCountPerSeciton*3-2), nil
		}
		return "", fmt.Errorf("Out of numbers with current count at %v", count)

	default:
		return "", fmt.Errorf("unsupported service type %v", serviceType)
	}
}
