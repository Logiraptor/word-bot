import "./core/main.go";
import "./index.scss";

import * as React from "react";
import * as ReactDOM from "react-dom";

import { App } from "./components/App";
import { GameService, StorageService } from "./services/game";
import { createStore } from "redux";
import { AppStore, reducer, setupSubscriptions } from "./models/store";

const store = createStore<AppStore>(reducer);
const gameService: GameService = new GameService();
setupSubscriptions(store, gameService);
let app = document.createElement("div");
document.body.appendChild(app);

ReactDOM.render(<App store={store} gameService={gameService} storage={new StorageService()} />, app);
