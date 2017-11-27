import { createStore } from "redux";
import { Tile, Move, Board } from "./core";

export interface AppStore {
    moves: Move[];
    scores: number[];
    rack: Tile[];
    board: Board;
}

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

export type Action = SetRack | UpdateMove | DeleteMove | AddMove;

export function setRack(value: Tile[]): SetRack {
    return {
        type: "setrack",
        value,
    };
}

export function addMove(value: Move): AddMove {
    return {
        type: "addmove",
        value,
    };
}

export function deleteMove(index: number): DeleteMove {
    return {
        type: "deletemove",
        index,
    };
}

export function updateMove(value: Move, index: number): UpdateMove {
    return {
        type: "updatemove",
        value,
        index,
    };
}

export const DefaultState = {
    moves: [],
    rack: [],
    board: Array(15).map(() => Array(15).map(() => null)),
    scores: [],
};

export function reducer(state: AppStore | undefined, action: Action): AppStore {
    if (!state) {
        return DefaultState;
    }
    switch (action.type) {
        case "addmove":
            return state;
        case "deletemove":
            return state;
        case "setrack":
            return state;
        case "updatemove":
            return state;
    }
}

// TODO: persist moves on change
// TODO: re-render board when moves change
// TODO: make sure there's always at least one move
// TODO: load moves from localstorage on app boot
