package rank

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/euscs/euscs-bot/internal/static"
)

func GetRankFromUsername(ctx context.Context, username string) (int, error) {
	corestrikeUrl := fmt.Sprintf("https://corestrike.gg/lookup/%s?region=Europe", url.PathEscape(username))
	resp, err := http.Get(corestrikeUrl)
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
		corestrikeUrl = fmt.Sprintf("https://corestrike.gg/lookup/%s", url.PathEscape(username))
		resp, err := http.Get(corestrikeUrl)
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
			return 0, static.ErrCorestrikeNotFound
		}
	}
}
