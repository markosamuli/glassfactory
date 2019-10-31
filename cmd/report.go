package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/markosamuli/glassfactory-cli/reporting"
	"github.com/markosamuli/glassfactory/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.AddCommand(monthlyReportsCmd)
	reportCmd.AddCommand(fiscalYearReportCmd)
}

func createReportingService(api *api.Service) (*reporting.Service, error) {
	ctx := context.Background()
	r, err := reporting.NewService(ctx, api)
	if err != nil {
		return nil, err
	}
	return r, nil
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Print time reports",
	Long:  `Print time reports for a user`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Time Reports")
	},
}

var monthlyReportsCmd = &cobra.Command{
	Use:   "monthly",
	Short: "Monthly time reports",
	Long:  `Print monthly time reports for the current calendar year`,
	Run: func(cmd *cobra.Command, args []string) {

		s, err := createApiService()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		member, err := s.GetCurrentMember()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		r, err := createReportingService(s)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		monthlyReports, err := r.MonthlyMemberTimeReports(member.ID, time.Now())
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		for _, r := range monthlyReports {
			r.RenderTable(os.Stdout)
		}

	},
}

var fiscalYearReportCmd = &cobra.Command{
	Use:   "fy",
	Short: "Fiscal year time report",
	Long:  `Print time reports for the current fiscal year`,
	Run: func(cmd *cobra.Command, args []string) {

		s, err := createApiService()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		member, err := s.GetCurrentMember()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		r, err := createReportingService(s)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fiscalYearFinalMonth := time.January
		fiscalYear := reporting.NewFiscalYear(time.Now().AddDate(-3, 0, 0), fiscalYearFinalMonth)
		annualReports, err := r.FiscalYearMemberTimeReports(member.ID, fiscalYear)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		for _, r := range annualReports {
			r.RenderTable(os.Stdout)
		}
	},
}