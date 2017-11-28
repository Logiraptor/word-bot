import { createStore, Store, applyMiddleware, MiddlewareAPI, Dispatch, Middleware } from "redux";
import { Tile, Move, Board } from "./core";
import { GameService } from "../services/game";

export interface AppStore {
    moves: Move[];
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

export class AppState {
    constructor(private gameService: GameService) {}

    createStore() {
        return createStore<AppStore>(
            this.reducer,
            applyMiddleware(this.validator as Middleware, this.renderer as Middleware),
        );
    }

    validator = (store: MiddlewareAPI<AppStore>) => (next: Dispatch<AppStore>) => (action: Action) => {
        next(action);
        if (action.type !== "receivevalidations" && action.type !== "receiverender") {
            this.gameService
                .validate({ moves: store.getState().moves, rack: store.getState().rack })
                .then((validations) => {
                    store.dispatch(receiveValidations(validations));
                });
        }
    };

    renderer = (store: MiddlewareAPI<AppStore>) => (next: Dispatch<AppStore>) => (action: Action) => {
        next(action);
        if (action.type !== "receiverender" && action.type !== "receivevalidations") {
            this.gameService.render({ moves: store.getState().moves, rack: store.getState().rack }).then((result) => {
                store.dispatch(receiveRender(result.Board, result.Scores));
            });
        }
    };

    reducer = (state: AppStore | undefined, action: Action): AppStore => {
        if (!state) {
            return DefaultState;
        }
        state = { ...state };
        switch (action.type) {
            case "addmove":
                state.moves = [ ...state.moves, action.value ];
                break;
            case "deletemove":
                state.moves = [ ...state.moves ];
                state.moves.splice(action.index, 1);
                if (state.moves.length === 0) {
                    state.moves.push(EmptyMove);
                }
                break;
            case "updatemove":
                state.moves = [ ...state.moves ];
                state.moves[action.index] = action.value;
                break;
            case "setrack":
                state.rack = action.value;
                break;
            case "receiverender":
                state.moves = [ ...state.moves ];
                state.board = action.board;
                state.moves.forEach((move, i) => {
                    move.score = action.scores[i];
                });
                break;
            case "receivevalidations":
                state.moves = [ ...state.moves ];
                state.moves.forEach((move, i) => {
                    move.valid = action.validations[i];
                });
                break;
        }

        return state;
    };
}

// TODO: persist moves on change
// TODO: load moves from localstorage on app boot
