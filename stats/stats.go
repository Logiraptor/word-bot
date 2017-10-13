package stats

import (
	"math"
)

func StatisticalSignificance(wins, losses, draws int) (winrate float64, chiSquare float64, significant bool) {
	var (
		W  = float64(wins)
		D  = float64(draws)
		L  = float64(losses)
		N  = W + D + L
		t1 = 2 * (math.Pow(W, 2) + math.Pow(L, 2))
		t2 = math.Pow(D+math.Sqrt(t1), 2)
		T  = sign(W, L) * ((t2 / N) - N)

		p = (W + (D / 2)) / N
	)
	return p, T, T > 3.841
}

func sign(a, b float64) float64 {
	if a > b {
		return 1
	}
	return -1
}
