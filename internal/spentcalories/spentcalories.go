package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"log"
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
		return 0, "", 0, fmt.Errorf("в строке недостаточно данных")
	}

	steps, err := strconv.Atoi(parts[0])

	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	activity := parts[1] 

	duration, err := time.ParseDuration(parts[2])

	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования: %w", err)
	}

	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	return steps, activity, duration, nil

}

func distance(steps int, height float64) float64 {

	stepLen := height * stepLengthCoefficient 

	return float64(steps) * stepLen / float64(mInKm)
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)

	return dist / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	
	steps, activity, duration, err := parseTraining(data) 

	if err != nil {
		log.Println(err)
		return "", err
	}

	var dist, avgSpeed, calories float64

	switch activity {
		case "Ходьба":
			dist = distance(steps, height)
			avgSpeed = meanSpeed(steps, height, duration)
			calories, err = WalkingSpentCalories(steps, weight, height, duration)
			if err != nil { 
				fmt.Println(err)
				return "", err
			}
		case "Бег":
			dist = distance(steps, height)
			avgSpeed = meanSpeed(steps, height, duration)
			calories, err = RunningSpentCalories(steps, weight, height, duration)
			if err != nil { 
				fmt.Println(err)
				return "", err
			}
		default:
			return "", fmt.Errorf("неизвестный тип тренировки")
	}

	durationInHours := duration.Hours()

	result := fmt.Sprintf(
		"Тип тренировки: %s\n"+
		"Длительность: %.2f ч.\n"+
		"Дистанция: %.2f км.\n"+
		"Скорость: %.2f км/ч\n"+
		"Сожгли калорий: %.2f\n",
		activity, durationInHours, dist, avgSpeed, calories,
	)

	return result, nil

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	avgSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	return (weight * avgSpeed * durationInMinutes) / float64(minInH), nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	avgSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * avgSpeed * durationInMinutes) / float64(minInH)

	return calories * walkingCaloriesCoefficient, nil

}
