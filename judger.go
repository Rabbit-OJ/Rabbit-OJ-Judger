package judger

import (
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/compare"
	JudgerModels "github.com/Rabbit-OJ/Rabbit-OJ-Judger/models"
)

const (
	StatusOK = "OK"
)

func JudgeOneCase(testResult *JudgerModels.TestResult, stdout, rightStdout, compMode string) *JudgerModels.JudgeResult {
	result := &JudgerModels.JudgeResult{}

	if testResult.Status != StatusOK {
		result.Status = testResult.Status
	} else {
		isAC := false
		if compMode == "STDIN_F" {
			isAC, _ = compare.ModeStdinFloat64(stdout, rightStdout)
		} else if compMode == "STDIN_S" {
			isAC, _ = compare.ModeStdinString(stdout, rightStdout)
		} else {
			isAC = compare.ModeCMP(stdout, rightStdout)
		}

		if isAC {
			result.Status = "AC"
		} else {
			result.Status = "WA"
		}
	}

	result.TimeUsed, result.SpaceUsed = testResult.TimeUsed, testResult.SpaceUsed
	return result
}
