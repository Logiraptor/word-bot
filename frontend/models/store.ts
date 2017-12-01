import { createStore, applyMiddleware, MiddlewareAPI, Dispatch, Middleware, Reducer } from "redux";
import { Tile, Move, Board, TileFlag } from "./core";
import { receiveValidations, Action, receiveRender, receivePlay, receiveRemainingTiles } from "./actions";
import { GameService, LocalStorage } from "../services/game";

export interface AppStore {
    moves: Move[];
    rack: Tile[];
    board: Board;
    play: Move;
    remainingTiles: Tile[];
}

export const EmptyMove: Move = {
    tiles: [],
    row: 0,
    col: 0,
    player: 1,
    direction: "horizontal",
};

export const DefaultState: AppStore = {
    moves: [ EmptyMove ],
    rack: [],
    board: Array(15).map(() => Array(15).map(() => null)),
    play: EmptyMove,
    remainingTiles: [],
};

export class AppState {
    constructor(private gameService: GameService, private storage: LocalStorage<AppStore>) {}

    createStore() {
        return createStore<AppStore>(
            this.reducer as Reducer<AppStore>,
            applyMiddleware(
                this.validator as Middleware,
                this.renderer as Middleware,
                this.persister as Middleware,
                this.player as Middleware,
                this.remainingTiler as Middleware,
            ),
        );
    }

    persister = (store: MiddlewareAPI<AppStore>) => (next: Dispatch<AppStore>) => (action: Action) => {
        next(action);
        const state = store.getState();
        this.storage.save(state);
    };

    player = (store: MiddlewareAPI<AppStore>) => (next: Dispatch<AppStore>) => (action: Action) => {
        next(action);
        if (action.changesBoard && action.type != "receiveplay") {
            this.gameService.play(store.getState()).then((play) => {
                store.dispatch(receivePlay(play));
            });
        }
    };

    validator = (store: MiddlewareAPI<AppStore>) => (next: Dispatch<AppStore>) => (action: Action) => {
        next(action);
        if (action.changesBoard) {
            const state = store.getState();
            this.gameService.validate(state).then((validations) => {
                store.dispatch(receiveValidations(validations));
            });
        }
    };

    remainingTiler = (store: MiddlewareAPI<AppStore>) => (next: Dispatch<AppStore>) => (action: Action) => {
        next(action);
        if (action.changesBoard) {
            const state = store.getState();
            const tiles = this.gameService.remainingTiles(state);
            store.dispatch(receiveRemainingTiles(tiles));
        }
    };

    renderer = (store: MiddlewareAPI<AppStore>) => (next: Dispatch<AppStore>) => (action: Action) => {
        next(action);
        if (action.changesBoard) {
            const state = store.getState();
            const movesWithPlay = [
                ...state.moves,
                {
                    ...state.play,
                    tiles: state.play.tiles.map((t) => ({
                        ...t,
                        Flags: [ TileFlag.NextAIMove ],
                    })),
                },
            ];
            this.gameService
                .render({
                    moves: movesWithPlay,
                    rack: state.rack,
                })
                .then((result) => {
                    console.log("Moves", movesWithPlay);
                    console.log("result", result);
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
                return state;
            case "deletemove":
                state.moves = [ ...state.moves ];
                state.moves.splice(action.index, 1);
                if (state.moves.length === 0) {
                    state.moves.push(EmptyMove);
                }
                return state;
            case "updatemove":
                state.moves = [ ...state.moves ];
                state.moves[action.index] = action.value;
                return state;
            case "setrack":
                state.rack = action.value;
                return state;
            case "receiverender":
                state.moves = [ ...state.moves ];
                state.board = action.board;
                state.moves.forEach((move, i) => {
                    move.score = action.scores[i];
                });
                return state;
            case "receivevalidations":
                state.moves = [ ...state.moves ];
                state.moves.forEach((move, i) => {
                    move.valid = action.validations[i];
                });
                return state;
            case "receiveplay":
                state.play = action.play;
                return state;
            case "receiveremainingtiles":
                state.remainingTiles = action.tiles;
                return state;
        }
    };
}
