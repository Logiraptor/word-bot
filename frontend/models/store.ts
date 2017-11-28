import { createStore, Store } from "redux";
import { Tile, Move, Board } from "./core";
import { GameService } from "../services/game";

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

export interface UpdateBoard {
    type: "updateboard";
    board: Board;
    scores: number[];
}

export type Action = SetRack | UpdateMove | DeleteMove | AddMove | UpdateBoard;

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

export function updateBoard(board: Board, scores: number[]): UpdateBoard {
    return {
        type: "updateboard",
        board,
        scores,
    };
}

export const EmptyMove: Move = {
    tiles: [],
    row: 0,
    col: 0,
    player: undefined,
    direction: "horizontal",
};

export const DefaultState = {
    moves: [ EmptyMove ],
    rack: [],
    board: Array(15).map(() => Array(15).map(() => null)),
    scores: [],
};

export function reducer(state: AppStore | undefined, action: Action): AppStore {
    if (!state) {
        return DefaultState;
    }
    state = { ...state };
    switch (action.type) {
        case "addmove":
            state.moves = [ ...state.moves, action.value ];
            return state;
        case "deletemove":
            state.moves = [ ...state.moves ];
            state.moves.splice(action.index, 1);
            if (state.moves.length === 0) {
                state.moves.push(EmptyMove);
            }
            return state;
        case "setrack":
            state.rack = action.value;
            return state;
        case "updatemove":
            state.moves = [ ...state.moves ];
            state.moves[action.index] = action.value;
            return state;
        case "updateboard":
            state.board = action.board;
            state.scores = action.scores;
            return state;
    }
}

export function setupSubscriptions(store: Store<AppStore>, gameService: GameService) {
    let moves = [];
    store.subscribe(async () => {
        const state = store.getState();
        if (state.moves !== moves) {
            moves = state.moves;
            const board = await gameService.render({
                moves: state.moves,
                rack: state.rack,
            });
            console.log(board.Board);
            store.dispatch(updateBoard(board.Board, board.Scores));
        }
    });
}

// TODO: persist moves on change
// TODO: load moves from localstorage on app boot
