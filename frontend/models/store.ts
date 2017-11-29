import { createStore, Store, applyMiddleware, MiddlewareAPI, Dispatch, Middleware } from "redux";
import { Tile, Move, Board } from "./core";
import { receiveValidations, Action, receiveRender } from "./actions";
import { GameService, LocalStorage } from "../services/game";

export interface AppStore {
    moves: Move[];
    rack: Tile[];
    board: Board;
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
    constructor(private gameService: GameService, private storage: LocalStorage<AppStore>) {}

    createStore() {
        return createStore<AppStore>(
            this.reducer,
            applyMiddleware(this.validator as Middleware, this.renderer as Middleware, this.persister as Middleware),
        );
    }

    persister = (store: MiddlewareAPI<AppStore>) => (next: Dispatch<AppStore>) => (action: Action) => {
        next(action);
        console.log(action);
        const state = store.getState();
        this.storage.save(state);
    };

    validator = (store: MiddlewareAPI<AppStore>) => (next: Dispatch<AppStore>) => (action: Action) => {
        next(action);
        if (action.type !== "receivevalidations" && action.type !== "receiverender") {
            const state = store.getState();
            this.gameService.validate(state).then((validations) => {
                store.dispatch(receiveValidations(validations));
            });
        }
    };

    renderer = (store: MiddlewareAPI<AppStore>) => (next: Dispatch<AppStore>) => (action: Action) => {
        next(action);
        if (action.type !== "receiverender" && action.type !== "receivevalidations") {
            const state = store.getState();
            this.gameService.render(state).then((result) => {
                store.dispatch(receiveRender(result.Board, result.Scores));
            });
        }
    };

    reducer = (state: AppStore | undefined, action: Action): AppStore => {
        if (!state) {
            return this.storage.load();
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
