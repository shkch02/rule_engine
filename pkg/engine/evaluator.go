package engine

import (
	"log"
	"rule_engine/pkg/config"
	"rule_engine/pkg/models"
	"strings"
)

// flagMap은 YAML의 문자열 플래그를 실제 정수 값(비트마스크)으로 매핑합니다.
// (참고: 실제 리눅스 값은 0, 1, 2가 아닐 수 있으나, 룰 파일의 주석을 따름)
var flagMap = map[string]int{
	"O_WRONLY": 1, // 룰 파일 주석 기준
	"O_RDWR":   2, // 룰 파일 주석 기준 (실제 리눅스 값은 2)
	//룰따라 추가 필요
}

type Evaluator struct{}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// Check는
func (e *Evaluator) Check(event models.Event, cond *config.Condition) bool {

	//이벤트(로그)내용중, 검사할 컨디션의 필드(ex syscall_name,flag)에 매칭되는 정보 저장
	eventValue, ok := event[cond.Field]

	if !ok {
		// 이벤트에 해당 필드가 없으면 조건 불일치
		return false
	}

	// 2. 필드 타입에 따라 평가 로직 분기
	switch cond.Field {
	case "syscall_name", "file_path":
		// 문자열 타입 필드 평가
		valStr, ok := eventValue.(string) //eventValue : 로그(이벤트)의 field 값
		if !ok {
			log.Printf("경고: 필드 타입 불일치 (필요: string, 실제: %T) for field %s", eventValue, cond.Field)
			return false
		}
		//   				이벤트.field 	룰.operator   룰.value
		return e.checkString(valStr, cond.Operator, cond.Value) //이벤트의 field가 룰의 condition에 만족하는지 검사 (만족 시 위반)

	case "flags":
		// 숫자(플래그) 타입 필드 평가
		// JSON은 정수를 float64로 파싱하므로 float64로 받음
		valFloat, ok := eventValue.(float64)
		if !ok {
			log.Printf("경고: 필드 타입 불일치 (필요: float64/number, 실제: %T) for field %s", eventValue, cond.Field)
			return false
		}
		// 실제 비트 연산은 정수(int)로 수행
		return e.checkFlags(int(valFloat), cond.Operator, cond.Value)
	}
	return false
}

// checkString은 문자열 필드를 평가합니다.
func (e *Evaluator) checkString(eventValue string, op string, ruleValue []string) bool {
	switch op {
	case "equals":
		// 룰에 정의된 값(들) 중 하나라도 일치하면
		for _, v := range ruleValue {
			if eventValue == v {
				return true
			}
		}
		return false
	case "starts_with_any":
		//이벤트 값 v가 ruleValue로 시작하는 문자열인지 확인
		for _, v := range ruleValue {
			if strings.HasPrefix(eventValue, v) {
				return true
			}
		}
		return false
	case "ends_with_any":
		//이벤트 값 v가 ruleValue로 끝나는 문자열인지 확인
		for _, v := range ruleValue {
			if strings.HasSuffix(eventValue, v) {
				return true
			}
		}
		return false
	}
	return false
}

// checkFlags는 비트 연산이 필요한 플래그를 평가합니다.
func (e *Evaluator) checkFlags(eventFlags int, op string, ruleFlags []string) bool {
	// 룰의 문자열 플래그(예: "O_WRONLY")를 비트마스크(정수)로 변환
	var ruleMask int
	for _, fStr := range ruleFlags {
		if val, ok := flagMap[fStr]; ok {
			ruleMask |= val // OR 연산으로 마스크 누적
		}
	}

	switch op {
	case "contains_any":
		// 이벤트 플래그와 룰 마스크를 AND 연산했을 때 0이 아니면
		// 룰이 요구하는 플래그 중 하나라도 포함된 것임 (예: 66 & 1 = 0, 66 & 2 = 2)
		return (eventFlags & ruleMask) != 0
	case "not_contains_any":
		// 이벤트 플래그와 룰 마스크를 AND 연산했을 때 0이어야 함
		// 룰이 요구하는 플래그가 하나도 없어야 함
		return (eventFlags & ruleMask) == 0
	}
	return false
}
