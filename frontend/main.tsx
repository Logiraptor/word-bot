import "./core/main.go";
import "./index.scss";

import * as React from "react";
import * as ReactDOM from "react-dom";

import { App } from "./components/App";
import { GameService, StorageService } from "./services/game";
import { createStore } from "redux";
import { AppStore, reducer } from "./models/store";

const store = createStore<AppStore>(reducer);
let app = document.createElement("div");
document.body.appendChild(app);
ReactDOM.render(<App store={store} gameService={new GameService()} storage={new StorageService()} />, app);
