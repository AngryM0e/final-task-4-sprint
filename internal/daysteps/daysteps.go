package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parsedData := strings.Split(data, ",")
	if len(parsedData) != 2 {
		return 0, 0, fmt.Errorf("требуется 2 параметра: шаги и продолжительность")
	}
	steps, err := strconv.Atoi(parsedData[0])
	if err != nil {
		err = fmt.Errorf("ошибка преобразования строки в int: %v", err)
		return 0, 0, err
	}
	if steps <= 0 {
		err = fmt.Errorf("количество шагов должно быть положительным")
		return 0, 0, err
	}
	duration, err := time.ParseDuration(parsedData[1])
	if err != nil {
		err = fmt.Errorf("ошибка преобразования строки в time.Duration")
		return 0, 0, err
	}

	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть положительной")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	if weight <= 0 || height <= 0 {
		log.Printf("Некорректные параметры: вес %.2f, рост %.2f", weight, height)
		return ""
	}

	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("Ошибка парсинга данных: %v", err)
        return ""
	}
	if steps <= 0 {
		log.Printf("Некорректное количество шагов: %d", steps)
        return ""
	}
	distance := float64(steps) * stepLength
	kilometres := distance / mInKm
	lostCalories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("Ошибка расчета калорий: %v", err)
        return ""
	}
	dayResults := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, kilometres, lostCalories)
	return dayResults
}
