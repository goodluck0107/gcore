package gtime_test

import (
	"gitee.com/monobytes/gcore/gutils/gtime"
	"testing"
)

func TestNow(t *testing.T) {
	t.Log(gtime.Now().Format(gtime.DateTime))
}

func TestToday(t *testing.T) {
	t.Log(gtime.Today())
}

func TestDay(t *testing.T) {
	t.Log(gtime.Day())
	t.Log(gtime.Day(-1))
	t.Log(gtime.Day(1))
}

func TestDayHead(t *testing.T) {
	t.Log(gtime.DayHead())
	t.Log(gtime.DayHead(-1))
	t.Log(gtime.DayHead(1))
}

func TestDayTail(t *testing.T) {
	t.Log(gtime.DayTail())
	t.Log(gtime.DayTail(-1))
	t.Log(gtime.DayTail(1))
}

func TestWeek(t *testing.T) {
	t.Log(gtime.Week())
	t.Log(gtime.Week(-1))
	t.Log(gtime.Week(1))
}

func TestWeekHead(t *testing.T) {
	t.Log(gtime.WeekHead())
	t.Log(gtime.WeekHead(-1))
	t.Log(gtime.WeekHead(1))
}

func TestWeekTail(t *testing.T) {
	t.Log(gtime.WeekTail())
	t.Log(gtime.WeekTail(-1))
	t.Log(gtime.WeekTail(1))
}

func TestMonth(t *testing.T) {
	for i := 0; i <= 100; i++ {
		t.Log(gtime.Month(0 - i))
	}
}

func TestMonthHead(t *testing.T) {
	for i := 0; i <= 100; i++ {
		t.Log(gtime.MonthHead(0 - i))
	}
}

func TestMonthTail(t *testing.T) {
	for i := 0; i <= 100; i++ {
		t.Log(gtime.MonthTail(0 - i))
	}
}
