export type Board = (Tile | null)[][];

export enum TileFlag {
    NextAIMove,
}

export interface Tile {
    Letter: string;
    Blank: boolean;
    Value: number;
    Bonus: 'DW' | 'TW' | 'TL' | 'DL' | '';
    Flags: TileFlag[] | null;
}

export interface Move {
    tiles: Tile[];
    row: number;
    col: number;
    player: number;
    direction: "horizontal" | "vertical";
    score?: number;
    valid?: boolean;
}

export interface MoveRequest {
    moves: Move[];
    rack: Tile[];
}

export interface RenderedBoard {
    Board: Board;
    Scores: number[];
}
