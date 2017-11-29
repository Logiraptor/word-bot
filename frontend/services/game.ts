import { Move, MoveRequest, RenderedBoard } from "../models/core";

declare const core: {
    RenderBoard(m: string): string;
};

export class GameService {
    async render(req: MoveRequest): Promise<RenderedBoard> {
        return JSON.parse(core.RenderBoard(JSON.stringify(req)));
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
            return JSON.parse(movesString);
        } catch (e) {
            return this.defaultValue;
        }
    }

    save(game: T) {
        console.log("saving", game);
        localStorage.setItem(this.key, JSON.stringify(game));
    }
}
