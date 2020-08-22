package model

import (
	"encoding/json"
	"strings"
	"testing"

	"gotest.tools/assert"
)

func removeSpace(in string) string {
	out := strings.ReplaceAll(in, " ", "")
	out = strings.ReplaceAll(out, "\t", "")
	out = strings.ReplaceAll(out, "\n", "")
	return out
}

func TestBillableStatus_MarshalJSON(t *testing.T) {
	var target BillableStatus
	target = 999
	_, err := json.Marshal(&target)
	assert.ErrorContains(t, err, "invalid BillableStatus")
}

func TestBillableStatus_UnmarshalJSON(t *testing.T) {
	var target BillableStatus
	data := `"invalid_status"`
	err := json.Unmarshal([]byte(data), &target)
	assert.ErrorContains(t, err, "invalid BillableStatus")
}

func TestBillableStatus(t *testing.T) {
	type mockProject struct {
		BillableStatus BillableStatus `json:"billable_status"`
	}
	var tests = []struct {
		data       string
		wantValue  BillableStatus
		wantString string
	}{
		{
			data: `{
			  "billable_status": null
			}`,
			wantValue:  Unknown,
			wantString: "Unknown",
		},
		{
			data: `{
			  "billable_status": "billable"
			}`,
			wantValue:  Billable,
			wantString: "Billable",
		},
		{
			data: `{
			  "billable_status": "non_billable"
			}`,
			wantValue:  NonBillable,
			wantString: "Non Billable",
		},
		{
			data: `{
			  "billable_status": "new_business"
			}`,
			wantValue:  NewBusiness,
			wantString: "New Business",
		},
	}
	for _, tt := range tests {
		t.Run(tt.wantString, func(t *testing.T) {
			var target mockProject
			err := json.Unmarshal([]byte(tt.data), &target)
			assert.NilError(t, err)
			assert.Equal(t, target.BillableStatus, tt.wantValue)
			assert.Equal(t, target.BillableStatus.String(), tt.wantString)
			encoded, err := json.Marshal(&target)
			assert.Equal(t, string(encoded), removeSpace(tt.data))
		})
	}
}
