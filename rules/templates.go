// (c) Copyright 2016 Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rules

import (
	"go/ast"

	"github.com/securego/gas"
)

type templateCheck struct {
	gas.MetaData
	calls gas.CallList
}

func (t *templateCheck) ID() string {
	return t.MetaData.ID
}

func (t *templateCheck) Match(n ast.Node, c *gas.Context) (*gas.Issue, error) {
	if node := t.calls.ContainsCallExpr(n, c); node != nil {
		for _, arg := range node.Args {
			if _, ok := arg.(*ast.BasicLit); !ok { // basic lits are safe
				return gas.NewIssue(c, n, t.ID(), t.What, t.Severity, t.Confidence), nil
			}
		}
	}
	return nil, nil
}

// NewTemplateCheck constructs the template check rule. This rule is used to
// find use of tempaltes where HTML/JS escaping is not being used
func NewTemplateCheck(id string, conf gas.Config) (gas.Rule, []ast.Node) {

	calls := gas.NewCallList()
	calls.Add("html/template", "HTML")
	calls.Add("html/template", "HTMLAttr")
	calls.Add("html/template", "JS")
	calls.Add("html/template", "URL")
	return &templateCheck{
		calls: calls,
		MetaData: gas.MetaData{
			ID:         id,
			Severity:   gas.Medium,
			Confidence: gas.Low,
			What:       "this method will not auto-escape HTML. Verify data is well formed.",
		},
	}, []ast.Node{(*ast.CallExpr)(nil)}
}
