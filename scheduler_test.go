package judger

import (
	"context"
	"fmt"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/config"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/docker"
	JudgerModels "github.com/Rabbit-OJ/Rabbit-OJ-Judger/models"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/protobuf"
	"sync"
	"testing"
)

var (
	alreadyInit = false
	initMu      sync.Mutex
)

func MockGetStorage(tid uint32, version string) ([]*JudgerModels.TestCaseType, error) {
	if tid == uint32(1) {
		testCases := []*JudgerModels.TestCaseType{
			{
				Id:         1,
				Stdin:      []byte("1 2"),
				Stdout:     []byte("3"),
				//StdinPath:  "/Users/yangziyue/Downloads/case/1.in",
				//StdoutPath: "/Users/yangziyue/Downloads/case/1.out",
				StdinPath:  "/home/case/1.in",
				StdoutPath: "/home/case/1.out",
			},
			{
				Id:         2,
				Stdin:      []byte("3 5"),
				Stdout:     []byte("8"),
				//StdinPath:  "/Users/yangziyue/Downloads/case/2.in",
				//StdoutPath: "/Users/yangziyue/Downloads/case/2.out",
				StdinPath:  "/home/case/2.in",
				StdoutPath: "/home/case/2.out",
			},
		}

		return testCases, nil
	}
	return make([]*JudgerModels.TestCaseType, 0), nil
}

func initJudger() {
	initMu.Lock()
	defer initMu.Unlock()
	if alreadyInit {
		return
	}

	alreadyInit = true
	ctx, _ := context.WithCancel(context.Background())
	cfg := &JudgerModels.JudgerConfigType{
		Kafka: JudgerModels.KafkaConfig{
			Brokers: []string{
				"localhost:9092",
			},
		},
		Rpc: "",
		AutoRemove: JudgerModels.AutoRemoveType{
			Containers: true,
			Files:      true,
		},
		Concurrent: JudgerModels.ConcurrentType{
			Judge: 2,
		},
		LocalImages: []string{
			"alpine_tester:latest",
		},
		Languages: []JudgerModels.LanguageType{
			{
				ID:      "cpp17",
				Name:    "C++17",
				Enabled: true,
				Args: JudgerModels.CompileInfo{
					BuildArgs: []string{
						"g++",
						"-std=c++17",
						"/home/code.cpp",
						"-Wall",
						"-lm",
						"-fno-asm",
						"--static",
						"-O2",
						"-o",
						"/home/code.o",
					},
					Source:      "/home/code.cpp",
					NoBuild:     false,
					BuildTarget: "/home/code.o",
					BuildImage:  "gcc:10.2.0",
					Constraints: JudgerModels.Constraints{
						CPU:          1000000000,
						Memory:       1073741824,
						BuildTimeout: 120,
						RunTimeout:   120,
					},
					RunArgs:     []string{"/home/code.o"},
					RunArgsJSON: "[\"/home/code.o\"]",
					RunImage:    "alpine_tester:latest",
				},
			},
		},
		Extensions: JudgerModels.ExtensionsType{
			HostBind: false,
			AutoPull: true,
			CheckJudge: JudgerModels.CheckJudgeType{
				Enabled:  false,
				Interval: 0,
				Requeue:  false,
			},
			Expire: JudgerModels.ExpireType{
				Enabled:  false,
				Interval: 0,
			},
		},
	}

	InitJudger(ctx, cfg, MockGetStorage, true, false, "Judge")

	OnJudgeResponse = append(OnJudgeResponse, func(sid uint32, isContest bool, judgeResult []*JudgerModels.JudgeResult) {
		fmt.Println(sid, isContest, judgeResult)
	})
}

func TestInitJudger(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%+v \n", err)
			t.Fail()
		}
	}()

	initJudger()
	needImages := docker.GetNeedImages()
	for image, need := range needImages {
		if need {
			fmt.Printf("Need to image %s", image)
			t.Fail()
		}
	}
}

func testJudgeHelper(code []byte) (string, []*protobuf.JudgeCaseResult, error) {
	initJudger()

	config.Global.Extensions.HostBind = true
	status1, result1, err1 := Scheduler(&protobuf.JudgeRequest{
		Sid:        1,
		Tid:        1,
		Version:    "1",
		Language:   "cpp17",
		TimeLimit:  1000,
		SpaceLimit: 128,
		CompMode:   "STDIN_S",
		Code:       code,
		Time:       0,
		IsContest:  false,
	})

	config.Global.Extensions.HostBind = false
	status2, result2, err2 := Scheduler(&protobuf.JudgeRequest{
		Sid:        1,
		Tid:        1,
		Version:    "1",
		Language:   "cpp17",
		TimeLimit:  1000,
		SpaceLimit: 128,
		CompMode:   "STDIN_S",
		Code:       code,
		Time:       0,
		IsContest:  false,
	})

	b1, b2 := err1 == nil, err2 == nil
	if (b1 && !b2) || (!b1 && b2) {
		panic("Inconsistency error state")
	}

	if status1 != status2 {
		panic("Inconsistency state state")
	}

	if len(result1) != len(result2) {
		panic("Inconsistency result length")
	}

	totalLength := len(result1)
	for i := 0; i < totalLength; i++ {
		if (result1[i] == nil && result2[i] != nil) || (result1[i] != nil && result2[i] == nil) {
			panic("Inconsistency test case result")
		}

		if result1[i] == nil || result2[i] == nil {
			continue
		}

		if result1[i].Status != result2[i].Status {
			panic("Inconsistency test case result")
		}
	}

	return status1, result1, err1
}

func TestShouldEmitCE(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%+v \n", err)
			t.Fail()
		}
	}()

	code := []byte("#include <iostream> \n" +
		"int mian() { \n" +
		"    return 0; \n" +
		"}")

	status, _, _ := testJudgeHelper(code)
	if status != "CE" {
		t.Fail()
	}
}

func TestShouldEmitRE(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%+v \n", err)
			t.Fail()
		}
	}()

	code := []byte("#include <iostream> \n" +
		"int main() { \n" +
		"    exit(9); \n" +
		"    return 0; \n" +
		"}")
	status, judgeResult, _ := testJudgeHelper(code)

	if status != "OK" {
		fmt.Println("[Should Emit RE] Status NOT OK")
		t.Fail()
	}
	for _, result := range judgeResult {
		if result.Status != "RE" {
			fmt.Println("[Should Emit RE] Some Case Status NOT RE", result)
			t.Fail()
		}
	}
}

func TestShouldEmitTLE(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%+v \n", err)
			t.Fail()
		}
	}()

	initJudger()

	code := []byte("#include <iostream> \n" +
		"int main() { \n" +
		"    while (1) {} \n" +
		"    return 0; \n" +
		"}")
	status, judgeResult, _ := testJudgeHelper(code)

	if status != "OK" {
		fmt.Println("[Should Emit TLE] Status NOT OK")
		t.Fail()
	}
	for _, result := range judgeResult {
		if result.Status != "TLE" {
			fmt.Println("[Should Emit TLE] Some Case Status NOT TLE", result)
			t.Fail()
		}
	}
}

func TestShouldEmitAC(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%+v \n", err)
			t.Fail()
		}
	}()

	code := []byte("#include <iostream> \n" +
		"int main() { \n" +
		"    int x, y; \n" +
		"    std::cin >> x >> y; \n" +
		"    std::cout << (x + y) << std::endl; \n" +
		"    return 0; \n" +
		"}")
	status, judgeResult, _ := testJudgeHelper(code)

	if status != "OK" {
		fmt.Println("[Should Emit AC] Status NOT OK")
		t.Fail()
	}
	for _, result := range judgeResult {
		if result.Status != "AC" {
			fmt.Println("[Should Emit AC] Some Case Status NOT AC", result)
			t.Fail()
		}
	}
}

func TestShouldEmitWA(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%+v \n", err)
			t.Fail()
		}
	}()

	code := []byte("#include <iostream> \n" +
		"int main() { \n" +
		"    int x, y; \n" +
		"    std::cin >> x >> y; \n" +
		"    std::cout << (x * y) << std::endl; \n" +
		"    return 0; \n" +
		"}")
	status, judgeResult, _ := testJudgeHelper(code)

	if status != "OK" {
		fmt.Println("[Should Emit WA] Status NOT OK")
		t.Fail()
	}
	for _, result := range judgeResult {
		if result.Status != "WA" {
			fmt.Println("[Should Emit WA] Some Case Status NOT WA", result)
			t.Fail()
		}
	}
}
