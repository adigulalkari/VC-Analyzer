package subcommands

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/adigulalkari/VC-Analyzer/pkg/format"
	"github.com/adigulalkari/VC-Analyzer/pkg/utils"
)

var (
	Repository string // Exported to be accessible outside the package
)

var GetCmd = &cobra.Command{ // Exported to be accessible outside the package
	Use:   "get <detail>",
	Short: "Display one or many repositories details",
	Example: heredoc.Doc(`
        Get stars count for a given repository
        $ vc-analyze get stars -r golang/go
    `),
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a detail argument")
		}

		possibleDetails := []string{"stars"}
		validResource := false

		for _, resource := range possibleDetails {
			if args[0] == resource {
				validResource = true
				break
			}
		}

		if !validResource {
			return fmt.Errorf(`detail "%s" is invalid`, args[0])
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if !format.IsRepositoryValid(Repository) {
			return errors.New("repository flag is required and must be in the format 'owner/repo'")
		}

		detail := args[0]

		switch detail {
		case "stars":
			{
				starsCount, err := utils.GetStarCount(Repository)
				if err != nil {
					return fmt.Errorf("could not get stars count for this repository: %w", err)
				}

				fmt.Printf("%s has %d stars\n", Repository, *starsCount)
			}
		}

		return nil
	},
}
