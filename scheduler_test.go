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
				Id:     1,
				Stdin:  []byte("1 2"),
				Stdout: []byte("3"),
				//StdinPath:  "/Users/yangziyue/Downloads/case/1.in",
				//StdoutPath: "/Users/yangziyue/Downloads/case/1.out",
				StdinPath:  "/home/case/1.in",
				StdoutPath: "/home/case/1.out",
			},
			{
				Id:     2,
				Stdin:  []byte("3 5"),
				Stdout: []byte("8"),
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
		BuildImages: []string{
			"alpine_tester:latest",
			"python_tester:latest",
			"java_tester:latest",
		},
		Languages: []JudgerModels.LanguageType{
			{
				ID:      "cpp17",
				Name:    "C++/17",
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
			{
				ID:      "rust",
				Name:    "Rust/1.46",
				Enabled: true,
				Args: JudgerModels.CompileInfo{
					BuildArgs: []string{
						"rustc",
						"-O",
						"/home/code.rs",
						"-o",
						"/home/code.o",
						"--target",
						"x86_64-unknown-linux-musl",
					},
					Source:      "/home/code.rs",
					NoBuild:     false,
					BuildTarget: "/home/code.o",
					BuildImage:  "rust:alpine",
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
			{
				ID:      "java11",
				Name:    "Java/11",
				Enabled: true,
				Args: JudgerModels.CompileInfo{
					BuildArgs: []string{
						"javac",
						"/home/Main.java",
					},
					Source:      "/home/Main.java",
					NoBuild:     false,
					BuildTarget: "/home/Main.class",
					BuildImage:  "openjdk:11",
					Constraints: JudgerModels.Constraints{
						CPU:          1000000000,
						Memory:       1073741824,
						BuildTimeout: 120,
						RunTimeout:   120,
					},
					RunArgs: []string{
						"java",
						"-cp",
						"/home",
						"Main",
					},
					RunArgsJSON: "[\"java\",\"-cp\",\"/home\",\"Main\"]",
					RunImage:    "java_tester:latest",
				},
			},
			{
				ID:      "python3",
				Name:    "Python/3",
				Enabled: true,
				Args: JudgerModels.CompileInfo{
					BuildArgs:   []string{},
					Source:      "/home/code.py",
					NoBuild:     true,
					BuildTarget: "",
					BuildImage:  "-",
					Constraints: JudgerModels.Constraints{
						CPU:          1000000000,
						Memory:       1073741824,
						BuildTimeout: 120,
						RunTimeout:   120,
					},
					RunArgs: []string{
						"python",
						"/home/code.py",
					},
					RunArgsJSON: "[\"python\",\"/home/code.py\"]",
					RunImage:    "python_tester:latest",
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

	//os.Setenv("DEV", "1")
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
			return
		}
	}
}

func testJudgeHelper(code []byte, language string) (string, []*protobuf.JudgeCaseResult, error) {
	initJudger()

	config.Global.Extensions.HostBind = true
	status1, result1, err1 := Scheduler(&protobuf.JudgeRequest{
		Sid:        1,
		Tid:        1,
		Version:    "1",
		Language:   language,
		TimeLimit:  1000,
		SpaceLimit: 128,
		CompMode:   "STDIN_S",
		Code:       code,
		Time:       0,
		IsContest:  false,
	})

	fmt.Printf("[Result1] %+v \n", result1)

	config.Global.Extensions.HostBind = false
	status2, result2, err2 := Scheduler(&protobuf.JudgeRequest{
		Sid:        1,
		Tid:        1,
		Version:    "1",
		Language:   language,
		TimeLimit:  1000,
		SpaceLimit: 128,
		CompMode:   "STDIN_S",
		Code:       code,
		Time:       0,
		IsContest:  false,
	})

	fmt.Printf("[Result2] %+v \n", result2)

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

	status, _, _ := testJudgeHelper(code, "cpp17")
	if status != "CE" {
		t.Fail()
		return
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
	status, judgeResult, _ := testJudgeHelper(code, "cpp17")

	if status != "OK" {
		fmt.Println("[Should Emit RE] Status NOT OK")
		t.Fail()
		return
	}
	for _, result := range judgeResult {
		if result.Status != "RE" {
			fmt.Println("[Should Emit RE] Some Case Status NOT RE", result)
			t.Fail()
			return
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
	status, judgeResult, _ := testJudgeHelper(code, "cpp17")

	if status != "OK" {
		fmt.Println("[Should Emit TLE] Status NOT OK")
		t.Fail()
		return
	}
	for _, result := range judgeResult {
		if result.Status != "TLE" {
			fmt.Println("[Should Emit TLE] Some Case Status NOT TLE", result)
			t.Fail()
			return
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
	status, judgeResult, _ := testJudgeHelper(code, "cpp17")

	if status != "OK" {
		fmt.Println("[Should Emit AC] Status NOT OK")
		t.Fail()
		return
	}
	for _, result := range judgeResult {
		if result.Status != "AC" {
			fmt.Println("[Should Emit AC] Some Case Status NOT AC", result)
			t.Fail()
			return
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
	status, judgeResult, _ := testJudgeHelper(code, "cpp17")

	if status != "OK" {
		fmt.Println("[Should Emit WA] Status NOT OK")
		t.Fail()
		return
	}
	for _, result := range judgeResult {
		if result.Status != "WA" {
			fmt.Println("[Should Emit WA] Some Case Status NOT WA", result)
			t.Fail()
			return
		}
	}
}

func TestShouldEmitMLE(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%+v \n", err)
			t.Fail()
		}
	}()

	initJudger()

	code := []byte("#include <iostream> \n" +
		"#include <cstring> \n" +
		"using namespace std; \n" +
		"int main() { \n" +
		"    for (int i = 0; i < 10; i++) { \n" +
		"        int* a = new int[10000000]; \n" +
		"        memset(a, 0xff, 10000000 * sizeof(int)); \n" +
		"    } \n" +
		"    while (1) {} \n" +
		"    return 0; \n" +
		"}")
	status, judgeResult, _ := testJudgeHelper(code, "cpp17")

	if status != "OK" {
		fmt.Println("[Should Emit MLE] Status NOT OK")
		t.Fail()
		return
	}
	for _, result := range judgeResult {
		if result.Status != "MLE" {
			fmt.Println("[Should Emit MLE] Some Case Status NOT MLE", result)
			t.Fail()
			return
		}
	}
}

func TestJava11ShouldEmitAC(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%+v \n", err)
			t.Fail()
		}
	}()

	code := []byte("import java.io.*;\n " +
		"import java.util.*;\n " +
		"public class Main {\n " +
		"    public static class Rabbit {} \n" +
		"    public static void main(String args[]) throws Exception {\n " +
		"        Rabbit r = new Rabbit(); \n" +
		"        Rabbit rabbit = new Rabbit(); \n" +
		"        Scanner cin=new Scanner(System.in);\n " +
		"        int a = cin.nextInt(), b = cin.nextInt();\n " +
		"        System.out.println(a+b);\n " +
		"    } \n" +
		"}")
	status, judgeResult, _ := testJudgeHelper(code, "java11")

	if status != "OK" {
		fmt.Println("[Should Emit AC] Status NOT OK")
		t.Fail()
		return
	}
	for _, result := range judgeResult {
		if result.Status != "AC" {
			fmt.Println("[Should Emit AC] Some Case Status NOT AC", result)
			t.Fail()
			return
		}
	}
}

func TestPython3ShouldEmitAC(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%+v \n", err)
			t.Fail()
		}
	}()

	code := []byte("s = input().split()\n" +
		"print(int(s[0]) + int(s[1]))\n")
	status, judgeResult, _ := testJudgeHelper(code, "python3")

	if status != "OK" {
		fmt.Println("[Should Emit AC] Status NOT OK")
		t.Fail()
		return
	}
	for _, result := range judgeResult {
		if result.Status != "AC" {
			fmt.Println("[Should Emit AC] Some Case Status NOT AC", result)
			t.Fail()
			return
		}
	}
}

func TestRustShouldEmitAC(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%+v \n", err)
			t.Fail()
		}
	}()

	code := []byte("use std::io;\n\n" +
		"fn main(){\n" +
		"    let mut input=String::new();\n" +
		"    io::stdin().read_line(&mut input).unwrap();\n" +
		"    let mut s=input.trim().split(' ');\n\n" +
		"    let a:i32=s.next().unwrap()\n" +
		"               .parse().unwrap();\n" +
		"    let b:i32=s.next().unwrap()\n" +
		"               .parse().unwrap();\n" +
		"    println!(\"{}\",a+b);" +
		"\n}")
	status, judgeResult, _ := testJudgeHelper(code, "rust")

	if status != "OK" {
		fmt.Println("[Should Emit AC] Status NOT OK")
		t.Fail()
		return
	}
	for _, result := range judgeResult {
		if result.Status != "AC" {
			fmt.Println("[Should Emit AC] Some Case Status NOT AC", result)
			t.Fail()
			return
		}
	}
}
