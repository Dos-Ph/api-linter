// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0158

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestResponseRepeatedFirstField(t *testing.T) {
	tests := []struct {
		testName    string
		FirstField  string
		SecondField string
		problems    testutils.Problems
	}{
		{
			"Valid",
			"repeated Book books = 1;",
			"string next_page_token = 2;",
			testutils.Problems{}},
		{
			// should have at least 1 field
			"InvalidType",
			"",
			"",
			testutils.Problems{{}}},
		{
			// first field by number should be repeated.
			"InvalidType",
			"repeated string next_page_token = 2;",
			"Book books = 1;",
			testutils.Problems{{}}},
		{
			// first field by position should be repeated.
			"InvalidType",
			"Book books = 2;",
			"repeated string next_page_token = 1;",
			testutils.Problems{{}}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create the proto message.
			f := testutils.ParseProto3Tmpl(t, `
				message Book {
					string name = 1;
				}
				message ListBooksResponse {
					{{.FirstField}}
					{{.SecondField}}
				}
			`, test)

			// Determine the descriptor that a failing test will attach to.
			if f.GetMessageTypes()[1].FindFieldByName("books") != nil {
				test.problems.SetDescriptor(f.GetMessageTypes()[1].FindFieldByName("books"))
			}

			// Run the lint rule and establish we get the correct problems.
			problems := responseRepeatedFirstField.Lint(f)
			if diff := test.problems.Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}

}