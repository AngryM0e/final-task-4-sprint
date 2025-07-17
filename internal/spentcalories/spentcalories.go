package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	dataSlice := strings.Split(data, ",")
	if len(dataSlice) != 3 {
		err := fmt.Errorf("некорректное количество полученных данных")
		return 0, "0 шагов", 0, err
	}
	
	steps, err := strconv.Atoi(dataSlice[0])
	if err != nil {
		return 0, "0 шагов", 0, err
	}

	activity := dataSlice[1]
	if activity != "Ходьба" && activity != "Бег" {
		err := fmt.Errorf("неизвестный тип тренировки: %s", activity)
		return 0, "", 0, err
	}

	duration, err := time.ParseDuration(dataSlice[2])
	if err != nil {
		return 0, "0 шагов", 0, err
	}
	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceKm := (float64(steps) * stepLength) / mInKm
	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	
	hours := duration.Hours()

	speed := dist / hours

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	
	if err != nil {
		log.Println("Ошибка парсинга данных тренировки:", err)
		return "", fmt.Errorf("ошибка парсинга данных: %v", err)
	}

	var dist, speed, calories float64
	var calcErr error

	switch activityType {
	case "бег":
		dist = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, calcErr = RunningSpentCalories(steps, weight, height, duration)
	case "ходьба":
		dist = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, calcErr = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activityType)
	}

	if calcErr != nil {
		log.Println("Ошибка расчета показателей тренировки:", calcErr)
		return "", fmt.Errorf("ошибка расчета показателей: %v", calcErr)
	}

	durationHours := fmt.Sprintf("%.2f", duration.Hours())
	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %s ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activityType,
		durationHours,
		dist,
		speed,
		calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной")
	}
	
	speed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	calories := (weight * speed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной")
	}

	speed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	baseCalories := (weight * speed * durationInMinutes) / minInH


	calories := baseCalories * walkingCaloriesCoefficient

return calories, nil
}
