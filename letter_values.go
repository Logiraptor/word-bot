package main

var letters = struct {
	A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V, W, X, Y, Z int
}{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25,
}

// From https://scrabble.hasbro.com/en-us/faq
// (1 point)-A, E, I, O, U, L, N, S, T, R
// (2 points)-D, G
// (3 points)-B, C, M, P
// (4 points)-F, H, V, W, Y
// (5 points)-K
// (8 points)- J, X
// (10 points)-Q, Z

var letterValues [26]Score

func init() {
	letterValues[letters.A] = 1
	letterValues[letters.B] = 3
	letterValues[letters.C] = 3
	letterValues[letters.D] = 2
	letterValues[letters.E] = 1
	letterValues[letters.F] = 4
	letterValues[letters.G] = 2
	letterValues[letters.H] = 4
	letterValues[letters.I] = 1
	letterValues[letters.J] = 8
	letterValues[letters.K] = 5
	letterValues[letters.L] = 1
	letterValues[letters.M] = 3
	letterValues[letters.N] = 1
	letterValues[letters.O] = 1
	letterValues[letters.P] = 3
	letterValues[letters.Q] = 10
	letterValues[letters.R] = 1
	letterValues[letters.S] = 1
	letterValues[letters.T] = 1
	letterValues[letters.U] = 1
	letterValues[letters.V] = 4
	letterValues[letters.W] = 4
	letterValues[letters.X] = 8
	letterValues[letters.Y] = 4
	letterValues[letters.Z] = 10
}
