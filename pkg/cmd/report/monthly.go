package report

import (
	"fmt"
	"os"
	"time"

	"github.com/markosamuli/glassfactory-cli/pkg/auth"
	"github.com/spf13/cobra"
)

// MonthlyReportOptions for the report command
type MonthlyReportOptions struct {
}

// NewMonthlyReportCommand creates new command
func NewMonthlyReportCommand() *cobra.Command {
	var o = &MonthlyReportOptions{}
	var c = &cobra.Command{
		Use:   "monthly",
		Short: "Monthly time reports",
		Long:  `Print monthly time reports for the current calendar year`,
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
func (o *MonthlyReportOptions) Run(cmd *cobra.Command) error {
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

	monthlyReports, err := r.MonthlyMemberTimeReports(member.ID, time.Now())
	if err != nil {
		return err
	}
	for _, r := range monthlyReports {
		r.RenderTable(os.Stdout)
	}
	return nil
}
