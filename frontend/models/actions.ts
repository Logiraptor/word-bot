import { Tile, Move, Board } from "./core";

export interface SetRack {
    isUserInput: true;
    type: "setrack";
    value: Tile[];
}

export interface UpdateMove {
    isUserInput: true;
    type: "updatemove";
    value: Move;
    index: number;
}

export interface DeleteMove {
    isUserInput: true;
    type: "deletemove";
    index: number;
}

export interface AddMove {
    isUserInput: true;
    type: "addmove";
    value: Move;
}

export interface ReceiveRender {
    isUserInput: false;
    type: "receiverender";
    board: Board;
    scores: number[];
}

export interface ReceiveValidations {
    isUserInput: false;
    type: "receivevalidations";
    validations: boolean[];
}

export interface ReceivePlay {
    isUserInput: false;
    type: "receiveplay";
    play: Move;
}

export type Action = SetRack | UpdateMove | DeleteMove | AddMove | ReceiveRender | ReceiveValidations | ReceivePlay;

export function setRack(value: Tile[]): SetRack {
    return { type: "setrack", value, isUserInput: true };
}

export function addMove(value: Move): AddMove {
    return { type: "addmove", value, isUserInput: true };
}

export function deleteMove(index: number): DeleteMove {
    return { type: "deletemove", index, isUserInput: true };
}

export function updateMove(value: Move, index: number): UpdateMove {
    return { type: "updatemove", value, index, isUserInput: true };
}

export function receiveRender(board: Board, scores: number[]): ReceiveRender {
    return { type: "receiverender", board, scores, isUserInput: false };
}

export function receiveValidations(validations: boolean[]): ReceiveValidations {
    return { type: "receivevalidations", validations, isUserInput: false };
}

export function receivePlay(play: Move): ReceivePlay {
    return { type: "receiveplay", play, isUserInput: false };
}
