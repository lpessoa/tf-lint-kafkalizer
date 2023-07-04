package rules

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

type TopicNameRule struct {
	tflint.DefaultRule
}

func NewSchemaRule() *TopicNameRule {
	return &TopicNameRule{}
}

func (s *TopicNameRule) Name() string {
	return "Kafkalizer topic name rule"
}

func (s *TopicNameRule) Enabled() bool {
	return true
}

func (s *TopicNameRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (s *TopicNameRule) Check(runner tflint.Runner) error {
	runner.WalkExpressions(tflint.ExprWalkFunc(func(expr hcl.Expression) hcl.Diagnostics {
		switch expr.(type) {
		case *hclsyntax.TupleConsExpr:
			{
				x := cty.Value{}
				err := runner.EvaluateExpr(expr, &x, &tflint.EvaluateExprOption{})

				if err != nil {
					return nil
				}

				var invalid string
				if x.CanIterateElements() {
					stop := x.ForEachElement(func(key cty.Value, val cty.Value) (stop bool) {
						m := val.AsValueMap()
						if d, ok := m["topic_name"]; ok {
							if d.AsString() == "sandbox.user" {
								stop = true
								invalid = d.AsString()
							}
						}
						return stop
					})
					if stop {
						err := runner.EmitIssue(s, fmt.Sprintf("Naming convention for topic_name is invalid: '%s'", invalid), expr.Range())
						if err != nil {
							return nil
						}
						return nil
					}
				}
			}
		}
		return nil
	}))
	/*// This rule is an example to get attributes of blocks other than resources.
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "topics_list_entry",
				Body: &hclext.BodySchema{},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeExpand})
	if err != nil {
		return err
	}

	for _, variable := range content.Blocks {
		fmt.Printf("%s", variable.Body.Attributes["topic_name"])
	}*/
	return nil
}
