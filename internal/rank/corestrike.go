package rank

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/haashi/omega-strikers-bot/internal/models"
)

func GetRankFromUsername(username string) (int, error) {
	url := fmt.Sprintf("https://corestrike.gg/lookup/%s", url.PathEscape(username))
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	reg := regexp.MustCompile(`Rating: (\d+) \(`)
	matches := reg.FindStringSubmatch(string(html))
	if len(matches) > 0 {
		rating, err := strconv.ParseInt(matches[1], 10, 0)
		return int(rating), err
	} else {
		return 0, &models.RankUpdateUsernameError{Username: username}
	}
}
