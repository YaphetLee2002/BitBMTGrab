package main

import (
	"badminton/api"
	"badminton/models"
	"fmt"
	"sort"
	"time"
)

// formatAvailableVenues 格式化并输出场地可用信息
func formatAvailableVenues(venues map[string]*models.Venue, hours []*models.SportEventsHour, booked *models.SportScheduleBooked) {
	var venueKeys []string
	for key := range venues {
		venueKeys = append(venueKeys, key)
	}
	sort.Strings(venueKeys)

	timeSlots := make(map[string]map[string]bool)
	for _, hour := range hours {
		timeSlot := fmt.Sprintf("%s-%s", hour.BegintimeText, hour.EndtimeText)
		timeSlots[timeSlot] = make(map[string]bool)
		for _, venueKey := range venueKeys {
			// 构造场地-时段键
			bookingKey := fmt.Sprintf("%s-%d", venueKey, hour.Id)
			// 检查是否被预订（0表示未预订）
			timeSlots[timeSlot][venueKey] = booked.Data[bookingKey] == 0
		}
	}

	fmt.Println("\n场地预订情况:")
	var timeSlotKeys []string
	for timeSlot := range timeSlots {
		timeSlotKeys = append(timeSlotKeys, timeSlot)
	}
	sort.Strings(timeSlotKeys)

	fmt.Printf("%-15s", "时间段")
	for _, venueKey := range venueKeys {
		fmt.Printf("%-8s", venues[venueKey].Name)
	}
	fmt.Println()

	for _, timeSlot := range timeSlotKeys {
		fmt.Printf("%-15s", timeSlot)
		for _, venueKey := range venueKeys {
			if timeSlots[timeSlot][venueKey] {
				fmt.Printf("%-8s", "可用")
			} else {
				fmt.Printf("%-8s", "已订")
			}
		}
		fmt.Println()
	}
}

// getWeekday 根据日期字符串获取星期几（1-7）
func getWeekday(dateStr string) (int, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, fmt.Errorf("日期格式错误: %v", err)
	}

	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return weekday, nil
}

// findVenueIDByName 根据场地名称查找场地ID
func findVenueIDByName(venues map[string]*models.Venue, name string) (string, error) {
	for id, venue := range venues {
		if venue.Name == name {
			return id, nil
		}
	}
	return "", fmt.Errorf("未找到名为 %s 的场地", name)
}

// findHourIDByTime 根据时间（小时）查找对应的时间段ID
func findHourIDByTime(hours []*models.SportEventsHour, hour int) (int, error) {
	targetTime := fmt.Sprintf("%02d:00", hour)

	for _, h := range hours {
		if h.BegintimeText == targetTime {
			return h.Id, nil
		}
	}
	return 0, fmt.Errorf("未找到开始时间为 %s 的时间段", targetTime)
}

// waitUntilNextAvailableTime 检查当前时间是否在允许下单的时间范围内，如果不在则等待到下一个可用时间
func waitUntilNextAvailableTime() {
	for {
		now := time.Now()
		currentTime := now.Format("15:04:05")

		if currentTime >= "07:00:00" && currentTime <= "23:30:00" {
			return
		}

		nextAvailable := now
		if currentTime > "23:30:00" {
			nextAvailable = now.Add(24 * time.Hour)
		}
		nextAvailable = time.Date(
			nextAvailable.Year(),
			nextAvailable.Month(),
			nextAvailable.Day(),
			7, 0, 0, 0,
			nextAvailable.Location(),
		)

		waitDuration := nextAvailable.Sub(now)

		fmt.Printf("\r当前时间: %s, 等待时间: %s",
			currentTime,
			waitDuration.Round(time.Second),
		)

		time.Sleep(time.Second)
	}
}

func main() {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
	}

	fmt.Print("请输入场馆ID（25-良乡，51-中关村，54-中关村早晚场）: ")
	var venueID int
	fmt.Scanln(&venueID)

	venues, err := api.GetSportEventsField(venueID, headers)
	if err != nil {
		fmt.Printf("获取场馆信息失败: %v\n", err)
		return
	}

	hours, err := api.GetSportEventsHour(venueID, headers)
	if err != nil {
		fmt.Printf("获取时间段信息失败: %v\n", err)
		return
	}

	fmt.Print("请输入预订日期（格式：YYYY-MM-DD）: ")
	var dateStr string
	fmt.Scanln(&dateStr)

	booked, err := api.GetSportScheduleBooked(venueID, dateStr, headers)
	if err != nil {
		fmt.Printf("获取场地预订信息失败: %v\n", err)
		return
	}

	formatAvailableVenues(venues, hours, booked)

	week, err := getWeekday(dateStr)
	if err != nil {
		fmt.Printf("计算星期失败: %v\n", err)
		return
	}

	_, err = api.GetSportEventsPrice(venueID, week, dateStr, headers)
	if err != nil {
		fmt.Printf("获取价格信息失败: %v\n", err)
		return
	}

	fmt.Print("请输入场地名称: ")
	var venueName string
	fmt.Scanln(&venueName)

	venueFieldID, err := findVenueIDByName(venues, venueName)
	if err != nil {
		fmt.Printf("查找场地失败: %v\n", err)
		return
	}

	fmt.Print("请输入预订时间（小时，例如9表示09:00开始的时段）: ")
	var hour int
	fmt.Scanln(&hour)

	hourID, err := findHourIDByTime(hours, hour)
	if err != nil {
		fmt.Printf("查找时间段失败: %v\n", err)
		return
	}

	scene := fmt.Sprintf("[{\"day\":\"%s\",\"fields\":{\"%s\":[%d]}}]", dateStr, venueFieldID, hourID)

	check, err := api.CheckSportSchedule(venueID, scene, headers)
	if err != nil {
		fmt.Printf("检查订单失败: %v\n", err)
		return
	}

	fmt.Println("\n等待可下单时间...")
	waitUntilNextAvailableTime()

	for {
		orderList, err := api.GetOrderList("makeappointment", "created", "", 1, 20, headers)
		if err != nil {
			fmt.Printf("获取订单列表失败: %v\n", err)
			return
		}

		for _, order := range orderList.Data.List {
			if order.SportEventsID == venueID {
				fmt.Printf("已检测到订单 %d，正在等待支付...\n", order.ID)
				return
			}
		}

		_, err = api.SubmitOrder(venueID, check.Data.TotalAmount, scene, headers)
		if err != nil {
			fmt.Printf("提交订单失败: %v\n", err)
		}
		time.Sleep(5 * time.Second)
	}

}
