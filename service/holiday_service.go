package service

import (
	"log"
	"time"
)

type HolidayServiceIf interface {
	IsTodayBlackFriday() bool
}

type holidayService struct{}

func NewHolidayService() HolidayServiceIf {
	return &holidayService{}
}

func (h *holidayService) IsTodayBlackFriday() bool {
	dtBlackFriday, err := findBlackFridayDay()
	if err != nil {
		log.Println(err)
	}
	dtNow := time.Now()

	return dtBlackFriday.Day() == dtNow.Day() && dtBlackFriday.Month() == dtNow.Month()
}

//como estamos usando um repositorio em memoria, retornamos de maneira fixa a data, deveria ser oriundo de uma base de feriados
func findBlackFridayDay() (time.Time, error) {
	layout := "02-01-2006"
	return time.Parse(layout, "26-11-2021")
}
