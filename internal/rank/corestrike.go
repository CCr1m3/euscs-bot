package rank

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/euscs/euscs-bot/internal/static"
)

type CorestrikeResponse struct {
	RankedStats struct {
		Username      string  `json:"username"`
		Rating        int     `json:"rating"`
		RatingDisplay string  `json:"rating_display"`
		Rank          int     `json:"rank"`
		Role          string  `json:"role"`
		Wins          int     `json:"wins"`
		Losses        int     `json:"losses"`
		Winpercent    string  `json:"winpercent"`
		Toppercent    string  `json:"toppercent"`
		Verified      bool    `json:"verified"`
		IsRanked      bool    `json:"is_ranked"`
		LpHistory     [][]any `json:"lp_history"`
	} `json:"rankedStats"`
	CharacterStats struct {
		Forwards []struct {
			Name        string `json:"name"`
			DisplayName string `json:"display_name"`
			Wins        int    `json:"wins"`
			Losses      int    `json:"losses"`
			Assists     int    `json:"assists"`
			Mvp         int    `json:"mvp"`
			Knockouts   int    `json:"knockouts"`
			Scores      int    `json:"scores"`
			Saves       int    `json:"saves"`
		} `json:"forwards"`
		Goalies []struct {
			Name        string `json:"name"`
			DisplayName string `json:"display_name"`
			Wins        int    `json:"wins"`
			Losses      int    `json:"losses"`
			Assists     int    `json:"assists"`
			Mvp         int    `json:"mvp"`
			Knockouts   int    `json:"knockouts"`
			Scores      int    `json:"scores"`
			Saves       int    `json:"saves"`
		} `json:"goalies"`
	} `json:"characterStats"`
	OverallStats struct {
		Forwards struct {
			Name        string `json:"name"`
			DisplayName string `json:"display_name"`
			Wins        int    `json:"wins"`
			Losses      int    `json:"losses"`
			Assists     int    `json:"assists"`
			Mvp         int    `json:"mvp"`
			Knockouts   int    `json:"knockouts"`
			Scores      int    `json:"scores"`
			Saves       int    `json:"saves"`
		} `json:"forwards"`
		Goalies struct {
			Name        string `json:"name"`
			DisplayName string `json:"display_name"`
			Wins        int    `json:"wins"`
			Losses      int    `json:"losses"`
			Assists     int    `json:"assists"`
			Mvp         int    `json:"mvp"`
			Knockouts   int    `json:"knockouts"`
			Scores      int    `json:"scores"`
			Saves       int    `json:"saves"`
		} `json:"goalies"`
	} `json:"overallStats"`
	Error string `json:"error"`
}

func GetCorestrikeInfoFromUsername(ctx context.Context, username string) (*CorestrikeResponse, error) {
	corestrikeUrl := fmt.Sprintf("https://corestrike.gg/lookup/%s?region=Global&json=true", url.PathEscape(username))
	resp, err := http.Get(corestrikeUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response CorestrikeResponse
	err = json.Unmarshal(jsonBytes, &response)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(response.Error, "") {
		if !response.RankedStats.IsRanked {
			corestrikeUrl := fmt.Sprintf("https://corestrike.gg/lookup/%s?region=Europe&json=true", url.PathEscape(username))
			resp, err := http.Get(corestrikeUrl)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			jsonBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(jsonBytes, &response)
			if err != nil {
				return nil, err
			}
		}
		if response.RankedStats.IsRanked {
			return &response, nil
		} else {
			return &response, static.ErrUnrankedUser
		}
	} else if strings.EqualFold(response.Error, "Invalid username") {
		return nil, static.ErrUsernameInvalid
	} else {
		return nil, static.ErrCorestrikeNotFound
	}
}
