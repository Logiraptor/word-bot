import "./core/main.go";
import "./index.scss";

import * as React from "react";
import * as ReactDOM from "react-dom";

import { App } from "./components/App";
import { GameService, StorageService } from "./services/game";

let app = document.createElement("div");
document.body.appendChild(app);
ReactDOM.render(<App gameService={new GameService()} storage={new StorageService()} />, app);
