// step-by-step/internal/spentcalories/spentCalories.go

package spentcalories

import (
	"fmt"
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
	// Разделить строку на слайс строк
	parts := strings.Split(data, ",")

	// Проверить, чтобы длина слайса была равна 3
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неправильный формат данных: ожидается 3 поля, а найдено %d", len(parts))
	}

	// Преобразовать первый элемент слайса (количество шагов) в тип int
	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("невозможно преобразовать количество шагов в число: %w", err)
	}

	// Вид активности (второй элемент слайса)
	activity := strings.TrimSpace(parts[1])

	// Преобразовать третий элемент слайса в time.Duration
	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("невозможно преобразовать продолжительность в time.Duration: %w", err)
	}

	// Если всё прошло без ошибок, вернуть количество шагов, вид активности, продолжительность и nil
	return steps, activity, duration, nil
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
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return fmt.Sprintf("Ошибка при парсинге данных: %v", err)
	}

	// Проверка, чтобы количество шагов было больше 0
	if steps <= 0 {
		return "Количество шагов должно быть положительным."
	}

	// Определение типа тренировки
	var calories float64
	switch activity {
	case "Ходьба":
		calories = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories = RunningSpentCalories(steps, weight, duration)
	default:
		return "Неизвестный тип тренировки."
	}

	// Расчет дистанции
	dist := distance(steps)

	// Расчет средней скорости
	meanSpd := meanSpeed(steps, duration)

	// Формирование строки с информацией
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity, duration.Hours(), dist, meanSpd, calories)
}
