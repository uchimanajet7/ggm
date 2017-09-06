package utils

import "time"

var jstLocation *time.Location

// GetUserConfig is get the now epoch time in milliseconds
func GetNowEpoch() int64 {
	return getEpochFromTime(time.Now())
}

func getJstTime(time time.Time) time.Time {
	jstLoc := getJstLocation()

	return time.In(jstLoc)
}

func getDefaultLastDate() time.Time {
	// 15 minutes ago
	return time.Now().Add(-15 * time.Minute)
}

func getTimeFromEpoch(epoch int64) time.Time {
	// epoch ms to time
	return time.Unix(0, epoch*int64(time.Millisecond))
}

func getEpochFromTime(t time.Time) int64 {
	// time to epoch ms
	return t.UnixNano() / 1000000
}

func getJstLocation() *time.Location {
	if jstLocation != nil {
		return jstLocation
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		loc = time.FixedZone("Asia/Tokyo", 9*60*60)
	}

	jstLocation = loc
	time.Local = jstLocation

	return jstLocation
}
