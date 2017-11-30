import { Tile, Move, Board } from "./core";

export interface SetRack {
    changesBoard: true;
    type: "setrack";
    value: Tile[];
}

export interface UpdateMove {
    changesBoard: true;
    type: "updatemove";
    value: Move;
    index: number;
}

export interface DeleteMove {
    changesBoard: true;
    type: "deletemove";
    index: number;
}

export interface AddMove {
    changesBoard: true;
    type: "addmove";
    value: Move;
}

export interface ReceiveRender {
    changesBoard: false;
    type: "receiverender";
    board: Board;
    scores: number[];
}

export interface ReceiveValidations {
    changesBoard: false;
    type: "receivevalidations";
    validations: boolean[];
}

export interface ReceivePlay {
    changesBoard: true;
    type: "receiveplay";
    play: Move;
}

export interface ReceiveRemainingTiles {
    changesBoard: false;
    type: "receiveremainingtiles";
    tiles: Tile[];
}

export type Action =
    | SetRack
    | UpdateMove
    | DeleteMove
    | AddMove
    | ReceiveRender
    | ReceiveValidations
    | ReceivePlay
    | ReceiveRemainingTiles;

export function setRack(value: Tile[]): SetRack {
    return { type: "setrack", value, changesBoard: true };
}

export function addMove(value: Move): AddMove {
    return { type: "addmove", value, changesBoard: true };
}

export function deleteMove(index: number): DeleteMove {
    return { type: "deletemove", index, changesBoard: true };
}

export function updateMove(value: Move, index: number): UpdateMove {
    return { type: "updatemove", value, index, changesBoard: true };
}

export function receiveRender(board: Board, scores: number[]): ReceiveRender {
    return { type: "receiverender", board, scores, changesBoard: false };
}

export function receiveValidations(validations: boolean[]): ReceiveValidations {
    return { type: "receivevalidations", validations, changesBoard: false };
}

export function receivePlay(play: Move): ReceivePlay {
    return { type: "receiveplay", play, changesBoard: true };
}

export function receiveRemainingTiles(tiles: Tile[]): ReceiveRemainingTiles {
    return { type: "receiveremainingtiles", changesBoard: false, tiles };
}
