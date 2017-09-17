import * as React from "react";
import { Tile, TileView } from "./rack";

export type Board = Tile[][];

export class BoardView extends React.Component<{ tiles: Board }> {
    renderRow = (row: Tile[], i: number) => {
        let tiles = row.map((tile, i) => {
            return <TileView key={i} tile={tile} />;
        });

        return <div key={i}>{tiles}</div>;
    };

    render() {
        return <div className="board">{this.props.tiles.map(this.renderRow)}</div>;
    }
}
