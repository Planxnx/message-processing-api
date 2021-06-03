package scheduleconstant

type ScheduleType string

const (
	ScheduleType_HOURLY ScheduleType = "HOURLY"
	ScheduleType_DAILY  ScheduleType = "DAILY"
	ScheduleType_WEEKLY ScheduleType = "WEEKLY"
)

type WeekDay string

const (
	WeekDay_SUNDAY    WeekDay = "Sunday"
	WeekDay_MONDAY    WeekDay = "Monday"
	WeekDay_TUESDAY   WeekDay = "Tuesday"
	WeekDay_WEDNESDAY WeekDay = "Wednesday"
	WeekDay_THURSDAY  WeekDay = "Thursday"
	WeekDay_FRIDAY    WeekDay = "Friday"
	WeekDay_SATURDAY  WeekDay = "Saturday"
)
