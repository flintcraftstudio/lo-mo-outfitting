package admin

// TripLabel returns a short human-readable label for a trip type value.
var TripLabel = map[string]string{
	"full_day_single": "Full Day",
	"half_day_single": "Half Day",
	"early_season":    "Early Season",
	"winter":          "Winter",
	"multiple_boats":  "Multiple Boats",
	"heroes":          "Heroes Rate",
}

// ExperienceLabel returns a human-readable label for an experience level value.
var ExperienceLabel = map[string]string{
	"never":       "First time",
	"some":        "Some experience",
	"comfortable": "Comfortable",
	"advanced":    "Advanced",
}

// LodgingLabel returns a human-readable label for a lodging value.
var LodgingLabel = map[string]string{
	"craig":       "Craig",
	"wolf_creek":  "Wolf Creek",
	"helena":      "Helena",
	"great_falls": "Great Falls",
	"not_sure":    "Not sure yet",
	"other":       "Other",
}

// StatusLabel returns a human-readable label for a booking status.
var StatusLabel = map[string]string{
	"new":          "New",
	"contacted":    "Contacted",
	"deposit_sent": "Deposit Sent",
	"confirmed":    "Confirmed",
	"complete":     "Complete",
	"cancelled":    "Cancelled",
}

// StatusColor returns a Tailwind color class for a status dot.
func StatusColor(status string) string {
	switch status {
	case "new":
		return "bg-teal"
	case "contacted":
		return "bg-blue-500"
	case "deposit_sent":
		return "bg-yellow-500"
	case "confirmed":
		return "bg-green-500"
	case "complete":
		return "bg-stone"
	case "cancelled":
		return "bg-red-500"
	default:
		return "bg-stone"
	}
}

// GetLabel is a helper to safely look up a map value with a fallback.
func GetLabel(m map[string]string, key string) string {
	if label, ok := m[key]; ok {
		return label
	}
	return key
}
