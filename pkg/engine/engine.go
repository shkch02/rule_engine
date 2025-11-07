package engine

import (
	"rule_engine/pkg/config"
	"rule_engine/pkg/models"
)

// RuleEngine은 로드된 룰셋과 평가기를 보유합니다.
type RuleEngine struct {
	ruleset *config.RulesetConfig
	eval    *Evaluator
}

// NewRuleEngine은 새 룰 엔진 인스턴스를 생성합니다.
func NewRuleEngine(ruleset *config.RulesetConfig) *RuleEngine {
	return &RuleEngine{
		ruleset: ruleset,
		eval:    NewEvaluator(), // 룰 평가기 생성
	}
}

// Evaluate는 map[key:value ...]구조체인 event를 인자로 받아 단일 이벤트를 모든 룰과 비교하여 위반 목록을 반환합니다.
func (e *RuleEngine) Evaluate(event models.Event) []models.Violation {
	var violations []models.Violation

	// 1. 이벤트 정규화 (필드 이름 통일)
	e.normalize(event)

	// 2. 모든 룰에 대해 검사
	for _, rule := range e.ruleset.Rules {
		if e.checkRule(event, &rule) {
			violations = append(violations, models.Violation{
				RuleID:      rule.RuleID,
				Description: rule.Description,
				Event:       event,
			})
		}
	}
	return violations
}

// checkRule은 이벤트가 단일 룰의 *모든* 조건(AND)을 만족하는지 검사합니다.
func (e *RuleEngine) checkRule(event models.Event, rule *config.Rule) bool {
	for _, cond := range rule.Conditions {
		// Evaluator가 실제 조건 평가를 수행
		if !e.eval.Check(event, &cond) {
			// 하나라도 조건이 틀리면(false) 이 룰은 위반이 아님
			return false
		}
	}
	// 모든 조건을 통과하면 룰 위반
	return true
}

// normalize는 eBPF 로그의 필드 이름(예: "type", "pathname")을
// 룰셋이 정의한 필드 이름(예: "syscall_name", "file_path")으로 통일합니다.
func (e *RuleEngine) normalize(event models.Event) {
	// 로그: {"type": "openat"} -> 룰: {"syscall_name": "openat"}
	if val, ok := event["type"]; ok {
		event["syscall_name"] = val
	}
	// 로그: {"pathname": "/etc/passwd"} -> 룰: {"file_path": "/etc/passwd"}
	if val, ok := event["pathname"]; ok {
		event["file_path"] = val
	}
	// 'flags'는 로그와 룰 필드 이름이 동일하므로 변환 필요 없음
}
