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
ORDER BY p2;

DELETE FROM leave_weights;


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

INSERT INTO leave_weights (leave, weight) SELECT moves.leave, AVG(win)
	FROM matchups
	INNER JOIN moves ON moves.game_id = matchups.game_id
	GROUP BY moves.leave;

