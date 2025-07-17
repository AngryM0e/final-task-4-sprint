package daysteps

import (
	"fmt"
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
		return 0, 0, fmt.Errorf("недостаточное количество данных")
	}
	steps, err := strconv.Atoi(parsedData[0])
	if err != nil {
		err = fmt.Errorf("ошибка преобразования строки в int")
		return 0, 0, err
	}
	if steps <= 0 {
		err = fmt.Errorf("недостаточное количество шагов для вычисления")
		return 0, 0, err
	}
	duration, err := time.ParseDuration(parsedData[1])
	if err != nil {
		err = fmt.Errorf("ошибка преобразования строки в time.Duration")
		return 0, 0, err
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		return ""
	}
	if steps <= 0 {
		return ""
	}
	distance := float64(steps) * stepLength
	kilometres := distance / mInKm
	lostCalories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return "Ошибка расчёт калорий"
	}
	dayResults := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %f км.\nВы сожгли %f ккал.", steps, kilometres, lostCalories)
	return dayResults
}
