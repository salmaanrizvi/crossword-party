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
const isFromServer = action => action && action.isFromServer

const connect = () => {
  const channel = getChannel()
  if (!channel) {
    console.log("[connect]: crossword party is not active!")
    return null
  }

  const ws = new WebSocket(process.env.__API_BASE_URL)
  ws.from = uuidv4()
  ws.channel = channel
  ws.clientVersion = process.env.__CWP_APP_VERSION

  ws.onopen = () => {
    console.log("[onopen]: websocket connection opened")
    ws.send(
      JSON.stringify({
        type: __CROSSWORD_PARTY_REGISTER,
        from: ws.from,
        channel: ws.channel,
        timestamp: (new Date).toISOString(),
        clientVersion: ws.clientVersion,
        gameId: ws.gameId,
      })
    )
  }
  
  ws.onmessage = msg => {
    const { data } = msg
  
    let action
    try {
      action = JSON.parse(data);
    } catch (e) {
      console.error('[onmessage]: error parsing message data', data);
      return
    }
  
    // extra check to not dispatch our own actions if we
    // somehow receive them on the socket
    if (action.from !== ws.from) {
      if (!ws.dispatch) return

      console.log("[onmessage]: received message from server, dispatching locally", action)
      ws.dispatch(action)
    }
  }

  // TODO: what should we do if we receive a close event?
  ws.onclose = event => {
    console.log('[onclose]: received close event', event)
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
    console.log("[setGameIdMiddleware]: unable to get game id from state yet")
    return next(action)
  }

  websocket.gameId = gameId

  console.log("[setGameIdMiddleware]: setting game id")
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

// const handleSyncGameRequest = store => {
//   const state = store.getState()
//   const {
//     gamePageData: {
//       cells,
//       status,
//       timer,
//       modal
//     } = {}
//   } = state

  
// }

// const handleCwpActionMiddleware = websocket => store => next => action => {
//   if (!isValidWs(websocket)) {
//     return next(action)
//   }

//   switch(action.type) {
//     case '__CWP_SYNC_GAME': {
//       handleSyncGameRequest(store)
//     }
//   }

//   return next(action)
// }

const postActionMiddleware = websocket => store => next => action => {
  if (!isValidWs(websocket)) {
    console.log("[postActionMiddleware]: invalid websocket, skipping postActionMiddleware!!!!!!")
    return next(action)
  }

  // Check if action was sent to us from websocket
  if (!isFromServer(action)) {
    console.log("[postActionMiddleware]: action was not from server, sending up...", action.type)
    action.from = websocket.from
    action.channel = websocket.channel
    action.clientVersion = websocket.clientVersion
    action.gameId = websocket.gameId
    action.timestamp = (new Date).toISOString()
    websocket.send(JSON.stringify(action));  
  } else {
    console.log("[postActionMiddleware]: action IS from server, not sending up", action.type, action.isFromServer)
  }

  return next(action);
}

// setDispatchMiddleware is solely meant to keep the referenced store
// on the websocket instance up to date for when the 
const setDispatchMiddleware = websocket => store => next => action => {
  if (websocket && !websocket.dispatch) {
    console.log("[setDispatchMiddleware]: setting dispatch on websocket")
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
    console.log("[handleActionMiddleware]: invalid websocket, skipping middleware")
    return next(action)
  }

  if (!isFromServer(action)) {
    console.log("[handleActionMiddleware]: received local action, skipping...", action.type)
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
  console.log('client_id', ws.from)
  let mwares = [
    setDispatchMiddleware(ws),
    setGameIdMiddleware(ws),
    handleActionMiddleware(ws),
    postActionMiddleware(ws),
  ]

  // if (process.env.NODE_ENV === 'development') {
  //   mwares.unshift(logger)
  // }


  // window.addEventListener('message', event => {
  //   if (!event.data || event.data.soure !== '__cwp') {
  //     return
  //   }

  //   switch (event.data.type) {
  //     case '__CWP_SYNC_GAME': {
  //       ws.dispatch()
  //     }
  //   }
  //     console.log('this is good data', event.data)
  // })

  // key to backdoor :)
  window.devToolsExtension = () => applyMiddleware(...mwares)
}
