package constants

var SubmissionPublisherNames = map[int8]string{
	1:  "Other (Please specify)",
	2:  "DC Comics",
	3:  "Marvel Comics",
	4:  "Image",
	5:  "Studiocomix Press",
	6:  "Lucha",
	7:  "Boom! Studios",
	8:  "Dark Horse Comics",
	9:  "IDW Publishing",
	10: "Archie Comics",
	11: "MUCAHI",
	12: "Novedades Editores SA de CV",
	13: "Marvel Comics Mexico",
	14: "Panini",
	15: "Aardvark-Vanaheim",
	16: "Scattered Comics",
	17: "Xochitl",
	18: "CCBA",
	19: "Dell",
	20: "Goldkey",
	21: "Charlton Comics",
	22: "Harvey Comics",
	23: "Star Comics",
	24: "Now Comics",
	25: "Mirage Studios",
	26: "Eclipse Comics",
	27: "Malibu Comics",
}

const (
	SubmissionPublisherNameOther = 1
)

var SubmissionOverallLetterGrades = map[string]string{
	"pr": "Poor",
	"fr": "Fair",
	"gd": "Good",
	"vg": "Very good",
	"fn": "Fine",
	"vf": "Very Fine",
	"nm": "Near Mint",
	"PR": "Poor",
	"FR": "Fair",
	"GD": "Good",
	"VG": "Very good",
	"FN": "Fine",
	"VF": "Very Fine",
	"NM": "Near Mint",
}

var SubmissionPrimaryLabelDetails = map[int8]string{
	1: "Other",
	2: "Regular Edition",
	3: "Direct Edition",
	4: "Newsstand Edition",
	5: "Variant Cover",
	6: "Canadian Price Variant",
	7: "Facsimile",
	8: "Reprint",
}

var SubmissionKeyIssue = map[int8]string{
	1:  "Other",
	2:  "1st appearance of",
	3:  "1st cameo of",
	4:  "1st Team Up of",
	5:  "1st appearance in solo title",
	6:  "Death of",
	7:  "Wedding of",
	8:  "Birth of",
	9:  "Iconic cover by",
	10: "Origin of",
}
