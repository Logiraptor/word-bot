import * as React from "react";

export interface Tile {
    Letter: string;
    Blank: boolean;
    Value: number;
    Bonus: string;
}

export interface Props {
    mini?: boolean;
    Tiles: Tile[];
    onChange(tiles: Tile[]);
    onMove(dRow: number, dCol: number);
    onFlip();
}

export class RackInput extends React.Component<Props> {
    constructor(props) {
        super(props);
    }

    renderTile = (tile: Tile, i: number) => {
        return <TileView key={i} tile={tile} />;
    };

    render() {
        return (
            <div
                tabIndex={0}
                className={"rack" + (this.props.mini ? " mini" : "")}
                onKeyDown={(event) => {
                    let newTile: Tile;
                    if (
                        event.keyCode <= 90 &&
                        event.keyCode >= 65 &&
                        !(event.altKey || event.shiftKey || event.ctrlKey || event.metaKey)
                    ) {
                        newTile = {
                            Blank: false,
                            Letter: event.key,
                            Value: 0,
                            Bonus: "",
                        };
                        this.props.onChange([ ...this.props.Tiles, newTile ]);
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
                        default:
                            console.log(event.keyCode);
                    }
                }}
            >
                {this.props.Tiles.map(this.renderTile)}
            </div>
        );
    }
}

export function TileView({ tile }: { tile: Tile }) {
    let letter: string;
    if (tile.Blank && !tile.Letter) {
        letter = " ";
    } else {
        letter = tile.Letter;
    }

    let hasTile = tile.Value !== -1;
    if (hasTile) {
        return (
            <span className={`space tile ${tile.Blank ? " blank" : ""}`}>
                <span className="letter">{letter}</span>
                <span className="score">{tile.Value}</span>{" "}
            </span>
        );
    }

    return <span className={`space ${tile.Bonus}`} />;
}
