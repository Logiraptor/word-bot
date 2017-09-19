import * as React from "react";
import * as ReactDOM from "react-dom";
import { RackInput, Tile } from "./rack";
import { Board, BoardView } from "./board";
import "./index.scss";

interface Move {
    tiles: Tile[];
    row: number;
    col: number;
    direction: "horizontal" | "vertical";
}

interface MoveRequest {
    moves: Move[];
    rack: Tile[];
}

interface State {
    moves: Move[];
    scores: number[];
    rack: Tile[];
    board: Board;
}

interface RenderedBoard {
    Board: Board;
    Scores: number[];
}

async function render(req: MoveRequest): Promise<RenderedBoard> {
    return (await fetch("/render", {
        method: "POST",
        body: JSON.stringify(req),
    }).then((x) => x.json())) as RenderedBoard;
}

async function play(req: MoveRequest): Promise<Move> {
    return await fetch("/play", {
        method: "POST",
        body: JSON.stringify(req),
    }).then((x) => x.json());
}

class App extends React.Component<{}, State> {
    constructor() {
        super();

        let movesString = localStorage.getItem("moves");
        if (!movesString) {
            movesString = "[]";
        }

        let moves: Move[];
        try {
            moves = JSON.parse(movesString);
        } catch (e) {
            moves = [];
        }

        this.state = {
            moves: moves,
            rack: [],
            board: Array(15).map(() => Array(15).map(() => null)),
            scores: [],
        };

        let resp = render({ moves: moves, rack: this.state.rack }).then((resp) => {
            console.log(resp.Board);
            this.setState({
                board: resp.Board,
                scores: resp.Scores,
            });
        });
    }

    updateMoves = async (state: Partial<State>) => {
        this.setState({
            moves: state.moves,
        });

        localStorage.setItem("moves", JSON.stringify(state.moves));

        let resp = await render({ moves: state.moves, rack: this.state.rack });
        console.log(resp.Board);
        this.setState({
            board: resp.Board,
            scores: resp.Scores,
        });
    };

    renderMove = (move: Move, i: number) => {
        const changeMove = (f: (move: Move) => void) => {
            let moves = [ ...this.state.moves ];
            moves[i] = { ...move };
            f(moves[i]);
            this.updateMoves({
                moves: moves,
            });
        };

        return (
            <form key={i} className="form" onSubmit={(e) => e.preventDefault()}>
                <button
                    className="btn primary"
                    onClick={() => {
                        let moves = [ ...this.state.moves ];
                        moves.splice(i, 1);
                        this.updateMoves({ moves: moves });
                    }}
                >
                    Remove
                </button>
                {this.state.scores[i]} Points
                <RackInput
                    mini
                    Tiles={move.tiles}
                    onChange={(tiles) => {
                        changeMove((m) => (m.tiles = tiles));
                    }}
                    onMove={(row, col) => {
                        changeMove((m) => {
                            m.row += row;
                            m.col += col;
                        });
                    }}
                    onFlip={() => {
                        changeMove((m) => {
                            m.direction = m.direction === "horizontal" ? "vertical" : "horizontal";
                        });
                    }}
                />
            </form>
        );
    };

    render() {
        return (
            <div>
                <div className="panel">
                    <BoardView tiles={this.state.board} />
                </div>

                <div className="panel">
                    <h1>Scrabble</h1>

                    <div className="scroll">{this.state.moves.map(this.renderMove)}</div>

                    <button
                        className="btn primary"
                        onClick={() => {
                            this.setState({
                                moves: [
                                    ...this.state.moves,
                                    {
                                        tiles: [],
                                        row: 0,
                                        col: 0,
                                        direction: "horizontal",
                                    },
                                ],
                            });
                        }}
                    >
                        Add Move
                    </button>

                    <hr />

                    <RackInput
                        Tiles={this.state.rack}
                        onChange={(tiles) => {
                            this.setState({
                                rack: tiles,
                            });
                        }}
                        onMove={(row, col) => {}}
                        onFlip={() => {}}
                    />

                    <hr />

                    <button
                        className="btn primary"
                        onClick={async () => {
                            let resp: Move = await play({
                                moves: this.state.moves,
                                rack: this.state.rack,
                            });

                            this.updateMoves({
                                moves: [ ...this.state.moves, resp ],
                            });
                        }}
                    >
                        Play For Me
                    </button>
                </div>
            </div>
        );
    }
}

let app = document.createElement("div");
document.body.appendChild(app);
ReactDOM.render(<App />, app);
