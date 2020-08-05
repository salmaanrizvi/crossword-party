import { applyMiddleware } from 'redux'
import { v4 as uuidv4 } from 'uuid'

const __CROSSWORD_PARTY_CHANNEL_PARAMETER = 'cwp_channel'
const __CROSSWORD_PARTY_REGISTER = '__CROSSWORD_PARTY_REGISTER'
const __CROSSWORD_PARTY_SET_GAME_ID = '__CROSSWORD_PARTY_SET_GAME_ID'

const getChannel = () => {
  const search = window.location.search.slice(1)
  const params = search.split('&')

  for (const param of params) {
    const [key, value] = param.split('=')
    if (key === __CROSSWORD_PARTY_CHANNEL_PARAMETER) {
      return value
    }
  }

  return null
}

const isValidWs = websocket => websocket && websocket.readyState == 1

const connect = () => {
  const channel = getChannel()
  if (!channel) {
    console.log("crossword party is not active!")
    return null
  }

  const ws = new WebSocket(process.env.__API_BASE_URL)
  ws.from = uuidv4()
  ws.channel = channel
  ws.clientVersion = process.env.__CWP_APP_VERSION

  ws.onopen = () => ws.send(
    JSON.stringify({
      type: __CROSSWORD_PARTY_REGISTER,
      from: ws.from,
      channel: ws.channel,
      timestamp: (new Date).toISOString(),
      clientVersion: ws.clientVersion,
      gameId: ws.gameId,
    })
  )
  
  ws.onmessage = msg => {
    const { data } = msg
  
    let action
    try {
      action = JSON.parse(data);
    } catch (e) {
      console.error('error parsing message data', data);
      return
    }
  
    // extra check to not dispatch our own actions if we
    // somehow receive them on the socket
    if (action.from !== ws.from) {
      if (!ws.dispatch) return

      ws.dispatch(action)
    }
  }

  // TODO: what should we do if we receive a close event?
  ws.onclose = (...args) => {
    console.log('received close event', ...args)
  }

  return ws
}

const setGameIdMiddleware = websocket => store => next => action => {
  if (!isValidWs(websocket) || websocket.gameId) {
    return next(action)
  }

  const {
    gamePageData: {
      meta: {
        id: gameId
      } = {}
    } = {}
  } = store.getState()

  if (!gameId) {
    return next(action)
  }

  websocket.gameId = gameId

  websocket.send(
    JSON.stringify({
      type: __CROSSWORD_PARTY_SET_GAME_ID,
      from: websocket.from,
      channel: websocket.channel,
      timestamp: (new Date).toISOString(),
      clientVersion: websocket.clientVersion,
      gameId,
    })
  )

  return next(action)
}

const postActionMiddleware = websocket => store => next => action => {
  if (!isValidWs(websocket)) {
    return next(action)
  }

  // Check if action was sent to us from websocket
  if (!action.channel && websocket.readyState == 1) {
    action.from = websocket.from
    action.channel = websocket.channel
    action.clientVersion = websocket.clientVersion
    action.gameId = websocket.gameId
    action.timestamp = (new Date).toISOString()
    websocket.send(JSON.stringify(action));  
  }

  return next(action);
}

// setDispatchMiddleware is solely meant to keep the referenced store
// on the websocket instance up to date for when the 
const setDispatchMiddleware = websocket => store => next => action => {
  if (isValidWs(websocket) && !websocket.dispatch) {
    websocket.dispatch = store.dispatch
  }

  return next(action);
}

// handleActionMiddleware acts only on messages received on the socket
// by validating the action has the "channel" key
//
// GUESS actions are handled as a special case to not update the
// current players selected cell when another user makes a guess 
const handleActionMiddleware = websocket => store => next => action => {
  if (!isValidWs(websocket)) {
    return next(action)
  }

  if (!action.channel) {
    return next(action)
  }

  let result
  switch (action.type) {
    case 'GUESS': {
      const currentSelection = store.getState().gamePageData.selection
      result = next(action)
      store.dispatch(getSelectCellPayload(currentSelection))
      break
    }
    default:
      result = next(action)
  }

  return result
}

const getSelectCellPayload = selection => {
  return {
    type: "SELECT_CELL",
    payload: {
      index: selection.cell,
      isMiddleClick: false
    },
    selection,
  }
}

const logger = store => next => action => {
  console.group(action.type);
  console.log('prev state', store.getState());
  console.log('action', action);
  let result = next(action);
  console.log('next state', store.getState());
  console.groupEnd();
  return result;
}

const ws = connect()

if (ws) {
  let mwares = [
    setGameIdMiddleware(ws),
    handleActionMiddleware(ws),
    setDispatchMiddleware(ws),
    postActionMiddleware(ws),
  ]

  if (process.env.NODE_ENV === 'development') {
    mwares.unshift(logger)
  }

  // key to backdoor :)
  window.devToolsExtension = () => applyMiddleware(...mwares)
}
