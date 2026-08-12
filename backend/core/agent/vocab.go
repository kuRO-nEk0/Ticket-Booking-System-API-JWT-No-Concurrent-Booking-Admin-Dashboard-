package agent

var Categories = []string{
	"music", "comedy", "theatre", "sports",
	"workshop", "conference", "film", "festival",
}

var Tags = []string{
	"acoustic", "jazz", "indie", "electronic", "classical",
	"standing", "seated", "outdoor", "late_night",
	"family_friendly", "all_ages", "18_plus",
}

var MoodTags = map[string][]string{
	"chill":     {"acoustic", "jazz", "seated"},
	"big_night": {"electronic", "standing", "late_night"},
	"family":    {"family_friendly", "all_ages"},
	"cultural":  {"theatre", "classical"},
}
