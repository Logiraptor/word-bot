import "./core/main.go";
import "./index.scss";

import * as React from "react";
import * as ReactDOM from "react-dom";

import { App } from "./components/App";
import { GameService, StorageService } from "./services/game";
import { createStore } from "redux";
import { AppStore, reducer, setupSubscriptions, setRack } from "./models/store";

const store = createStore<AppStore>(reducer);
const gameService: GameService = new GameService();
const storageService: StorageService = new StorageService();
setupSubscriptions(store, gameService);
store.dispatch(setRack([]));
let app = document.createElement("div");
document.body.appendChild(app);

ReactDOM.render(<App store={store} gameService={gameService} storage={storageService} />, app);
