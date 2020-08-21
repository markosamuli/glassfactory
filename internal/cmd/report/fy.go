package report

import (
	"fmt"
	"os"
	"time"

	"github.com/markosamuli/glassfactory/internal/auth"
	"github.com/markosamuli/glassfactory/reporting"
	"github.com/spf13/cobra"
)

// FiscalYearReportOptions for the report command
type FiscalYearReportOptions struct {
}

// NewFiscalYearReportCommand creates new command
func NewFiscalYearReportCommand() *cobra.Command {
	var o = &FiscalYearReportOptions{}
	var c = &cobra.Command{
		Use:   "fy",
		Short: "Fiscal year time report",
		Long:  `Print time reports for the current fiscal year`,
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Run(cmd)
			if err != nil {
				fmt.Print(err)
			}
		},
	}
	return c
}

// Run the command
func (o *FiscalYearReportOptions) Run(cmd *cobra.Command) error {
	gfAuth, ok := auth.FromContext(cmd.Context())
	if !ok {
		return fmt.Errorf("failed to get authentication details")
	}

	s, err := gfAuth.NewService()
	if err != nil {
		return err
	}

	member, err := s.GetCurrentMember()
	if err != nil {
		return err
	}

	r, err := createReportingService(s)
	if err != nil {
		return err
	}

	fiscalYearFinalMonth := time.January
	fiscalYear := reporting.NewFiscalYear(time.Now().AddDate(-3, 0, 0), fiscalYearFinalMonth)
	annualReports, err := r.FiscalYearMemberTimeReports(member.ID, fiscalYear)
	if err != nil {
		return err
	}

	for _, r := range annualReports {
		r.RenderTable(os.Stdout)
	}
	return nil
}
