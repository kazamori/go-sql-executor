package query

import (
	"errors"
	"regexp"
	"strings"

	"github.com/cybozu-go/log"
)

const (
	commentSeparator = "--"
	comma            = ","
)

var (
	errNoBindVariable     = errors.New("not found bind variable")
	errNoParameter        = errors.New("not found parameter")
	errUnmatchVarAndParam = errors.New(
		"the number of bind variables and parameters is not same",
	)

	reBindVariable = regexp.MustCompile(
		`[[:space:]]*([$:][[:alnum:]]+)[[:space:]]*[,;\)]*`,
	)
	reParameter = regexp.MustCompile(
		`[[:space:]]*\[(.*?)\][[:space:]]*`,
	)
)

func GetBindVariable(sql string) []string {
	matches := reBindVariable.FindAllStringSubmatch(sql, -1)
	vars := make([]string, 0, len(matches))
	for _, m := range matches {
		// FIXME: incomplete, separate string literal
		vars = append(vars, m[1])
	}
	log.Debug("bind variable", map[string]interface{}{
		"sql":  sql,
		"vars": vars,
	})
	return vars
}

func HasBindVariable(sql string) bool {
	// FIXME: incomplete, separate string literal
	return reBindVariable.MatchString(sql)
}

type SimpleParser struct {
	sql    string
	vars   []string
	params []string
}

func (p *SimpleParser) getBindVariable() error {
	p.vars = GetBindVariable(p.sql)
	if len(p.vars) == 0 {
		return errNoBindVariable
	}
	return nil
}

func (p *SimpleParser) getParameter() error {
	s := strings.Split(p.sql, commentSeparator)
	if len(s) < 2 {
		return errNoParameter
	}

	paramStr := s[1]
	matches := reParameter.FindAllStringSubmatch(paramStr, -1)
	if len(matches) == 0 {
		return errNoParameter
	}

	_values := strings.Split(matches[0][1], comma)
	p.params = make([]string, 0, len(_values))
	for _, value := range _values {
		v := strings.ReplaceAll(value, `"`, `'`)
		v = strings.TrimSpace(v)
		p.params = append(p.params, v)
	}

	log.Debug("parameter", map[string]interface{}{
		"s":       paramStr,
		"matches": matches,
		"params":  p.params,
	})
	return nil
}

func (p *SimpleParser) Parse() error {
	if err := p.getBindVariable(); err != nil {
		return err
	}
	if err := p.getParameter(); err != nil {
		return err
	}
	return nil
}

func (p *SimpleParser) Replace() (string, error) {
	if len(p.vars) != len(p.params) {
		log.Error(
			"number of bind variables and parameters",
			map[string]interface{}{
				"vars_length":   len(p.vars),
				"params_length": len(p.params),
			})
		return "", errUnmatchVarAndParam
	}

	_s := strings.Split(p.sql, commentSeparator)
	s := _s[0]
	for _, vp := range Zip(p.vars, p.params) {
		bindVar := vp[0]
		param := vp[1]
		s = strings.Replace(s, bindVar, param, 1)
	}

	s = strings.TrimSpace(s)
	log.Debug("replace", map[string]interface{}{
		"s": s,
	})
	return s, nil
}

func NewSimpleParser(sql string) *SimpleParser {
	return &SimpleParser{
		sql: sql,
	}
}
