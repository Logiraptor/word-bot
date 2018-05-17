import './core/main.go'
import './index.scss'

import * as React from 'react'
import * as ReactDOM from 'react-dom'

import {App} from './components/App'
import {GameService, LocalStorage} from './services/game'
import {AppState, DefaultState} from './models/store'
import {setRack} from './models/actions'
import {GameState} from './models/gamestate'

const gameService = new GameService()
const storage = new LocalStorage('appstate', DefaultState)
const appstate = new AppState(gameService, storage)

const store = appstate.createStore()
store.dispatch(setRack([]))
let app = document.createElement('div')
document.body.appendChild(app)

ReactDOM.render(<App gameState={new GameState(store)}/>, app)
