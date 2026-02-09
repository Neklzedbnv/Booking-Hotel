package domain


type RoomAvailability struct {
	RoomID       int64  `json:"room_id"`
	Room         *Room  `json:"room,omitempty"`
	IsAvailable  bool   `json:"is_available"`
	AvailableRooms int  `json:"available_rooms"`
}


type BookingStats struct {
	TotalBookings      int64   `json:"total_bookings"`
	ConfirmedBookings  int64   `json:"confirmed_bookings"`
	CanceledBookings   int64   `json:"canceled_bookings"`
	AverageStayDays    float64 `json:"average_stay_days"`
	TotalRevenue       float64 `json:"total_revenue"`
}
