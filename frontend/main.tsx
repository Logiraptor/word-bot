import "./core/main.go";
import "./index.scss";

import * as React from "react";
import * as ReactDOM from "react-dom";

import { App } from "./components/App";
import { GameService, LocalStorage } from "./services/game";
import { createStore, applyMiddleware } from "redux";
import { AppStore, AppState, DefaultState } from "./models/store";
import { Move } from "./models/core";
import { setRack } from "./models/actions";

const gameService = new GameService();
const storage = new LocalStorage("appstate", DefaultState);
const appstate = new AppState(gameService, storage);

const store = appstate.createStore();
store.dispatch(setRack([]));
let app = document.createElement("div");
document.body.appendChild(app);

ReactDOM.render(<App store={store} gameService={gameService} />, app);
