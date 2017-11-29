import { Tile, Move, Board } from "./core";

export interface SetRack {
    type: "setrack";
    value: Tile[];
}

export interface UpdateMove {
    type: "updatemove";
    value: Move;
    index: number;
}

export interface DeleteMove {
    type: "deletemove";
    index: number;
}

export interface AddMove {
    type: "addmove";
    value: Move;
}

export interface ReceiveRender {
    type: "receiverender";
    board: Board;
    scores: number[];
}

export interface ReceiveValidations {
    type: "receivevalidations";
    validations: boolean[];
}

export type Action = SetRack | UpdateMove | DeleteMove | AddMove | ReceiveRender | ReceiveValidations;

export function setRack(value: Tile[]): SetRack {
    return { type: "setrack", value };
}

export function addMove(value: Move): AddMove {
    return { type: "addmove", value };
}

export function deleteMove(index: number): DeleteMove {
    return { type: "deletemove", index };
}

export function updateMove(value: Move, index: number): UpdateMove {
    return { type: "updatemove", value, index };
}

export function receiveRender(board: Board, scores: number[]): ReceiveRender {
    return { type: "receiverender", board, scores };
}

export function receiveValidations(validations: boolean[]): ReceiveValidations {
    return { type: "receivevalidations", validations };
}
