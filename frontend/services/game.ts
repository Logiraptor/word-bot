import { MoveRequest, RenderedBoard, Move } from "../models/core";


export class GameService  {
    async render(req: MoveRequest): Promise<RenderedBoard> {
        return (await fetch("/render", {
            method: "POST",
            body: JSON.stringify(req),
        }).then((x) => x.json())) as RenderedBoard;
    }
    
    async play(req: MoveRequest): Promise<Move> {
        return await fetch("/play", {
            method: "POST",
            body: JSON.stringify(req),
        }).then((x) => x.json());
      }
}


export class StorageService {
    async load(): Promise<Move[]> {
        let movesString = localStorage.getItem("moves");
        if (!movesString) {
            movesString = "[]";
        }

        let moves: Move[];
        try {
            moves = JSON.parse(movesString);
        } catch (e) {
            moves = [];
        }
        return moves;
    }

    async save(game: Move[]) {
        localStorage.setItem("moves", JSON.stringify(game));
    }
}