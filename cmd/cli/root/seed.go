package root

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type createBeerStyleRequest struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	MinTemp float64 `json:"minTemp"`
	MaxTemp float64 `json:"maxTemp"`
}

func newSeedCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "seed",
		Short: "Seed beer styles into the running API",
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL := os.Getenv("API_BASE_URL")
			if baseURL == "" {
				baseURL = "http://localhost:8080"
			}

			client := &http.Client{
				Timeout: 5 * time.Second,
			}

			styles := []createBeerStyleRequest{
				{ID: "1", Name: "Weissbier", MinTemp: -1, MaxTemp: 3},
				{ID: "2", Name: "Pilsens", MinTemp: -2, MaxTemp: 4},
				{ID: "3", Name: "Weizenbier", MinTemp: -4, MaxTemp: 6},
				{ID: "4", Name: "Red Ale", MinTemp: -5, MaxTemp: 5},
				{ID: "5", Name: "IPA", MinTemp: -7, MaxTemp: 10},
				{ID: "6", Name: "Dunkel", MinTemp: -8, MaxTemp: 2},
				{ID: "7", Name: "Imperial Stouts", MinTemp: -10, MaxTemp: 13},
				{ID: "8", Name: "Brown Ale", MinTemp: 0, MaxTemp: 14},
			}

			for _, s := range styles {
				body, err := json.Marshal(s)
				if err != nil {
					return err
				}

				req, err := http.NewRequest(
					http.MethodPost,
					fmt.Sprintf("%s/beer-styles", baseURL),
					bytes.NewBuffer(body),
				)
				if err != nil {
					return err
				}

				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				// 201 = criado
				// 400 = já existe / inválido → ignoramos (idempotente)
				if resp.StatusCode != http.StatusCreated &&
					resp.StatusCode != http.StatusBadRequest {
					return fmt.Errorf(
						"failed to seed %s: status %d",
						s.Name,
						resp.StatusCode,
					)
				}
			}

			fmt.Println("Seed completed successfully")
			return nil
		},
	}
}
