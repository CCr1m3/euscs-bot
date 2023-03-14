package rank

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/euscs/euscs-bot/internal/static"
)

type CorestrikeRankedStats struct {
	Rating    int      `json:"rating"`
	Rank      int      `json:"rank"`
	Role      string   `json:"role"`
	Wins      int      `json:"wins"`
	Losses    int      `json:"losses"`
	LPHistory [][2]int `json:"lp_history"`
}
type CorestrikeResponse struct {
	RankedStats    CorestrikeRankedStats `json:"rankedStats"`
	PlayerID       string                `json:"playerID"`
	EquippedTitle  string                `json:"equippedTitle"`
	EquippedBanner string                `json:"equippedBanner"`
	EquippedAvatar string                `json:"equippedAvatar"`
	Error          string                `json:"error"`
}

func GetCorestrikeInfoFromUsername(ctx context.Context, username string) (*CorestrikeResponse, error) {
	corestrikeUrl := fmt.Sprintf("https://corestrike.gg/lookup/%s?region=Europe&json=true", url.PathEscape(username))
	resp, err := http.Get(corestrikeUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jsonStr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response CorestrikeResponse
	err = json.Unmarshal(jsonStr, &response)
	if err != nil {
		return nil, err
	}
	if response.Error == "" {
		return &response, nil
	} else {
		corestrikeUrl = fmt.Sprintf("https://corestrike.gg/lookup/%s&json=true", url.PathEscape(username))
		resp, err = http.Get(corestrikeUrl)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		jsonStr, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var response CorestrikeResponse
		err = json.Unmarshal(jsonStr, &response)
		if err != nil {
			return nil, err
		}
		if response.Error == "" {
			return &response, nil
		} else {
			return nil, static.ErrCorestrikeNotFound
		}
	}
}
