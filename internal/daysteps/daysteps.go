package daysteps

import (
	"errors"
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
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("Длина слайса не равна 2")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, errors.New("Длина слайса не равна 2")
	}

	if steps <= 0 {
		return 0, 0, errors.New("Количество шагов должно быть положительным")
	}

	parsedTime, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, errors.New("неверный формат продолжительности")
	}
	return steps, parsedTime, nil
}

func DayActionInfo(data string, weight, height float64) string {

	step, times, err := parsePackage(data)
	if err != nil {
		fmt.Errorf("Ошибка в парсинге строки: %w", err)
		return ""
	}

	if step <= 0 {
		fmt.Errorf("Количество шагов меньше 0: %d", step)
		return ""
	}

	dist := (stepLength * float64(step)) / mInKm
	calories, err := spentcalories.WalkingSpentCalories(step, weight, height, times)
	if err != nil {
		fmt.Errorf("Ошибка в расчете калорий: %w", err)
		return ""
	}

	result := fmt.Sprintf(
		"Количество шагов: %d\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли: %.2f ккал.",
		step, dist, calories)

	return result
}
