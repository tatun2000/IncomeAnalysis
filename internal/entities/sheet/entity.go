package sheet

import (
	"fmt"
	"time"

	"telegrammbot.core/internal/errs"
)

type ReqType string

var (
	AddValueToCell   ReqType = "Добавить в"
	GetValueFromCell ReqType = "Получить"
)

const (
	Health int = iota + 1
	Clothes
	EssentialGoods
	SecondHandGoods
	Flat
	Cafe
	Internet
	MobileComm
	OtherSubs
	Transport
	Devices
	Travelling
)

func ConvertCategoryTypeToCell(catType int) (result string, err error) {
	switch catType {
	case Health:
		return "19", nil
	case Clothes:
		return "20", nil
	case EssentialGoods:
		return "21", nil
	case SecondHandGoods:
		return "28", nil
	case Flat:
		return "22", nil
	case Cafe:
		return "23", nil
	case Internet:
		return "24", nil
	case MobileComm:
		return "25", nil
	case OtherSubs:
		return "26", nil
	case Transport:
		return "27", nil
	case Devices:
		return "29", nil
	case Travelling:
		return "30", nil
	default:
		return result, fmt.Errorf("ConvertCategoryTypeToCell: %w", errs.ErrInvalidCategoryType)
	}
}

func GetActualDayCell() (dayCell string, err error) {
	day := time.Now().Day()
	switch day {
	case 1:
		return "G", nil
	case 2:
		return "H", nil
	case 3:
		return "I", nil
	case 4:
		return "J", nil
	case 5:
		return "K", nil
	case 6:
		return "L", nil
	case 7:
		return "M", nil
	case 8:
		return "N", nil
	case 9:
		return "O", nil
	case 10:
		return "P", nil
	case 11:
		return "Q", nil
	case 12:
		return "R", nil
	case 13:
		return "S", nil
	case 14:
		return "T", nil
	case 15:
		return "U", nil
	case 16:
		return "V", nil
	case 17:
		return "W", nil
	case 18:
		return "X", nil
	case 19:
		return "Y", nil
	case 20:
		return "Z", nil
	case 21:
		return "AA", nil
	case 22:
		return "AB", nil
	case 23:
		return "AC", nil
	case 24:
		return "AD", nil
	case 25:
		return "AE", nil
	case 26:
		return "AF", nil
	case 27:
		return "AG", nil
	case 28:
		return "AH", nil
	case 29:
		return "AI", nil
	case 30:
		return "AJ", nil
	case 31:
		return "AK", nil
	default:
		return dayCell, fmt.Errorf("getActualDayCell: %w", errs.ErrInvalidDay)
	}
}

func GetActualMonthSheet() (monthSheet string, err error) {
	month := time.Now().Month()
	switch month {
	case 1:
		return "Январь", nil
	case 2:
		return "Февраль", nil
	case 3:
		return "Март", nil
	case 4:
		return "Апрель", nil
	case 5:
		return "Май", nil
	case 6:
		return "Июнь", nil
	case 7:
		return "Июль", nil
	case 8:
		return "Август", nil
	case 9:
		return "Сентябрь", nil
	case 10:
		return "Октябрь", nil
	case 11:
		return "Ноябрь", nil
	case 12:
		return "Декабрь", nil
	default:
		return monthSheet, fmt.Errorf("getActualMonthSheet: %w", errs.ErrInvalidMonth)
	}
}
