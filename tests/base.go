package tests

import (
	"log"

	"api-tests-template/internal/utils"
)

// SetupSuite действия перед запуском сьюта с тестами
func SetupSuite() {
	log.Println("Init environment variables")
	utils.LoadEnv()
}

// TearDownSuite действия после запуска сьюта с тестами
func TearDownSuite() {
	log.Println("Tear down suite")
}

// Precondition выводит текст в консоль с меткой [Precondition]
func Precondition(text string) {
	utils.LogWithLabelAndTimeStamp("Precondition", text)
}

// Step writes a test action to the console with a Step label.
func Step(text string) {
	utils.LogWithLabelAndTimeStamp("Step", text)
}

// Result writes a verified outcome to the console with a Result label.
func Result(text string) {
	utils.LogWithLabelAndTimeStamp("Result", text)
}
