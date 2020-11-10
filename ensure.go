package main

import (
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

func ensureServiceID(promql string, serviceID string) (string, error) {
	expr, err := parser.ParseExpr(promql)
	if err != nil {
		return "", err
	}

	ensureLabel(expr, serviceID)
	return expr.String(), nil
}

func ensureLabel(expr parser.Expr, serviceID string) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *parser.AggregateExpr:
		ensureLabel(e.Expr, serviceID)

	case *parser.Call:
		for _, arg := range e.Args {
			ensureLabel(arg, serviceID)
		}

	case *parser.ParenExpr:
		ensureLabel(e.Expr, serviceID)

	case *parser.UnaryExpr:
		ensureLabel(e.Expr, serviceID)

	case *parser.BinaryExpr:
		ensureLabel(e.LHS, serviceID)
		ensureLabel(e.RHS, serviceID)

	case *parser.NumberLiteral:
		return

	case *parser.VectorSelector:
		flag := false
		for _, lm := range e.LabelMatchers {
			if lm.Name == "service_id" {
				lm.Value = serviceID
				flag = true
			}
		}
		if !flag {
			lm, _ := labels.NewMatcher(labels.MatchEqual, "service_id", serviceID)
			e.LabelMatchers = append(e.LabelMatchers, lm)
		}
	case *parser.MatrixSelector:
		ensureLabel(e.VectorSelector, serviceID)

	case *parser.SubqueryExpr:
		ensureLabel(e.Expr, serviceID)

	case *parser.StringLiteral:
		return
	}
}
