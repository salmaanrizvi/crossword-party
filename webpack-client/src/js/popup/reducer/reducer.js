import { useReducer, createContext } from 'react'
import Chrome from '../chrome'

export const ReducerContext = createContext({ state: {}, dispatch: () => {} });

export const ACTIONS = {
  LoadSettings: 'LOAD_SETTINGS',
  SyncSettings: 'SYNC_SETTINGS',

  SetChannelId: 'SET_CHANNEL_ID',
  RemoveChannelId: 'REMOVE_CHANNEL_ID',
}

export const reducer = (state, action) => {
  switch (action.type) {
    case ACTIONS.LoadSettings: {
      return Object.assign({}, state, { loading: true })
    }
    case ACTIONS.SyncSettings: {
      const newState = Object.assign({}, state, action.data, { loading: false })
      Chrome.set({ channelId: newState.channelId })
      return newState
    }

    case ACTIONS.SetChannelId: {
      Chrome.set({ channelId: action.channelId })
      return Object.assign({}, state, { channelId: action.channelId, url: action.url })
    }
    case ACTIONS.RemoveChannelId: {
      Chrome.set({ channelId: '' })
      return Object.assign({}, state, { channelId: '' })
    }
    default:
      return state
  }
}
