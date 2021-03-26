package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"github.com/maczikasz/go-runs/internal/server/dto"
	"net/http"
)

//go:generate moq -out ../mocks/rules.go --skip-ensure . RuleSaver RuleFinder

type (
	RuleSaver interface {
		WriteRule(ruleType string, matcherType string, ruleContent string, runbookId string) error
		DeleteRule(ruleId string) error
	}

	RuleFinder interface {
		FindOneRule(ruleType string, matcherType string, ruleContent string) (*model.RuleEntity, error)
		ListRules() (*[]model.RuleEntity, error)
	}

	RuleReloader func()

	RuleHandler struct {
		ruleSaver    RuleSaver
		ruleFinder   RuleFinder
		ruleMatcher  runbooks.RuleMatcher
		ruleReloader RuleReloader
	}
)

func NewRuleHandler(ruleSaver RuleSaver, ruleFinder RuleFinder, ruleMatcher runbooks.RuleMatcher, ruleReloader RuleReloader) *RuleHandler {
	return &RuleHandler{ruleSaver: ruleSaver, ruleFinder: ruleFinder, ruleMatcher: ruleMatcher, ruleReloader: ruleReloader}
}

func (h RuleHandler) AddNewRule(context *gin.Context) {
	var ruleDto dto.RuleCreateDTO
	err := context.BindJSON(&ruleDto)

	if err != nil {
		//TODO fix all number status
		context.Status(http.StatusBadRequest)
		_ = context.Error(err)
		return
	}

	err = h.ruleSaver.WriteRule(ruleDto.RuleType, ruleDto.MatcherType, ruleDto.RuleContent, ruleDto.RunbookId)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		_ = context.Error(err)
		return
	}

	h.ruleReloader()

	context.Status(http.StatusCreated)
}

func (h RuleHandler) ListAllRules(context *gin.Context) {
	rules, err := h.ruleFinder.ListRules()
	if err != nil {
		context.Status(http.StatusInternalServerError)
		_ = context.Error(err)
		return
	}

	context.JSON(200, rules)
}

func (h RuleHandler) DisableRule(context *gin.Context) {
	ruleId := context.Param("ruleId")

	err := h.ruleSaver.DeleteRule(ruleId)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		_ = context.Error(err)
		return
	}

	context.Status(http.StatusNoContent)
}

func (h RuleHandler) TestRuleMatch(context *gin.Context) {
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

//TODO rewrite to use the same objectid
func (h RuleHandler) UpdateRule(context *gin.Context) {
	ruleId := context.Param("ruleId")

	var createDto dto.RuleCreateDTO
	err := context.BindJSON(&createDto)

	if err != nil {
		//TODO fix all number status
		context.Status(http.StatusBadRequest)
		_ = context.Error(err)
		return
	}

	err = h.ruleSaver.WriteRule(createDto.RuleType, createDto.MatcherType, createDto.RuleContent, createDto.RunbookId)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		_ = context.Error(err)
		return
	}

	err = h.ruleSaver.DeleteRule(ruleId)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		_ = context.Error(err)
		return
	}

	h.ruleReloader()

	context.Status(http.StatusCreated)

}
