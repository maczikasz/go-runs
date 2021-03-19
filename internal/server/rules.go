package server

import (
	"github.com/gin-gonic/gin"
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"net/http"
)

type RuleSaver interface {
	WriteRule(ruleType string, matcherType string, ruleContent string, runbookId string) error
	DeleteRule(ruleId string) error
}

type RuleFinder interface {
	FindOneRule(ruleType string, matcherType string, ruleContent string) (*model.RuleEntity, error)
	ListRules() (*[]model.RuleEntity, error)
}

type RuleReloader func()

type ruleHandler struct {
	ruleSaver    RuleSaver
	ruleFinder   RuleFinder
	ruleMatcher  runbooks.RuleMatcher
	ruleReloader RuleReloader
}

type RuleCreateDTO struct {
	RuleType    string `json:"rule_type"`
	MatcherType string `json:"matcher_type"`
	RuleContent string `json:"rule_content"`
	RunbookId   string `json:"runbook_id"`
}

func (h ruleHandler) AddNewRule(context *gin.Context) {
	var dto RuleCreateDTO
	err := context.BindJSON(&dto)

	if err != nil {
		//TODO fix all number status
		context.Status(http.StatusBadRequest)
		_ = context.Error(err)
		return
	}

	err = h.ruleSaver.WriteRule(dto.RuleType, dto.MatcherType, dto.RuleContent, dto.RunbookId)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		_ = context.Error(err)
		return
	}

	h.ruleReloader()

	context.Status(http.StatusCreated)
}

func (h ruleHandler) ListAllRules(context *gin.Context) {
	rules, err := h.ruleFinder.ListRules()
	if err != nil {
		context.Status(http.StatusInternalServerError)
		_ = context.Error(err)
		return
	}

	context.JSON(200, rules)
}

func (h ruleHandler) DisableRule(context *gin.Context) {
	ruleId := context.Param("ruleId")

	err := h.ruleSaver.DeleteRule(ruleId)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		_ = context.Error(err)
		return
	}

	context.Status(http.StatusNoContent)
}

func (h ruleHandler) TestRuleMatch(context *gin.Context) {
	var testError model.Error
	err := context.BindJSON(&testError)

	if err != nil {
		context.Status(http.StatusBadRequest)
		_ = context.Error(err)
		return
	}

	match, b := h.ruleMatcher.FindMatchingRunbook(testError)

	if !b {
		context.Status(http.StatusNotFound)
	} else {
		context.String(200, match)
	}
}
