package spentcalories

import (
	"fmt"
	// "log"
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
			return 0, "", 0, fmt.Errorf("неверный формат данных")
	}

	steps, err := strconv.Atoi(strings.TrimSpace(dataSlice[0]))
	if err != nil {
			return 0, "", 0, fmt.Errorf("неверное количество шагов")
	}
	if steps <= 0 {
			return 0, "", 0, fmt.Errorf("количество шагов должно быть положительным")
	}

	activity := strings.TrimSpace(dataSlice[1])
	if activity != "Ходьба" && activity != "Бег" {
			return 0, "", 0, fmt.Errorf("неизвестный тип тренировки")
	}

	duration, err := time.ParseDuration(strings.TrimSpace(dataSlice[2]))
	if err != nil {
			return 0, "", 0, fmt.Errorf("неверная продолжительность")
	}
	if duration <= 0 {
			return 0, "", 0, fmt.Errorf("продолжительность должна быть положительной")
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

	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
			return "", err
	}

	var calories float64
	var calcErr error

	switch activityType {
	case "Ходьба":
			calories, calcErr = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
			calories, calcErr = RunningSpentCalories(steps, weight, height, duration)
	default:
			return "", fmt.Errorf("неизвестный тип тренировки")
	}

	if calcErr != nil {
			return "", fmt.Errorf("ошибка расчета калорий: %w", err)
	}

	return fmt.Sprintf(
			"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			activityType,
			duration.Hours(),
			distance(steps, height),
			meanSpeed(steps, height, duration),
			calories,
	), nil
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
