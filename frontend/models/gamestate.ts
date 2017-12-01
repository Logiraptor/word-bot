import { Move, Tile } from "./core";
import { AppStore } from "./store";
import { Store } from "redux";
import { setRack, updateMove, addMove, deleteMove } from "./actions";

export class GameState {
    constructor(private store: Store<AppStore>) {}

    setRack(rack: Tile[]) {
        this.store.dispatch(setRack(rack));
    }
    updateMove(i: number, value: Move) {
        this.store.dispatch(updateMove(value, i));
    }
    addMove(value: Move) {
        this.store.dispatch(addMove(value));
    }
    removeMove(i: number) {
        this.store.dispatch(deleteMove(i));
    }

    subscribe(setState: (s: AppStore) => void) {
        this.store.subscribe(() => {
            setState(this.store.getState());
        });
    }
}
