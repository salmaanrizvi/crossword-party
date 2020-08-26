import React, { useContext } from 'react'
import { v4 as uuidv4 } from 'uuid'
import { Button, Typography } from '@material-ui/core'

import { ReducerContext, ACTIONS } from '../reducer'
import Chrome from '../chrome'

export const CreateChannel = () => {
  const { state, dispatch } = useContext(ReducerContext)

  if (state.channelId) {
    return null
  }

  const handleClick = () => {
    const channelId = uuidv4()
    const url = `${ state.url }?cwp_channel=${channelId}`

    Chrome.updateTab(null, { url }).then(() => {
      dispatch({ type: ACTIONS.SetChannelId, channelId, url })
    })
  }

  return (
    <div className="flex-center flex-col fh fw margin-12">
      <Typography variant="body2" paragraph>
        Start a party below and enjoy playing a multiplayer version of the NYTimes crossword.
      </Typography>
      <Button color="primary" variant="contained" onClick={handleClick} classes={{ root: 'create-party-button' }}>
        Start a new party!
      </Button>
    </div>
  )
}

CreateChannel.defaultProps = {
  isSessionActive: false
}
