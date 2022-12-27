package main

import (
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

func ensureServiceID(promql string, pod string) (string, error) {
	expr, err := parser.ParseExpr(promql)
	if err != nil {
		return "", err
	}

	ensureLabel(expr, pod)
	return expr.String(), nil
}

func ensureLabel(expr parser.Expr, pod string) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *parser.AggregateExpr:
		ensureLabel(e.Expr, pod)

	case *parser.Call:
		for _, arg := range e.Args {
			ensureLabel(arg, pod)
		}

	case *parser.ParenExpr:
		ensureLabel(e.Expr, pod)

	case *parser.UnaryExpr:
		ensureLabel(e.Expr, pod)

	case *parser.BinaryExpr:
		ensureLabel(e.LHS, pod)
		ensureLabel(e.RHS, pod)

	case *parser.NumberLiteral:
		return

	case *parser.VectorSelector:
		flag := false
		for _, lm := range e.LabelMatchers {
			if lm.Name == "pod" {
				lm.Value = pod
				flag = true
			}
		}
		if !flag {
			lm, _ := labels.NewMatcher(labels.MatchEqual, "pod", pod)
			e.LabelMatchers = append(e.LabelMatchers, lm)
		}
	case *parser.MatrixSelector:
		ensureLabel(e.VectorSelector, pod)

	case *parser.SubqueryExpr:
		ensureLabel(e.Expr, pod)

	case *parser.StringLiteral:
		return
	}
}
