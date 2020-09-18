package judger

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/config"
	JudgerModel "github.com/Rabbit-OJ/Rabbit-OJ-Judger/models"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/protobuf"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type CollectedStdout struct {
	Stdout      string
	RightStdout string
}

func Scheduler(request *protobuf.JudgeRequest) (string, []*protobuf.JudgeCaseResult, error) {
	sid := request.Sid

	fmt.Printf("========START JUDGE(%d)======== \n", sid)
	fmt.Printf("(%d) [Scheduler] Received judge request \n", sid)

	startSchedule := time.Now()
	defer func() {
		fmt.Printf("(%d) [Scheduler] total cost : %d ms \n", sid, time.Since(startSchedule).Milliseconds())
	}()

	// initialize files
	currentPath, err := utils.SubmissionGenerateDirWithMkdir(sid)
	if err != nil {
		return "Internal Error", nil, err
	}

	defer func() {
		fmt.Printf("(%d) [Scheduler] Cleaning files \n", sid)
		if config.Global.AutoRemove.Files {
			_ = os.RemoveAll(currentPath)
		}
	}()

	outputPath, err := utils.JudgeGenerateOutputDirWithMkdir(currentPath)
	if err != nil {
		//return false, err
		return "Internal Error", nil, err
	}

	codePath := fmt.Sprintf("%s/", currentPath)

	compileInfo, ok := config.CompileObject[request.Language]
	if !ok {
		return "Internal Error", nil, errors.New("language doesn't support")
	}

	fmt.Printf("(%d) [Scheduler] Init test cases \n", sid)
	// get case
	testCases, err := StorageGetFunc(request.Tid, request.Version)
	if err != nil {
		return "Internal Error", nil, err
	}
	testCaseCount := len(testCases)

	var buildProduction []byte
	if !compileInfo.NoBuild {
		// compile
		fmt.Printf("(%d) [Scheduler] Start Compile \n", sid)
		if buildProduction, err = Compiler(
			sid,
			codePath,
			request.Code,
			&compileInfo,
		); err != nil {
			fmt.Printf("(%d) [Scheduler] CE %+v \n", sid, err)
			//CallbackAllError("CE", sid, request.IsContest, testCaseCount)
			return "CE", make([]*protobuf.JudgeCaseResult, testCaseCount), err
		}

		fmt.Printf("(%d) [Scheduler] Compile OK \n", sid)
	}

	// run
	fmt.Printf("(%d) [Scheduler] Start Runner \n", sid)
	var runnerCollectedStdout map[string][]byte
	if runnerCollectedStdout, err = Runner(
		sid,
		codePath,
		&compileInfo,
		testCases,
		strconv.FormatUint(uint64(request.TimeLimit), 10),
		strconv.FormatUint(uint64(request.SpaceLimit), 10),
		outputPath,
		request.Code,
		buildProduction); err != nil {

		fmt.Printf("(%d) [Scheduler] RE %+v \n", sid, err)
		//CallbackAllError("RE", sid, request.IsContest, testCaseCount)
		return "RE", make([]*protobuf.JudgeCaseResult, testCaseCount), err
	}
	fmt.Printf("(%d) [Scheduler] Runner OK \n", sid)

	fmt.Printf("(%d) [Scheduler] Reading result \n", sid)
	jsonFileByte, err := ioutil.ReadFile(filepath.Join(codePath, "result.json"))
	if err != nil {
		//CallbackAllError("RE", sid, request.IsContest, testCaseCount)
		return "RE", make([]*protobuf.JudgeCaseResult, testCaseCount), err
	}

	var testResultArr []JudgerModel.TestResult
	if err := json.Unmarshal(jsonFileByte, &testResultArr); err != nil || testResultArr == nil {
		//CallbackAllError("RE", sid, request.IsContest, testCaseCount)
		return "RE", make([]*protobuf.JudgeCaseResult, testCaseCount), err
	}

	// collect std::out
	fmt.Printf("(%d) [Scheduler] Collecting stdout \n", sid)
	allStdin := make([]CollectedStdout, testCaseCount)
	for i := 1; i <= testCaseCount; i++ {
		allStdin[i-1].RightStdout = string(testCases[i-1].Stdout)
	}

	// optimize this: avoid writing, reading file in the disk (performance optimization)
	if runnerCollectedStdout == nil {
		for i := 1; i <= testCaseCount; i++ {
			path := fmt.Sprintf("%s/%d.out", outputPath, i)

			stdoutByte, err := ioutil.ReadFile(path)
			if err != nil {
				allStdin[i-1].Stdout = ""
			} else {
				allStdin[i-1].Stdout = string(stdoutByte)
			}
		}
	} else {
		for i := 1; i <= testCaseCount; i++ {
			key := fmt.Sprintf("%d.out", i)

			if stdoutByte, ok := runnerCollectedStdout[key]; ok {
				allStdin[i-1].Stdout = string(stdoutByte)
			} else {
				allStdin[i-1].Stdout = ""
			}
		}
	}

	// judge std::out
	fmt.Printf("(%d) [Scheduler] Judging stdout \n", sid)
	resultList := make([]*protobuf.JudgeCaseResult, testCaseCount)

	for index, item := range allStdin {
		testResult := &testResultArr[index]
		resultList[index] = &protobuf.JudgeCaseResult{}

		judgeResult := JudgeOneCase(testResult, item.Stdout, item.RightStdout, request.CompMode)
		resultList[index].Status = judgeResult.Status
		resultList[index].SpaceUsed = judgeResult.SpaceUsed
		resultList[index].TimeUsed = judgeResult.TimeUsed
	}
	//CallbackSuccess(sid, request.IsContest, resultList)

	fmt.Printf("(%d) [Scheduler] Finish \n", sid)
	return "OK", resultList, nil
}
