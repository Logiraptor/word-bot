import * as React from "react";
import { Tile } from "../models/core";

export interface Props {
    mini?: boolean;
    Tiles: Tile[];
    onChange(tiles: Tile[]);
    onMove(dRow: number, dCol: number);
    onFlip();
    onDelete();
    onChangePlayer(player: number);
    onSubmit();
    score?: number;
    player: number;
}

export class RackInput extends React.Component<Props> {
    constructor(props) {
        super(props);
    }

    renderTile = (tile: Tile, i: number) => {
        return (
            <TileView
                key={i}
                tile={tile}
                onClick={() => {
                    let tiles = [ ...this.props.Tiles ];
                    tiles[i] = { ...tile, Blank: !tile.Blank };
                    this.props.onChange(tiles);
                }}
            />
        );
    };

    render() {
        return (
            <div
                tabIndex={0}
                className={"rack" + (this.props.mini ? " mini" : "")}
                onKeyDown={(event) => {
                    let newTile: Tile;
                    if (event.keyCode <= 90 && event.keyCode >= 65) {
                        let letter = event.key;
                        if (event.shiftKey) {
                            letter = event.key.toLowerCase();
                        }
                        newTile = {
                            Blank: false,
                            Letter: letter,
                            Value: 0,
                            Bonus: "",
                        };
                        this.props.onChange([ ...this.props.Tiles, newTile ]);
                        return;
                    }
                    if (event.keyCode <= 57 && event.keyCode >= 48) {
                        this.props.onChangePlayer(event.keyCode - 48);
                        return;
                    }
                    switch (event.keyCode) {
                        case 37:
                            this.props.onMove(0, -1);
                            break;
                        case 38:
                            this.props.onMove(-1, 0);
                            break;
                        case 39:
                            this.props.onMove(0, 1);
                            break;
                        case 40:
                            this.props.onMove(1, 0);
                            break;
                        case 46:
                            this.props.onDelete();
                            break;
                        case 8:
                            this.props.onChange(this.props.Tiles.slice(0, -1));
                            break;
                        case 32:
                            newTile = {
                                Blank: true,
                                Letter: "",
                                Value: 0,
                                Bonus: "",
                            };
                            this.props.onChange([ ...this.props.Tiles, newTile ]);
                            break;
                        case 16:
                            this.props.onFlip();
                            break;
                        case 13:
                            this.props.onSubmit();
                            break;
                        default:
                            console.log(event.keyCode);
                    }
                }}
            >
                <span className={"move-score" + ` player-${this.props.player}`}>{this.props.score}</span>
                {this.props.Tiles.map(this.renderTile)}
            </div>
        );
    }
}

export function TileView({ tile, onClick }: { tile: Tile; onClick: () => void }) {
    let letter: string;
    if (tile.Blank && !tile.Letter) {
        letter = " ";
    } else {
        letter = tile.Letter;
    }

    let hasTile = tile.Value !== -1;
    if (hasTile) {
        return (
            <span onClick={onClick} className={`space tile ${tile.Blank ? " blank" : ""}`}>
                <span className="letter">{letter}</span>
                <span className="score">{tile.Value}</span>{" "}
            </span>
        );
    }

    return <span onClick={onClick} className={`space ${tile.Bonus}`} />;
}
