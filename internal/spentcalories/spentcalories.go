package spentcalories

import (
	"errors"
	"fmt"
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
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("Длина слайса не равна 3")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, errors.New("Количество шагов должно быть больше 0")
	}

	workout := parts[1]
	if workout != "Бег" && workout != "Ходьба" {
		return 0, "", 0, errors.New("неизвестный тип тренировки")
	}

	parsedTime, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, errors.New("неверный формат продолжительности")
	}

	return steps, workout, parsedTime, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	dist := (stepLength * float64(steps)) / float64(mInKm)
	return dist
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	averageSpeed := dist / duration.Hours()

	return averageSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, workout, duration, err := parseTraining(data)
	if err != nil {
		return "", errors.New("Ошибка при выполнении парсинга")
	}

	var dist, averageSpeed, calories float64

	switch workout {
	case "Ходьба":
		dist = distance(steps, height)
		averageSpeed = meanSpeed(steps, height, duration)
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("ошибка расчета калорий: %w", err)
		}
	case "Бег":
		dist = distance(steps, height)
		averageSpeed = meanSpeed(steps, height, duration)
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("ошибка расчета калорий: %w", err)
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f",
		workout, duration, dist, averageSpeed, calories)
	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("Входные параметры некорректны. Кол-во шагов, вес, рост и продолжительность бега должны быть больше 0.")
	}

	averageSpeed := meanSpeed(steps, weight, duration)

	durationInMinutes := duration.Minutes()

	calories := (weight * averageSpeed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("Входные параметры некорректны. Кол-во шагов, вес, рост и продолжительность бега должны быть больше 0.")
	}

	averageSpeed := meanSpeed(steps, weight, duration)

	durationInMinutes := duration.Minutes()

	calories := ((weight * averageSpeed * durationInMinutes) / minInH) * walkingCaloriesCoefficient

	return calories, nil
}
