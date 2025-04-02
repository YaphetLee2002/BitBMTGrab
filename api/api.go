package api

import (
	"badminton/config"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"badminton/models"
)

func makeRequest(url string, headers map[string]string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", config.ApiHost+url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	return body, nil
}

// GetSportEventsField 获取指定ID的场馆信息
func GetSportEventsField(id int, headers map[string]string) (map[string]*models.Venue, error) {
	// 打印请求信息
	fmt.Printf("正在获取场馆ID %d 的信息...\n", id)
	url := fmt.Sprintf("/api/sport_events/field/id/%d", id)

	body, err := makeRequest(url, headers)
	if err != nil {
		return nil, err
	}

	var response models.Response[map[string]*models.Venue]

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析场馆信息失败: %w", err)
	}

	if response.Code != 1 {
		return nil, fmt.Errorf("获取场馆信息失败: %s", response.Message)
	}

	// 打印场馆信息
	//fmt.Println("获取到的场馆信息:")
	//for key, venue := range response.Data {
	//	fmt.Printf("场地 %s:\n%s\n", key, venue)
	//}

	return response.Data, nil
}

// GetSportEventsHour 获取指定ID的时间段信息
func GetSportEventsHour(id int, headers map[string]string) ([]*models.SportEventsHour, error) {
	fmt.Printf("正在获取场馆ID %d 的时间段信息...\n", id)
	url := fmt.Sprintf("/api/sport_events/hour/id/%d", id)

	body, err := makeRequest(url, headers)
	if err != nil {
		return nil, err
	}

	var response models.Response[[]*models.SportEventsHour]

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析时间段信息失败: %w", err)
	}

	if response.Code != 1 {
		return nil, fmt.Errorf("获取时间段信息失败: %s", response.Message)
	}

	// 打印时间段信息
	//fmt.Println("获取到的时间段信息:")
	//for _, hour := range response.Data {
	//	fmt.Println(hour)
	//}

	return response.Data, nil
}

// GetSportScheduleBooked 获取指定日期的场地预订信息
func GetSportScheduleBooked(id int, day string, headers map[string]string) (*models.SportScheduleBooked, error) {
	fmt.Printf("正在获取场馆ID %d 在 %s 的预订信息...\n", id, day)
	url := fmt.Sprintf("/api/sport_schedule/booked/id/%d?day=%s", id, day)

	body, err := makeRequest(url, headers)
	if err != nil {
		return nil, err
	}

	var response models.SportScheduleBooked
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析场地预订信息失败: %w", err)
	}

	if response.Code != 1 {
		return nil, fmt.Errorf("获取场地预订信息失败: %s", response.Msg)
	}

	// 打印场地预订信息
	//fmt.Println("获取到的场地预订信息:")
	//for key, status := range response.Data {
	//	fmt.Printf("场地-时段 %s: %d\n", key, status)
	//}

	return &response, nil
}

// GetSportEventsPrice 获取指定场馆的价格信息
func GetSportEventsPrice(id int, week int, day string, headers map[string]string) (*models.SportEventsPrice, error) {
	fmt.Printf("正在获取场馆ID %d 在星期 %d, 日期 %s 的价格信息...\n", id, week, day)
	url := fmt.Sprintf("/api/sport_events/price/id/%d?week=%d&day=%s", id, week, day)

	headers["token"] = config.ApiToken
	body, err := makeRequest(url, headers)
	if err != nil {
		return nil, err
	}

	var response models.SportEventsPrice
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析价格信息失败: %w", err)
	}

	if response.Code != 1 {
		return nil, fmt.Errorf("获取价格信息失败: %s", response.Msg)
	}

	fmt.Println("获取到的价格信息:")
	for daytype, price := range response.Data {
		fmt.Printf("%s时段: 价格 %d元\n", daytype, price.Price)
	}

	return &response, nil
}

// GetOrderList 获取订单列表
func GetOrderList(orderType string, status string, orderID string, page int, limit int, headers map[string]string) (*models.OrderList, error) {
	fmt.Printf("正在获取订单列表，类型：%s，状态：%s，订单号：%s，页码：%d，每页数量：%d...\n", orderType, status, orderID, page, limit)
	url := fmt.Sprintf("/api/order/index?ordertype=%s&status=%s&orderid=%s&page=%d&limit=%d", orderType, status, orderID, page, limit)

	headers["token"] = config.ApiToken
	body, err := makeRequest(url, headers)
	if err != nil {
		return nil, err
	}

	var response models.OrderList
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析订单列表失败: %w", err)
	}

	if response.Code != 1 {
		return nil, fmt.Errorf("获取订单列表失败: %s", response.Msg)
	}

	fmt.Printf("获取到 %d 条未支付订单记录\n", response.Data.Total)
	for _, order := range response.Data.List {
		fmt.Printf("订单号: %s, 金额: %d元", order.OrderID, order.Amount)
	}
	return &response, nil
}

func makePostRequest(url string, data map[string]string, headers map[string]string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	formData := make([]string, 0)
	for key, value := range data {
		formData = append(formData, fmt.Sprintf("%s=%s", key, value))
	}
	formDataStr := strings.Join(formData, "&")

	req, err := http.NewRequest("POST", config.ApiHost+url, strings.NewReader(formDataStr))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	return body, nil
}

// CheckSportSchedule 检查订单是否可以下单
func CheckSportSchedule(id int, scene string, headers map[string]string) (*models.CheckOrderResponse, error) {
	fmt.Printf("正在检查场馆ID %d 的订单信息...\n", id)
	url := fmt.Sprintf("/api/sport_schedule/check/id/%d", id)

	data := map[string]string{
		"scene": scene,
	}

	headers["token"] = config.ApiToken
	body, err := makePostRequest(url, data, headers)
	if err != nil {
		return nil, err
	}

	var response models.CheckOrderResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析订单检查信息失败: %w", err)
	}

	if response.Code != 1 {
		return nil, fmt.Errorf("订单检查失败: %s", response.Msg)
	}

	fmt.Println("获取到的订单检查信息:")
	fmt.Printf("场馆名称: %s\n", response.Data.VenueName)
	fmt.Printf("总金额: %d元\n", response.Data.TotalAmount)
	fmt.Printf("支付有效期: %d分钟\n", response.Data.PayValidTime)

	return &response, nil
}

// SubmitOrder 提交订单
func SubmitOrder(sportEventsId int, money int, scene string, headers map[string]string) (*models.OrderResponse, error) {
	fmt.Printf("正在提交订单，场馆ID: %d, 金额: %d元...\n", sportEventsId, money)
	url := "/api/order/submit"

	data := map[string]string{
		"orderid":         "",
		"card_id":         "",
		"sport_events_id": fmt.Sprintf("%d", sportEventsId),
		"money":           fmt.Sprintf("%d", money),
		"ordertype":       "makeappointment",
		"paytype":         "bitpay",
		"scene":           scene,
		"openid":          config.OpenId,
	}

	headers["token"] = config.ApiToken
	body, err := makePostRequest(url, data, headers)
	if err != nil {
		return nil, err
	}

	var response models.OrderResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析订单提交响应失败: %w", err)
	}

	if response.Code != 1 {
		return nil, fmt.Errorf("订单提交失败: %s", response.Msg)
	}

	fmt.Println("订单提交成功")

	return &response, nil
}
