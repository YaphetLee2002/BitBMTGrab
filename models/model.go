package models

type Response[T any] struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
	TimeString string `json:"time"`
}

type SportEventsHour struct {
	Id            int    `json:"id"`
	SportEventsId int    `json:"sport_events_id"`
	Begintime     int    `json:"begintime"`
	Endtime       int    `json:"endtime"`
	Createtime    int    `json:"createtime"`
	BegintimeText string `json:"begintime_text"`
	EndtimeText   string `json:"endtime_text"`
	Daytype       string `json:"daytype"`
}

type Venue struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	MaxNums int    `json:"maxnums"`
}

type SportScheduleBooked struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Time string         `json:"time"`
	Data map[string]int `json:"data"`
}

type PriceInfo struct {
	Daytype   string `json:"daytype"`
	Price     int    `json:"price"`
	HalfPrice int    `json:"half_price"`
}

type SportEventsPrice struct {
	Code int                  `json:"code"`
	Msg  string               `json:"msg"`
	Time string               `json:"time"`
	Data map[string]PriceInfo `json:"data"`
}

type SceneItem struct {
	Name          string `json:"name"`
	Amount        int    `json:"amount"`
	Periods       int    `json:"periods"`
	Fields        int    `json:"fields"`
	SportEventsId string `json:"sport_events_id"`
	Weekday       string `json:"weekday"`
	Begintime     int64  `json:"begintime"`
	Endtime       int64  `json:"endtime"`
	Beginhourtime int64  `json:"beginhourtime"`
	Endhourtime   int64  `json:"endhourtime"`
}

type CheckOrderResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Time string `json:"time"`
	Data struct {
		SceneList       []SceneItem `json:"scene_list"`
		TotalAmount     int         `json:"total_amount"`
		ExtraAmount     *int        `json:"extra_amount"`
		SportEventsId   string      `json:"sport_events_id"`
		SportEventsName string      `json:"sport_events_name"`
		VenueId         int         `json:"venue_id"`
		VenueName       string      `json:"venue_name"`
		PayValidTime    int         `json:"pay_valid_time"`
		PayId           int         `json:"pay_id"`
	} `json:"data"`
}

type OrderResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Time string `json:"time"`
}

type OrderSceneItem struct {
	Name          string `json:"name"`
	Amount        int    `json:"amount"`
	Periods       int    `json:"periods"`
	Fields        int    `json:"fields"`
	SportEventsID string `json:"sport_events_id"`
	Weekday       string `json:"weekday"`
	Begintime     int64  `json:"begintime"`
	Endtime       int64  `json:"endtime"`
	Beginhourtime int64  `json:"beginhourtime"`
	Endhourtime   int64  `json:"endhourtime"`
}

type OrderScene struct {
	Day    string         `json:"day"`
	Fields map[string][]int `json:"fields"`
}

type OrderConfig struct {
	SceneList      []OrderSceneItem `json:"scene_list"`
	TotalAmount    int              `json:"total_amount"`
	ExtraAmount    *int             `json:"extra_amount"`
	SportEventsID  string           `json:"sport_events_id"`
	SportEventsName string          `json:"sport_events_name"`
	VenueID        int              `json:"venue_id"`
	VenueName      string           `json:"venue_name"`
	PayValidTime   int              `json:"pay_valid_time"`
	PayID          int              `json:"pay_id"`
	Scene          []OrderScene     `json:"scene"`
}

type OrderItem struct {
	ID             int          `json:"id"`
	OrderID        string       `json:"orderid"`
	OrderType      string       `json:"ordertype"`
	UserID         int          `json:"user_id"`
	Amount         int          `json:"amount"`
	PayAmount      int          `json:"payamount"`
	PayType        string       `json:"paytype"`
	PayTime        *string      `json:"paytime"`
	IP             string       `json:"ip"`
	UserAgent      string       `json:"useragent"`
	Memo           string       `json:"memo"`
	CreateTime     int64        `json:"createtime"`
	UpdateTime     int64        `json:"updatetime"`
	DeleteTime     *int64       `json:"deletetime"`
	Status         string       `json:"status"`
	ConfigJSON     string       `json:"configjson"`
	VenueID        int          `json:"venue_id"`
	SportEventsID  int          `json:"sport_events_id"`
	PayID          int          `json:"pay_id"`
	TransactionID  *string      `json:"transaction_id"`
	OpenID         string       `json:"openid"`
	SportEventsName string      `json:"sportevents_name"`
	VenueName      string       `json:"venue_name"`
	Config         OrderConfig  `json:"config"`
}

type OrderList struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Time string `json:"time"`
	Data struct {
		Total int         `json:"total"`
		List  []OrderItem `json:"list"`
	} `json:"data"`
}
