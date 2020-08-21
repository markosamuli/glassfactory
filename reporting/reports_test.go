package reporting

import (
	"testing"

	"github.com/markosamuli/glassfactory/model"
	"gotest.tools/assert"
)

func TestFormatBillableStatus(t *testing.T) {
	billable := &model.Project{BillableStatus: model.Billable}
	billableReport := &MonthlyTimeReport{
		Project: billable,
	}
	assert.Equal(t, billableReport.BillableStatus(), "Billable")

	nonBillable := &model.Project{BillableStatus: model.NonBillable}
	nonBillableReport := &AnnualTimeReport{
		Project: nonBillable,
	}
	assert.Equal(t, nonBillableReport.BillableStatus(), "Non Billable")

	newBusiness := &model.Project{BillableStatus: model.NewBusiness}
	newBusinessReport := &FiscalYearTimeReport{
		Project: newBusiness,
	}
	assert.Equal(t, newBusinessReport.BillableStatus(), "New Business")

	unknown := &model.Project{BillableStatus: model.Unknown}
	unknownReport := &FiscalYearTimeReport{
		Project: unknown,
	}
	assert.Equal(t, unknownReport.BillableStatus(), "Unknown")
}
