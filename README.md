# eBPF 시스템콜 로그 기반 보안 룰 엔진

## 1.프로젝트 개요
본 프로젝트는 GoLang으로 작성된 실시간 eBPF 시스템콜 로그 분석 및 보안 룰 엔진입니다. 본 프로젝트는 중앙집중형 stateful set pod으로 배포되어 eBPF 에이전트(Collector)에 의해 수집된 시스템콜 이벤트를(현재 PoC 단계에서는 로그 파일, 최종적으로는 Kafka 스트림) 입력받아, YAML에 정의된 보안 룰셋과 비교 분석합니다.

분석 결과, 룰을 위반하는 행위(예: 민감 파일 접근, 권한 상승 시도)가 감지되면 즉시 알림(현재 콘솔 출력, 추후 웹훅 부착)을 발생시키는 것을 목적으로 합니다.

## 2. 주요 기능
* **파일 기반 분석(PoC에서만 적용)** : main.go가 로그 파일을 직접 읽어 한 줄씩 파싱하고 룰 위반 여부를 즉시 분석합니다.

* **동적 로그 파싱** : eBPF가 생성하는 다양한 구조의 JSON 로그(시스템콜마다 필드가 다름)를 처리하기 위해 map[string]interface{}를 사용하여 유연하게 파싱합니다.

* **YAML 기반 룰셋** : ruleset_config.yaml 파일을 통해 보안 탐지 룰을 유연하게 정의하고 수정할 수 있습니다.

* **룰 평가 엔진** : equals, contains_any, not_contains_any, starts_with_any 등 다양한 operator를 지원하며, flags 필드에 대한 비트 연산(Bitwise operation) 평가를 지원합니다.

* **확장 가능한 아키텍처** : Alerter (알림) 및 Source (입력) 인터페이스를 정의하여, 향후 Kafka Consumer, Slack/Webhook 알림 등 새로운 모듈을 쉽게 추가할 수 있도록 설계되었습니다.

## 3. 요구사항
* **GoLang** : Go 1.21 이상

* **Go 의존성** : gopkg.in/yaml.v3 (룰셋 파싱용)

                gopkg.in/check.v1 (yaml.v3의 하위/간접 의존성)

                go mod download {의존성 패키지}

* **(향후) 스트리밍** : Apache Kafka (실시간 분석 모드 시)

* **(향후) 배포** : Kubernetes (StatefulSet 배포 시)


## 4. 사용 방법

현재 PoC(Proof of Concept)는 로컬 로그 파일을 인자로 받아 실행됩니다.

#### 1. (선택) 바이너리 빌드
```
# PoC 모듈 빌드
go build -o rule-engine ./cmd/engine/
```

#### 2. 룰 엔진 실행
go run 으로 빌드와 실행

```
# test_logs.txt 파일을 인자로 전달
go run ./cmd/engine/main.go test_logs.txt
```

1.의 바이너리 빌드를 수행한 경우
```
./rule-engine test_logs.txt
```

#### 3. 룰셋 수정

ruleset_config.yaml 파일의 룰을 수정한 뒤, 분석기를 다시 실행하여 변경된 룰이 적용되는지 확인할 수 있습니다.

## 5. 프로젝트 구조
```
/rule-engine
|
|-- go.mod                  # Go 모듈 파일
|
|-- cmd/engine/
|   |-- main.go             # (✔️ PoC) 현재의 실행 파일. 로그 파일 파싱 및 엔진 실행
|
|-- pkg/
|   |-- config/
|   |   |-- rules.go        # (✔️ PoC) ruleset_config.yaml 로딩
|   |
|   |-- models/
|   |   |-- models.go       # (✔️ PoC) 'Event' (map[string]interface{}) 및 'Violation' 정의
|   |
|   |-- engine/
|   |   |-- engine.go       # (✔️ PoC) 룰 엔진 코어 (이벤트 정규화, 룰 평가 요청)
|   |   |-- evaluator.go    # (✔️ PoC) 룰 조건(operator) 실제 평가 로직
|   |
|   |-- input/
|   |   |-- input.go        # (Staged) 향후 Kafka/File 소스용 인터페이스 (현재는 비어있음)
|   |   |-- kafka.go        # (Stub) 향후 Kafka Consumer 구현을 위한 스텁
|   |
|   |-- alerter/
|   |   |-- alerter.go      # (✔️ PoC) 알림용 인터페이스
|   |   |-- print.go        # (✔️ PoC) PoC용 콘솔(print) 알림 구현
|   |   |-- slack.go        # (Stub) 향후 Slack 알림 구현을 위한 스텁
|
|-- ruleset_config.yaml     # (✔️ PoC) 제공해주신 룰셋
|
|-- test_logs.txt           # (✔️ PoC) PoC 테스트용 샘플 로그
```

## 6. 향후 개선 사항



## 7. 라이센스

## ToDo

#### engine.normalize 개선 필요, 인자 값 처리 고민 필요

#### 실행시간 단축 고민

#### 어떻게 카프카로 스트림을 구성할 것인가

### main.json.unmarshal 여기 지금 0x주소값 구조체 파싱 실패중, 넘길때 0x땔지 여기서 처리할지 고민중