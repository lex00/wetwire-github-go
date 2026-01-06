package dependabot

// Schedule defines when Dependabot checks for updates.
type Schedule struct {
	// Interval is how often to check for updates.
	// Values: "daily", "weekly", "monthly".
	Interval string `yaml:"interval"`

	// Day is the day of week for weekly schedules.
	// Values: "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday".
	Day string `yaml:"day,omitempty"`

	// Time is the time of day to check (hh:mm format, UTC by default).
	Time string `yaml:"time,omitempty"`

	// Timezone is the IANA timezone identifier.
	Timezone string `yaml:"timezone,omitempty"`
}
