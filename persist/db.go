package persist

import (
	"encoding/json"

	"github.com/Logiraptor/word-bot/core"
	"github.com/jinzhu/gorm"
)

type DB struct {
	db *gorm.DB
}

func NewDB(filename string) (*DB, error) {
	db, err := gorm.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(Game{}, Move{})

	return &DB{
		db: db,
	}, nil
}

type Game struct {
	gorm.Model
	Moves []Move
}

func (g *Game) AddMove(player string, move core.ScoredMove) {
	buf, _ := json.Marshal(move.Word)

	g.Moves = append(g.Moves, Move{
		Tiles:  string(buf),
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
	Row, Col int
	Dir      core.Direction
	Score    core.Score

	Player string
	Game   Game
	GameID uint
}

func (db *DB) SaveGame(g Game) error {
	return db.db.Create(&g).Error
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
			(SELECT games.id as game_id, p1.player as p1, p2.player as p2, p1.score > p2.score as win
			FROM games
			INNER JOIN game_player_scores AS p1 ON p1.game_id = games.id
			INNER JOIN game_player_scores AS p2 ON p2.game_id = games.id AND p1.player < p2.player)

	SELECT REPLACE(p1, X'0A', ''), REPLACE(p2, X'0A', ''), SUM(win), COUNT(win), 100 * (CAST(SUM(win) AS float) / CAST(COUNT(win) AS float)) as winrate
	FROM matchups
	GROUP BY p1, p2
	ORDER BY winrate

	`

}
