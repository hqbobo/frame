package utils

import "time"

//GetNextMinute 获取下一分
func GetNextMinute(now time.Time) int64 {
	return now.Unix() - int64(now.Second()) + 60
}

//GetNextHour 获取下一小时
func GetNextHour(now time.Time) int64 {
	return now.Unix() - int64(now.Minute()*60) - int64(now.Second()) + int64(60*60)
}

//GetNextDay 获取下一天
func GetNextDay(now time.Time) int64 {
	return now.Unix() - int64(now.Hour()*3600) - int64(now.Minute()*60) - int64(now.Second()) + int64(3600*24)
}

//GetNextWeek 获取下一周
func GetNextWeek(now time.Time) int64 {
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset+7).Unix()
}

//GetNextMonth 获取下一月
func GetNextMonth(now time.Time) int64 {
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0).Unix()
}
