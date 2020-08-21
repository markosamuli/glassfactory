package report

import (
	"context"

	"github.com/markosamuli/glassfactory-cli/pkg/reporting"
	"github.com/markosamuli/glassfactory/api"
	"github.com/spf13/cobra"
)

func createReportingService(api *api.Service) (*reporting.Service, error) {
	ctx := context.Background()
	r, err := reporting.NewService(ctx, api)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// NewCommand creates new report command
func NewCommand() *cobra.Command {
	var c = &cobra.Command{
		Use:   "report",
		Short: "Print time reports",
		Long:  `Print time reports for a user`,
	}
	c.AddCommand(NewMonthlyReportCommand())
	c.AddCommand(NewFiscalYearReportCommand())
	return c

}
