import { applyMiddleware } from 'redux'
import { v4 as uuidv4 } from 'uuid'

const __CROSSWORD_PARTY_CHANNEL_PARAMETER = 'cwp_channel'
const __CROSSWORD_PARTY_REGISTER = '__CROSSWORD_PARTY_REGISTER'

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

const connect = () => {
  const channel = getChannel()
  if (!channel) {
    console.log("crossword party is not active!")
    return null
  }

  const ws = new WebSocket('wss://localhost:8000/ws');
  ws.from = uuidv4()
  ws.channel = channel //'58c4c90b-041d-4232-9ae9-e219679b1130' //uuidv4()
  ws.version = process.env.__CWP_APP_VERSION   
  ws.onopen = () => ws.send(
    JSON.stringify({
      type: __CROSSWORD_PARTY_REGISTER,
      from: ws.from,
      channel: ws.channel,
      timestamp: (new Date).toISOString(),
      clientVersion: ws.version,
    })
  )
  
  ws.onmessage = msg => {
    const { target: websocket, data } = msg
  
    let action
    try {
      action = JSON.parse(data);
    } catch (e) {
      console.error('error parsing message data', data);
      return
    }
  
    if (action.from !== websocket.from) {
      if (!websocket.store && websocket.store.dispatch) return

      ws.store.dispatch(action)
    }
  }

  ws.onclose = (...args) => {
    console.log('received close event', ...args)
  }

  return ws
}

const postActionMiddleware = websocket => store => next => action => {
  // Check if action was sent to us from websocket
  if (!action.channel && websocket.readyState == 1) {
    action.from = websocket.from
    action.channel = websocket.channel
    action.timestamp = (new Date).toISOString()
    action.dispatched = true
    websocket.send(JSON.stringify(action));  
  }
  return next(action);
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

const onmessageMiddleware = websocket => store => next => action => {
  if (websocket) {
    websocket.store = store
  }

  return next(action);
}

const disruptActionMiddleware = websocket => store => next => action => {
  // TODO: think about this

  // switch (action.type) {
  //   case "APPLY_PROGRESS":
  //     // if this is from nytimes app
  //     if (action.dispatched) {
  //       action.type = '_APPLY_PROGRESS'
  //     }
  // }

  next(action)
}

const handleActionMiddleware = websocket => store => next => action => {
  if (!websocket) {
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

const ws = connect()
const mwares = [
  logger,

  disruptActionMiddleware(ws),
  handleActionMiddleware(ws),

  onmessageMiddleware(ws),
  postActionMiddleware(ws),
]

if (ws) {
  // key to backdoor :)
  window.devToolsExtension = () => applyMiddleware(...mwares)
}
