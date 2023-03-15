package rank

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type LeaderboardPlayer struct {
	Username    string `json:"username" db:"username"`
	PlayerID    string `json:"playerId" db:"playerID"`
	LogoID      string `json:"logoId" db:"logoID"`
	TitleID     string `json:"titleId" db:"titleID"`
	NameplateID string `json:"nameplateId" db:"nameplateID"`
	EmoticonID  string `json:"emoticonId" db:"emoticonID"`
	Rating      int    `json:"rating" db:"rating"`
	TopRole     string `json:"topRole" db:"toprole"`
	Wins        int    `json:"wins" db:"wins"`
	Losses      int    `json:"losses" db:"losses"`
	Games       int    `json:"games" db:"games"`
}

type LeaderboardPaging struct {
	StartRank  int               `json:"startRank"`
	PageSize   int               `json:"pageSize"`
	Region     LeaderboardRegion `json:"specificRegion"`
	TotalItems int               `json:"totalItems"`
}

type LeaderboardResponse struct {
	Players []LeaderboardPlayer
	Paging  LeaderboardPaging
	Region  LeaderboardRegion `json:"specificRegion"`
}

type LeaderboardRegion string

const (
	World        LeaderboardRegion = "World"
	Europe       LeaderboardRegion = "Europe"
	NorthAmerica LeaderboardRegion = "NorthAmerica"
	SouthAmerica LeaderboardRegion = "SouthAmerica"
	Asia         LeaderboardRegion = "Asia"
)

func GetLeaderboardPage(ctx context.Context, from int, size int, region LeaderboardRegion) (*LeaderboardResponse, error) {
	leaderboardUrl := fmt.Sprintf("https://prometheus.odysseyinteractive.gg/api/v1/ranked/leaderboard/players?startRank=%d&pageSize=%d", from, size)
	if region != World {
		leaderboardUrl += "&specificRegion=" + string(region)
	}
	resp, err := http.Get(leaderboardUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jsonStr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response LeaderboardResponse
	err = json.Unmarshal(jsonStr, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func CopyLeaderboardsRoutine() {
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	regions := []LeaderboardRegion{Europe, NorthAmerica} //other regions ignored for now
	for {
		t := time.Now()
		for _, region := range regions {
			for i := 1; i < 10000; i += 25 {
				response, err := GetLeaderboardPage(ctx, i, 25, region)
				if err != nil {
					continue
				}
				db := db.GetInstance()
				tx, err := db.Beginx()
				if err != nil {
					tx.Rollback()
					continue
				}
				for _, player := range response.Players {
					_, err = tx.Exec("DELETE FROM leaderboard WHERE playerid=?", player.PlayerID)
					if err != nil {
						tx.Rollback()
						continue
					}
					_, err = tx.NamedExec("INSERT INTO leaderboard (username,playerID,logoID,titleID,nameplateID,emoticonID,rating,toprole,wins,losses,games) VALUES (:username,:playerID,:logoID,:titleID,:nameplateID,:emoticonID,:rating,:toprole,:wins,:losses,:games)", player)
					if err != nil {
						tx.Rollback()
						continue
					}
				}
				err = tx.Commit()
				if err != nil {
					err2 := tx.Rollback()
					if err2 != nil {
						continue
					}
					continue
				}
				time.Sleep(time.Millisecond * 100)
			}
		}
		log.WithFields(log.Fields{
			string(static.UUIDKey): ctx.Value(static.UUIDKey),
			"timeelapsed":          time.Since(t),
		}).Debug("finished sync leaderboard")
	}
}
