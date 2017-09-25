import { Move, Tile, Board } from "../models/core";
import * as React from "react";
import { RackInput } from "./RackInput";
import { BoardView } from "./BoardView";
import { GameService, StorageService } from "../services/game";

export interface State {
    moves: Move[];
    scores: number[];
    rack: Tile[];
    board: Board;
}

export class App extends React.Component<{ gameService: GameService; storage: StorageService }, State> {
    constructor(props) {
        super();

        this.state = {
            moves: [],
            rack: [],
            board: Array(15).map(() => Array(15).map(() => null)),
            scores: [],
        };
    }

    async componentDidMount() {
        let moves = await this.props.storage.load();
        this.setState({ moves });
    }

    componentDidUpdate(_, prevState: State) {
        if (this.state.moves !== prevState.moves) {
            this.props.storage.save(this.state.moves);
            this.props.gameService.render({ moves: this.state.moves, rack: this.state.rack }).then((resp) => {
                this.setState({
                    board: resp.Board,
                    scores: resp.Scores,
                });
            });

            if (this.state.moves.length === 0) {
                this.setState({ moves: [ newMove([]) ] });
            }
        }
    }

    render() {
        return (
            <div>
                <BoardView tiles={this.state.board} />
                <div className="right-panel">
                    <MovePanel
                        moves={this.state.moves}
                        scores={this.state.scores}
                        changeMoves={(moves) => {
                            this.setState({ moves });
                        }}
                    />
                    <ScoreBoard moves={this.state.moves} scores={this.state.scores} />
                    <div className="player-rack">
                        <RackInput
                            Tiles={this.state.rack}
                            onChange={(tiles) => {
                                this.setState({
                                    rack: tiles,
                                });
                            }}
                            onMove={(row, col) => {}}
                            onFlip={() => {}}
                            onDelete={() => {}}
                            onChangePlayer={(player) => {}}
                            player={undefined}
                            onSubmit={async () => {
                                let resp: Move = await this.props.gameService.play({
                                    moves: this.state.moves,
                                    rack: this.state.rack,
                                });

                                let newRack = removeFromRack(this.state.rack, resp.tiles);

                                this.setState({
                                    moves: [
                                        ...this.state.moves,
                                        { ...resp, player: getOtherPlayer(this.state.moves) },
                                    ],
                                    rack: newRack,
                                });
                            }}
                        />
                    </div>
                </div>
            </div>
        );
    }
}

interface Score {
    player: number;
    score: number;
}

class ScoreBoard extends React.Component<{ moves: Move[]; scores: number[] }> {
    renderRow = (score: Score, i: number) => {
        return (
            <tr className={`player-${score.player}`} key={i}>
                <td>{score.player}</td> <td>{score.score}</td>
            </tr>
        );
    };

    render() {
        let scores: Score[] = [];
        this.props.moves.forEach((element, i) => {
            if (!scores[element.player]) {
                scores[element.player] = { player: element.player, score: 0 };
            }
            scores[element.player].score += this.props.scores[i];
        });

        return (
            <div className="score-board">
                <table>
                    <tbody>{scores.map(this.renderRow)}</tbody>
                </table>
            </div>
        );
    }
}

class MovePanel extends React.Component<{ moves: Move[]; scores: number[]; changeMoves: (m: Move[]) => void }> {
    renderMove = (move: Move, i: number) => {
        const changeMove = (f: (move: Move) => void) => {
            let moves = [ ...this.props.moves ];
            moves[i] = { ...move };
            f(moves[i]);
            this.props.changeMoves(moves);
        };

        return (
            <div key={i}>
                <RackInput
                    score={this.props.scores[i]}
                    mini
                    Tiles={move.tiles}
                    onChange={(tiles) => {
                        changeMove((m) => (m.tiles = tiles));
                    }}
                    onMove={(row, col) => {
                        changeMove((m) => {
                            m.row = Math.max(0, Math.min(14, m.row + row));
                            m.col = Math.max(0, Math.min(14, m.col + col));
                        });
                    }}
                    onFlip={() => {
                        changeMove((m) => {
                            m.direction = m.direction === "horizontal" ? "vertical" : "horizontal";
                        });
                    }}
                    onDelete={() => {
                        let newMoves = [ ...this.props.moves ];
                        newMoves.splice(i, 1);
                        this.props.changeMoves(newMoves);
                    }}
                    player={move.player}
                    onChangePlayer={(player) => {
                        changeMove((m) => {
                            m.player = player;
                        });
                    }}
                    onSubmit={() => {
                        this.props.changeMoves([ ...this.props.moves, newMove(this.props.moves) ]);
                    }}
                />
            </div>
        );
    };

    render() {
        return (
            <div className="move-panel">
                <div className="scroll">{this.props.moves.map(this.renderMove)}</div>
            </div>
        );
    }
}

function newMove(moves: Move[]): Move {
    return {
        tiles: [],
        row: 0,
        col: 0,
        player: getOtherPlayer(moves),
        direction: "horizontal",
    };
}

function getOtherPlayer(moves: Move[]) {
    let lastPlayer = undefined;
    if (moves.length > 0) {
        let lastMove = moves[moves.length - 1];
        lastPlayer = lastMove.player === 1 ? 2 : 1;
    }
    return lastPlayer;
}

function removeFromRack(oldRack: Tile[], tiles: Tile[]): Tile[] {
    let rack = [ ...oldRack ];
    tiles.forEach((tile) => {
        let i = rack.findIndex((candidate) => {
            if (candidate.Blank && tile.Blank) {
                return true;
            }
            return candidate.Letter == tile.Letter;
        });
        rack.splice(i, 1);
    });
    return rack;
}
