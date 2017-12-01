import { Move, MoveRequest, RenderedBoard, Tile } from "../models/core";
import { DefaultState } from "../models/store";

declare const core: {
    RenderBoard(m: string): string;
    RemainingTiles(m: string): string;
};

export class GameService {
    async render(req: MoveRequest): Promise<RenderedBoard> {
        return JSON.parse(core.RenderBoard(JSON.stringify(req)));
    }

    remainingTiles(req: MoveRequest): Tile[] {
        return JSON.parse(core.RemainingTiles(JSON.stringify(req)));
    }

    async play(req: MoveRequest): Promise<Move> {
        return await fetch("/play", {
            method: "POST",
            body: JSON.stringify(req),
        }).then((x) => x.json());
    }

    async validate(req: MoveRequest): Promise<boolean[]> {
        return await fetch("/validate", {
            method: "POST",
            body: JSON.stringify(req),
        }).then((x) => x.json());
    }
}

export class LocalStorage<T> {
    constructor(private key: string, private defaultValue: T) {}

    load(): T {
        let movesString = localStorage.getItem(this.key);
        if (!movesString) {
            return this.defaultValue;
        }

        try {
            return Object.assign({}, DefaultState, JSON.parse(movesString));
        } catch (e) {
            return this.defaultValue;
        }
    }

    save(game: T) {
        console.log("saving", game);
        localStorage.setItem(this.key, JSON.stringify(game));
    }
}
