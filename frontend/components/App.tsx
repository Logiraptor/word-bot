import * as React from 'react'

import {Move, Tile} from '../models/core'
import {BoardView} from './BoardView'
import {RackInput} from './RackInput'
import {AppStore, DefaultState} from '../models/store'
import {GameState} from '../models/gamestate'

export interface State {
    store: AppStore;
}

interface Props {
    gameState: GameState;
}

export class App extends React.Component<Props, State> {
    constructor(props: Props) {
        super(props)
        this.state = {
            store: DefaultState,
        }
    }

    componentDidMount() {
        this.props.gameState.subscribe((s) => {
            this.setState({store: s})
        })
    }

    render() {
        const {store} = this.state

        let counts: { [x: string]: { tile: Tile; count: number } } = {}
        store.remainingTiles.forEach((tile) => {
            const name = tile.Blank ? 'blank' : tile.Letter

            if (!(name in counts)) {
                counts[name] = {
                    tile,
                    count: 0,
                }
            }

            counts[name].count++
        })

        return (
            <div className="app">
                <BoardView tiles={store.board}/>
                <div className="right-panel">
                    <MovePanel gameState={this.props.gameState} moves={store.moves}/>
                    <ScoreBoard moves={store.moves} scores={store.moves.map((x) => x.score)}/>
                    <div className="player-rack">
                        <RackInput
                            Tiles={store.rack}
                            onChange={(tiles) => {
                                this.props.gameState.setRack(tiles)
                            }}
                            onMove={(row, col) => {
                            }}
                            onFlip={() => {
                            }}
                            onDelete={() => {
                            }}
                            onChangePlayer={(player) => {
                            }}
                            player={undefined}
                            onSubmit={() => {
                                let newRack = removeFromRack(store.rack, store.play.tiles)
                                this.props.gameState.addMove({...store.play, player: getOtherPlayer(store.moves)})
                                this.props.gameState.setRack(newRack)
                            }}
                        />
                    </div>
                </div>
            </div>
        )
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
                <td>{score.player}</td>
                <td>{score.score}</td>
            </tr>
        )
    }

    render() {
        let scores: Score[] = []
        this.props.moves.forEach((element, i) => {
            if (!scores[element.player]) {
                scores[element.player] = {player: element.player, score: 0}
            }
            scores[element.player].score += this.props.scores[i]
        })

        return (
            <div className="score-board">
                <table>
                    <tbody>{scores.map(this.renderRow)}</tbody>
                </table>
            </div>
        )
    }
}

interface MovePanelProps {
    gameState: GameState;
    moves: Move[];
}

class MovePanel extends React.Component<MovePanelProps> {
    renderMove = (move: Move, i: number) => {
        const changeMove = (f: (move: Move) => void) => {
            const newMove = {...move}
            f(newMove)
            this.props.gameState.updateMove(i, newMove)
        }

        return (
            <div className={move.valid === false ? 'move-panel-move error-icon' : 'move-panel-move'} key={i}>
                <RackInput
                    score={move.score}
                    mini
                    Tiles={move.tiles}
                    onChange={(tiles) => {
                        changeMove((m) => (m.tiles = tiles))
                    }}
                    onMove={(row, col) => {
                        changeMove((m) => {
                            m.row = Math.max(0, Math.min(14, m.row + row))
                            m.col = Math.max(0, Math.min(14, m.col + col))
                        })
                    }}
                    onFlip={() => {
                        changeMove((m) => {
                            m.direction = m.direction === 'horizontal' ? 'vertical' : 'horizontal'
                        })
                    }}
                    onDelete={() => {
                        this.props.gameState.removeMove(i)
                    }}
                    player={move.player}
                    onChangePlayer={(player) => {
                        changeMove((m) => {
                            m.player = player
                        })
                    }}
                    onSubmit={() => {
                        this.props.gameState.addMove(newMove(this.props.moves))
                    }}
                />
            </div>
        )
    }

    render() {
        return (
            <div className="move-panel">
                <div className="scroll">{this.props.moves.map(this.renderMove)}</div>
            </div>
        )
    }
}

function newMove(moves: Move[]): Move {
    return {
        tiles: [],
        row: 0,
        col: 0,
        player: getOtherPlayer(moves),
        direction: 'horizontal',
    }
}

function getOtherPlayer(moves: Move[]) {
    let lastPlayer = undefined
    if (moves.length > 0) {
        let lastMove = moves[moves.length - 1]
        lastPlayer = lastMove.player === 1 ? 2 : 1
    }
    return lastPlayer
}

function removeFromRack(oldRack: Tile[], tiles: Tile[]): Tile[] {
    let rack = [...oldRack]
    tiles.forEach((tile) => {
        let i = rack.findIndex((candidate) => {
            if (candidate.Blank && tile.Blank) {
                return true
            }
            return candidate.Letter == tile.Letter
        })
        rack.splice(i, 1)
    })
    return rack
}
