package models

import "fmt"

// String 为Response实现String方法
func (r Response[T]) String() string {
	return fmt.Sprintf("Code: %s\nMessage: %s\nTime: %s\nData: %v", r.Code, r.Message, r.TimeString, r.Data)
}

// String 为SportEventsHour实现String方法
func (s SportEventsHour) String() string {
	return fmt.Sprintf("时间段ID: %d\n开始时间: %s\n结束时间: %s\n日期类型: %s",
		s.Id, s.BegintimeText, s.EndtimeText, s.Daytype)
}

// String 为Venue实现String方法
func (v Venue) String() string {
	return fmt.Sprintf("场馆ID: %d\n场馆名称: %s\n最大人数: %d", v.ID, v.Name, v.MaxNums)
}