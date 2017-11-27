import * as React from "react";

import { Board, Move, Tile } from "../models/core";
import { GameService, StorageService } from "../services/game";
import { BoardView } from "./BoardView";
import { RackInput } from "./RackInput";
import { Store, Dispatch } from "redux";
import { AppStore, Action, updateMove, deleteMove, addMove, setRack, DefaultState } from "../models/store";

export interface State {
    store: AppStore;
}

interface Props {
    store: Store<AppStore>;
    gameService: GameService;
    storage: StorageService;
}

export class App extends React.Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = {
            store: DefaultState,
        };
    }

    componentDidMount() {
        this.props.store.subscribe(() => {
            this.setState({
                store: this.props.store.getState(),
            });
        });
    }

    render() {
        const { store } = this.state;
        return (
            <div>
                <BoardView tiles={store.board} />
                <div className="right-panel">
                    <MovePanel
                        dispatch={this.props.store.dispatch}
                        moves={store.moves}
                        scores={store.scores}
                        service={this.props.gameService}
                    />
                    <ScoreBoard moves={store.moves} scores={store.scores} />
                    <div className="player-rack">
                        <RackInput
                            Tiles={store.rack}
                            onChange={(tiles) => {
                                this.props.store.dispatch(setRack(tiles));
                                // this.setState({
                                //     rack: tiles,
                                // });
                            }}
                            onMove={(row, col) => {}}
                            onFlip={() => {}}
                            onDelete={() => {}}
                            onChangePlayer={(player) => {}}
                            player={undefined}
                            onSubmit={async () => {
                                let resp: Move = await this.props.gameService.play({
                                    moves: store.moves,
                                    rack: store.rack,
                                });

                                let newRack = removeFromRack(store.rack, resp.tiles);
                                this.props.store.dispatch(addMove({ ...resp, player: getOtherPlayer(store.moves) }));
                                this.props.store.dispatch(setRack(newRack));
                                // this.setState({
                                //     moves: [ ...store.moves,  ],
                                //     rack: newRack,
                                // });
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

interface MovePanelProps {
    dispatch: Dispatch<Action>;
    moves: Move[];
    scores: number[];
    service: GameService;
}

class MovePanel extends React.Component<MovePanelProps, { valid: boolean[] }> {
    state = {
        valid: [],
    };

    componentDidMount() {
        this.revalidate();
    }

    componentDidUpdate(prevProps: MovePanelProps) {
        if (this.props.moves != prevProps.moves) {
            this.revalidate();
        }
    }

    async revalidate() {
        let valid = await this.props.service.validate({
            moves: this.props.moves,
            rack: [],
        });
        this.setState({ valid });
    }

    renderMove = (move: Move, i: number) => {
        const changeMove = (f: (move: Move) => void) => {
            // let moves = [ ...this.props.moves ];
            // moves[i] = { ...move };
            // f(moves[i]);
            // this.props.changeMoves(moves);

            const newMove = { ...move };
            f(newMove);
            this.props.dispatch(updateMove(newMove, i));
        };

        return (
            <div className={!this.state.valid[i] ? "move-panel-move error-icon" : "move-panel-move"} key={i}>
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
                        this.props.dispatch(deleteMove(i));
                        // let newMoves = [ ...this.props.moves ];
                        // newMoves.splice(i, 1);

                        // this.props.changeMoves(newMoves);
                    }}
                    player={move.player}
                    onChangePlayer={(player) => {
                        changeMove((m) => {
                            m.player = player;
                        });
                    }}
                    onSubmit={() => {
                        this.props.dispatch(addMove(newMove(this.props.moves)));
                        // this.props.changeMoves([ ...this.props.moves, newMove(this.props.moves) ]);
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
