package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// RulesetConfig는 ruleset_config.yaml 전체 구조입니다.
type RulesetConfig struct {
	Version     string `yaml:"ruleset_version"`
	Description string `yaml:"description"`
	Rules       []Rule `yaml:"rules"`
}

// Rule은 개별 룰입니다.
type Rule struct {
	RuleID      string      `yaml:"rule_id"`
	Description string      `yaml:"description"`
	Conditions  []Condition `yaml:"conditions"`
}

// Condition은 룰의 개별 조건입니다.
// YAML의 value는 항상 문자열 배열이므로 []string으로 받습니다.
type Condition struct {
	Field    string   `yaml:"field"`
	Operator string   `yaml:"operator"`
	Value    []string `yaml:"value"`
}

// LoadRules는 파일 경로에서 룰셋을 로드합니다.
func LoadRules(path string) (*RulesetConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config RulesetConfig //위에서 정의한 struct
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
