// go run ./cmd/engine/main.go test_logs.txt

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"rule_engine/pkg/alerter"
	"rule_engine/pkg/config"
	"rule_engine/pkg/engine"
	"rule_engine/pkg/models"
)

func init() {
	fmt.Println()
	fmt.Println(`  _____       _        ______             _            `)
	fmt.Println(` |  __ \     | |      |  ____|           (_)           `)
	fmt.Println(` | |__) |   _| | ___  | |__   _ __   __ _ _ _ __   ___ `)
	fmt.Println(` |  _  / | | | |/ _ \ |  __| | '_ \ / _' | | '_ \ / _ \`)
	fmt.Println(` | | \ \ |_| | |  __/ | |____| | | | (_| | | | | |  __/`)
	fmt.Println(` |_|  \_\__,_|_|\___| |______|_| |_|\__, |_|_| |_|\___|`)
	fmt.Println(`                                     __/ |             `)
	fmt.Println(`                                    |___/              `)
	fmt.Println()
	fmt.Println(" [ Rule Engine running... ]")
	fmt.Println()
}

const (
	// PoC를 위해 룰셋 파일 경로를 하드코딩합니다.
	// K8s 배포 시 ConfigMap을 통해 이 경로에 마운트할 수 있습니다.
	defaultRulesetPath = "ruleset_config.yaml"
)

func main() {
	log.Println("룰 엔진 PoC 시작...")

	// 1. 커맨드 라인 인자로부터 로그 파일 경로 획득
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "사용법: %s <log_file_path>\n", os.Args[0])
		os.Exit(1)
	}
	logFilePath := os.Args[1]

	// 2. 룰셋 로드
	ruleset, err := config.LoadRules(defaultRulesetPath)
	if err != nil {
		log.Fatalf("룰셋 파일 로드 실패 (%s): %v", defaultRulesetPath, err)
	}
	log.Printf("룰셋 로드 완료: %d개 룰", len(ruleset.Rules))

	// 3. 룰 엔진 및 Alerter 초기화 (PoC용 PrintAlerter 사용)
	ruleEngine := engine.NewRuleEngine(ruleset)
	pocAlerter := alerter.NewPrintAlerter()

	// 4. 로그 파일 열기
	file, err := os.Open(logFilePath)
	if err != nil {
		log.Fatalf("로그 파일 열기 실패 (%s): %v", logFilePath, err)
	}
	defer file.Close()

	// 5. 파일 스캔 및 이벤트 처리
	scanner := bufio.NewScanner(file)
	lineNum := 0
	log.Printf("로그 파일 스캔 시작: %s", logFilePath)

	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()

		if len(line) == 0 {
			continue
		}

		// 6. 동적 JSON 파싱 (핵심)
		// SyscallEvent 구조체 대신 map[string]interface{} 사용
		//ToDO : 여기 지금 0x주소값 구조체 파싱 실패중, 넘길떄 0x땔지 여기서 처리할지 고민중
		var event models.Event
		if err := json.Unmarshal(line, &event); err != nil {
			log.Printf("경고: 로그 파싱 실패 (라인 %d): %v", lineNum, err)
			continue
		}
		log.Printf("디버그: 이벤트 객체 (라인 %d): %+v", lineNum, event)

		// 7. 룰 엔진에 이벤트 전달 및 평가
		violations := ruleEngine.Evaluate(event)

		// 8. 위반 사항 알림 (PrintAlerter가 콘솔에 출력)
		if len(violations) > 0 {
			for _, v := range violations {
				pocAlerter.Alert(v)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("로그 파일 스캔 중 에러: %v", err)
	}

	log.Println("룰 엔진 PoC 종료.")
}
