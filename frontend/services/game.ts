import { Move, MoveRequest, RenderedBoard } from "../models/core";

function format(m: MoveRequest, playerNames: string[] = []): string {
    let copy: any = { ...m, moves: [ ...m.moves ] };
    copy.moves = copy.moves.map((move: Move): any => {
        return { ...move, player: playerNames[move.player - 1] };
    });
    return JSON.stringify(copy);
}

export class GameService {
    async render(req: MoveRequest): Promise<RenderedBoard> {
        return (await fetch("/render", {
            method: "POST",
            body: format(req),
        }).then((x) => x.json())) as RenderedBoard;
    }

    async play(req: MoveRequest): Promise<Move> {
        return await fetch("/play", {
            method: "POST",
            body: format(req),
        }).then((x) => x.json());
    }

    async save(req: MoveRequest): Promise<void> {
        let p1 = prompt("Player 1 Name");
        let p2 = prompt("Player 2 Name");
        await fetch("/save", {
            method: "POST",
            body: format(req, [ p1, p2 ]),
        });
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
