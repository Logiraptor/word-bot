import "./core/main.go";
import "./index.scss";

import * as React from "react";
import * as ReactDOM from "react-dom";

import { App } from "./components/App";
import { GameService, StorageService } from "./services/game";
import { createStore, applyMiddleware } from "redux";
import { AppStore, setRack, AppState } from "./models/store";

const gameService: GameService = new GameService();
const storageService: StorageService = new StorageService();
const appstate = new AppState(gameService);

const store = appstate.createStore();
store.dispatch(setRack([]));
let app = document.createElement("div");
document.body.appendChild(app);

ReactDOM.render(<App store={store} gameService={gameService} storage={storageService} />, app);
