import React, { useContext } from 'react'
import { v4 as uuidv4 } from 'uuid'
import { Button, Typography } from '@material-ui/core'

import { ReducerContext, ACTIONS } from '../reducer'

export const CreateChannel = () => {
  const { state, dispatch } = useContext(ReducerContext)

  if (state.channelId) {
    return null
  }

  const handleClick = () => {
    dispatch({ type: ACTIONS.SetChannelId, channelId: uuidv4() })
  }

  return (
    <div className="flex-center flex-col fh fw margin-12">
      <Typography variant="subtitle2" paragraph>
        Start a party below and enjoy playing a multiplayer version of the NYTimes crossword.
      </Typography>
      <Button color="primary" variant="contained" onClick={handleClick}>
        Start a new party!
      </Button>
    </div>
  )
}

CreateChannel.defaultProps = {
  isSessionActive: false
}
