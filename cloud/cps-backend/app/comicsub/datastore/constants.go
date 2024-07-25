package datastore

const (
	Printing1stPrint  = 1
	Printing2ndPrint  = 2
	Printing3rdPrint  = 3
	Printing4thPrint  = 4
	Printing5thPrint  = 5
	Printing6thPrint  = 6
	Printing7thPrint  = 7
	Printing8thPrint  = 8
	Printing9thPrint  = 9
	Printing10thPrint = 10
	PrintingAshcan    = 100
	PrintingUnknown   = 101

	KeyIssueFirstAppearance = 2
	KeyIssueFirstCameo      = 3
	KeyIssueFirstTeamUp     = 4
	KeyIssueFirstSoloTitle  = 5
	KeyIssueDeath           = 6
	KeyIssueWedding         = 7
	KeyIssueBirth           = 8
	KeyIssueIconicCover     = 9
	KeyIssueOrigin          = 10
	KeyIssueOther           = 1
)

var KeyIssueMap = map[int]string{
	2:  "1st appearance of",
	3:  "1st cameo of",
	4:  "1st Team Up of",
	5:  "1st appearance in solo title",
	6:  "Death of",
	7:  "Wedding of",
	8:  "Birth of",
	9:  "Iconic cover by",
	10: "Origin of",
	1:  "Other",
}

var ServiceTypeMap = map[int8]string{
	1: "Pre-Screening Service",
	2: "CPS Pedigree Service",
	3: "CPS Capsule",
	4: "CPS Capsule Indie Mint Gem",
	5: "CPS Capsule Signature Collection",
	6: "CPS Capsule U-Grade",
}
