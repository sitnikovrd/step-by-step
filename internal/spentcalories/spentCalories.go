// step-by-step/internal/spentcalories/spentCalories.go

package spentcalories

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                            = 0.65  // средняя длина шага.
	mInKm                              = 1000  // количество метров в километре.
	minInH                             = 60    // количество минут в часе.
	runningCaloriesMeanSpeedMultiplier = 18.0  // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0  // среднее количество сжигаемых калорий при беге.
	walkingCaloriesWeightMultiplier    = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier       = 0.029 // множитель роста.
)

// parseTraining парсит строку с данными тренировки
func parseTraining(data string) (int, string, time.Duration, error) {
	// Регулярные выражения для поиска количества шагов, вида активности и времени
	stepPattern := regexp.MustCompile(`^\d+`)
	typePattern := regexp.MustCompile(`,\s*([^,]+),\s*`)
	timePattern := regexp.MustCompile(`\d+h\d+m$`)

	// Найти количество шагов
	matchSteps := stepPattern.FindStringSubmatch(data)
	if matchSteps == nil {
		return 0, "", 0, fmt.Errorf("невозможно распознать количество шагов в строке: %q", data)
	}
	steps, err := strconv.Atoi(matchSteps[1])
	if err != nil {
		return 0, "", 0, fmt.Errorf("невозможно преобразовать количество шагов в число: %w", err)
	}

	// Найти вид активности
	matchType := typePattern.FindStringSubmatch(data)
	if matchType == nil {
		return 0, "", 0, fmt.Errorf("невозможно распознать вид активности в строке: %q", data)
	}

	// Найти время
	matchTime := timePattern.FindStringSubmatch(data)
	if matchTime == nil {
		return 0, "", 0, fmt.Errorf("невозможно распознать время в строке: %q", data)
	}

	// Парсинг времени в формате "XhYm"
	parts := strings.Split(matchTime[1], "h")
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("невозможно преобразовать часы в число: %w", err)
	}

	minutes := 0
	if len(parts) > 1 {
		minPart := strings.TrimSuffix(parts[1], "m")
		minutes, err = strconv.Atoi(minPart)
		if err != nil {
			return 0, "", 0, fmt.Errorf("невозможно преобразовать минуты в число: %w", err)
		}
	}

	// Преобразовать время в Duration
	duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute

	return steps, matchType[1], duration, nil
}

// distance возвращает дистанцию в километрах
func distance(steps int) float64 {
	return float64(steps) * lenStep / mInKm
}

// meanSpeed возвращает среднюю скорость
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration == 0 {
		return 0
	}
	dist := distance(steps)
	hours := duration.Hours()
	return dist / hours
}

// RunningSpentCalories возвращает калории, потраченные при беге
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	meanSpeed := meanSpeed(steps, duration)
	calories := ((runningCaloriesMeanSpeedMultiplier * meanSpeed) - runningCaloriesMeanSpeedShift) * weight
	return calories
}

// WalkingSpentCalories возвращает калории, потраченные при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	meanSpeed := meanSpeed(steps, duration)
	calories := ((walkingCaloriesWeightMultiplier * weight) +
		(meanSpeed*meanSpeed/height)*walkingSpeedHeightMultiplier) *
		float64(duration.Minutes()) / minInH
	return calories
}

// TrainingInfo формирует строку с информацией о тренировке
func TrainingInfo(data string, weight, height float64) string {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		return fmt.Sprintf("Ошибка при парсинге данных: %v", err)
	}

	// Проверка, чтобы количество шагов было больше 0
	if steps <= 0 {
		return ""
	}

	// Определение типа тренировки
	var activity string
	switch activityType {
	case "Walking":
		activity = "Ходьба"
	case "Running":
		activity = "Бег"
	default:
		return "Неизвестный тип тренировки"
	}

	// Расчет дистанции
	dist := distance(steps)

	// Расчет средней скорости
	meanSpd := meanSpeed(steps, duration)

	// Расчет калорий
	var calories float64
	switch activityType {
	case "Walking":
		calories = WalkingSpentCalories(steps, weight, height, duration)
	case "Running":
		calories = RunningSpentCalories(steps, weight, duration)
	default:
		return "Неизвестный тип тренировки"
	}

	// Формирование строки с информацией
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f ккал.",
		activity, duration.Hours(), dist, meanSpd, calories)
}
