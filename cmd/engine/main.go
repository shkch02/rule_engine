// go run ./cmd/engine/main.go test_logs.txt

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"rule_engine/pkg/alerter"
	"rule_engine/pkg/config"
	"rule_engine/pkg/engine"
	"rule_engine/pkg/input"
	"strings"
)

func init() {
	fmt.Println()
	fmt.Println(`                  %%@@@@@@@@%       %@@@@@@@@@&    `)
	fmt.Println(`                @@@@@@@@@@@@@@@  %@@@@@@@@@@@@@@@    _____       _        ______             _            `)
	fmt.Println(`    @@@@@@@@@ %@@@@@@%@%@@@@@@@@@@@@@@%@%%@@@@@@@@  |  __ \     | |      |  ____|           (_)           `)
	fmt.Println(`@@@@@@@@@@@@%%@@@@@%         %@@@@@@@       %@@@@@  | |__) |   _| | ___  | |__   _ __   __ _ _ _ __   ___ `)
	fmt.Println(`      @@@@@% %@@@@@   @@@@@@@@@@@@@@         @@@@@  |  _  / | | | |/ _ \ |  __| | '_ \ / _' | | '_ \ / _ \`)
	fmt.Println(`             %@@@@@   %%%%@@@@@@@@@@%       @@@@@@  | | \ \ |_| | |  __/ | |____| | | | (_| | | | | |  __/`)
	fmt.Println(`             %@@@@@@%%%%%@@@@@%@@@@@@%%  @@@@@@@@   |_|  \_\__,_|_|\___| |______|_| |_|\__, |_|_| |_|\___|`)
	fmt.Println(`              @@@@@@@@@@@@@@@@  @@@@@@@@@@@@@%@@                                        __/ |             `)
	fmt.Println(`                %@@@@@@@@%%       %@@@@@@@@@@%                                         |___/              `)
	fmt.Println()
	fmt.Println(" [ Rule Engine running... ]")
	fmt.Println()
}

const (
	defaultRulesetPath = "/etc/rules/rule.yaml" // 룰셋 파일 경로, 배포 야믈에서 연결 해줘야함
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. 룰셋에서 룰 야믈 입력받아 야믈에 맞는 구조체로 파싱
	ruleset, err := config.LoadRules(defaultRulesetPath)
	if err != nil {
		log.Fatalf("룰셋 파일 로드 실패 (%s): %v", defaultRulesetPath, err)
	}
	log.Printf("룰셋 로드 완료: %d개 룰", len(ruleset.Rules))

	// 2. 룰 엔진 및 Alerter 초기화 (PoC용 PrintAlerter 사용)
	ruleEngine := engine.NewRuleEngine(ruleset)
	pocAlerter := alerter.NewPrintAlerter()

	brokerStr := os.Getenv("KAFKA_BROKERS")
	if brokerStr == "" {
		log.Fatalf("KAFKA_BROKERS 환경 변수가 설정되지 않았습니다.")
	}
	brokers := strings.Split(brokerStr, ",")

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		log.Fatalf("KAFKA_TOPIC 환경 변수가 설정되지 않았습니다.")
	}
	groupID := "rule-engine-group" //카프카 컨슈머 그룹 아이디

	kafkaSource := input.NewKafkaSource(brokers, topic, groupID)
	eventCh, err := kafkaSource.Stream(ctx)
	if err != nil {
		log.Fatalf("카프카 소스 스트림 생성 실패: %v", err)
	}

	log.Println("이벤트 스트리밍 시작...")

	for event := range eventCh {
		violations := ruleEngine.Evaluate(event)

		if len(violations) > 0 {
			for _, v := range violations {
				pocAlerter.Alert(v)
			}
		}
	}
}
