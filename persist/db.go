package persist

import (
	"github.com/Logiraptor/word-bot/core"
	"github.com/jinzhu/gorm"
)

type DB struct {
	DB *gorm.DB
}

func NewDB(filename string) (*DB, error) {
	db, err := gorm.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	return NewDBConn(db)
}

func NewDBConn(db *gorm.DB) (*DB, error) {
	err := db.AutoMigrate(Game{}, Move{}, LeaveWeight{}).Error
	if err != nil {
		return nil, err
	}

	return &DB{
		DB: db,
	}, nil
}

type LeaveWeight struct {
	Leave  string
	Weight float64
}

type Game struct {
	gorm.Model
	Moves []Move
}

func (g *Game) AddMove(player string, leave []core.Tile, move core.ScoredMove) {
	g.Moves = append(g.Moves, Move{
		Tiles:  core.Tiles2String(move.Word),
		Leave:  core.Tiles2String(leave),
		Row:    move.Row,
		Col:    move.Col,
		Player: player,
		Dir:    move.Direction,
		Score:  move.Score,
	})
}

type Move struct {
	ID       uint `gorm:"primary_key"`
	Tiles    string
	Leave    string
	Row, Col int
	Dir      core.Direction
	Score    core.Score

	Player string
	Game   Game
	GameID uint
}

func (db *DB) SaveGame(g Game) error {
	return db.DB.Create(&g).Error
}

func (db *DB) LoadLeaveWeights() ([]LeaveWeight, error) {
	var weights []LeaveWeight
	err := db.DB.Find(&weights).Error
	if err != nil {
		return nil, err
	}
	return weights, nil
}

type Matchup struct {
	Player1, Player2 string
	NumGames         int
}

func (db *DB) GetMatchups(matchups []Matchup) (error) {
	//TODO: this query is returning invalid results.
	//I need a player id concept
	rows, err := db.DB.Raw(`
	WITH
	game_player_scores AS
		(SELECT games.id AS game_id, player, SUM(moves.score) AS score
		FROM games
		INNER JOIN moves ON moves.game_id = games.id
		GROUP BY games.id, moves.player)
	SELECT DISTINCT p1.player AS p1, p2.player AS p2, COUNT(games.id)
		FROM games
		INNER JOIN game_player_scores AS p1 ON p1.game_id = games.id
		INNER JOIN game_player_scores AS p2 ON p2.game_id = games.id
		GROUP BY p1.player, p2.player
	`).Rows()
	if err != nil {
		return err
	}

	for rows.Next() {
		m := Matchup{}
		err = rows.Scan(&m.Player1, &m.Player2, &m.NumGames)
		if err != nil {
			return err
		}

		for i, other := range matchups {
			if other.Player1 == m.Player1 && other.Player2 == m.Player2 {
				matchups[i].NumGames = m.NumGames
				break
			}
		}
	}
	return nil
}

func (db *DB) PrintStats() {
	// score by game, player
	_ = `
	SELECT games.id, player, SUM(moves.score)
	FROM games
	INNER JOIN moves ON moves.game_id = games.id
	GROUP BY games.id, moves.player

	`

	// number of players per game
	_ = `
	SELECT games.id, COUNT(DISTINCT player)
	FROM games
	INNER JOIN moves ON moves.game_id = games.id
	GROUP BY games.id

	`

	// gameid, score1, score2, p1Wins
	_ = `
	WITH game_player_scores
	AS (SELECT games.id AS game_id, player, SUM(moves.score) as score
	FROM games
	INNER JOIN moves ON moves.game_id = games.id
	GROUP BY games.id, moves.player)
	SELECT games.id, p1.score, p2.score, p1.score > p2.score
	FROM games
	INNER JOIN game_player_scores AS p1 ON p1.game_id = games.id
	INNER JOIN game_player_scores AS p2 ON p2.game_id = games.id AND p1.player > p2.player

	`

	// gameid, score1, score2, p1Wins
	_ = `
	WITH game_player_scores
	AS (SELECT games.id AS game_id, player, SUM(moves.score) as score
	FROM games
	INNER JOIN moves ON moves.game_id = games.id
	GROUP BY games.id, moves.player)
	SELECT games.id, p1.score, p2.score, p1.player, p2.player, p1.score > p2.score
	FROM games
	INNER JOIN game_player_scores AS p1 ON p1.game_id = games.id
	INNER JOIN game_player_scores AS p2 ON p2.game_id = games.id AND p1.player > p2.player

	`

	// win percent for each matchup
	_ = `
	WITH
		game_player_scores AS 
			(SELECT games.id AS game_id, player, SUM(moves.score) as score
			FROM games
			INNER JOIN moves ON moves.game_id = games.id
			GROUP BY games.id, moves.player),
		matchups AS
			(SELECT games.id as game_id, p1.player as p1, p2.player as p2, p1.score - p2.score as diff, p1.score > p2.score as win
			FROM games
			INNER JOIN game_player_scores AS p1 ON p1.game_id = games.id
			INNER JOIN game_player_scores AS p2 ON p2.game_id = games.id AND p1.player < p2.player)

	SELECT REPLACE(p1, X'0A', ''), REPLACE(p2, X'0A', ''), SUM(win), COUNT(win), 100 * (CAST(SUM(win) AS float) / CAST(COUNT(win) AS float)) as winrate
	FROM matchups
	GROUP BY p1, p2
	ORDER BY winrate DESC

	`

	// leave + diff + win
	_ = `
	WITH
	game_player_scores AS 
		(SELECT games.id AS game_id, player, SUM(moves.score) AS score
		FROM games
		INNER JOIN moves ON moves.game_id = games.id
		GROUP BY games.id, moves.player),
	matchups AS
		(SELECT games.id AS game_id, p1.player AS p1, p2.player AS p2, p1.score - p2.score AS diff, p1.score > p2.score AS win
		FROM games
		INNER JOIN game_player_scores AS p1 ON p1.game_id = games.id
		INNER JOIN game_player_scores AS p2 ON p2.game_id = games.id AND p1.player < p2.player)
	SELECT moves.leave, diff, win
	FROM matchups
	INNER JOIN moves ON moves.game_id = matchups.game_id
	`

}
