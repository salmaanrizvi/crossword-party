import React, { useReducer } from 'react';

import '../../css/popup.css'
import { reducer, ReducerContext } from './reducer'
import { useLoadSettings } from './hooks'

import Header from './Header'
import Loader from './Loader'
import { ActiveChannel, CreateChannel } from './Channels'

const PopupApp = () => {
  const [state, dispatch] = useReducer(reducer, { channelId: '', loading: true })
  useLoadSettings(dispatch)

  return (
    <ReducerContext.Provider value={{ state, dispatch }}>
      <div className="app flex-center flex-col">
        <Header loading={state.loading} />
        { state.loading 
          ? <Loader />
          : (<React.Fragment>
              <ActiveChannel />
              <CreateChannel />
            </React.Fragment>
          )}
      </div>
    </ReducerContext.Provider>
  )
}

export default PopupApp
