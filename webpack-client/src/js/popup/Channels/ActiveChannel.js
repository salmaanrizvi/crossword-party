import React, { useContext } from 'react'
import { Avatar, Button, TextField, Typography, makeStyles } from '@material-ui/core'

import { ReducerContext, ACTIONS } from '../reducer'
import Chrome from '../chrome';
import { trimChannelId } from '../hooks'

const connected = ['S', 'I']
const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
    '& > *': {
      margin: theme.spacing(0.5),
    },
  },
  small: {
    width: theme.spacing(3),
    height: theme.spacing(3),
    fontSize: '14px'
  },
}));

export const ActiveChannel = () => {
  const { state, dispatch } = useContext(ReducerContext)
  const classes = useStyles()

  if (!state.channelId) {
    return null
  }

  const getUrl = () => {
    return `https://nytimes.com/crosswords/2020/8/10?cwp_channel=${state.channelId}`
  }

  const handleClick = () => {
    dispatch({ type: ACTIONS.RemoveChannelId })
    Chrome.activeTab().then(tab => {
      const url = trimChannelId(tab.url)
      Chrome.updateTab(tab.id, { url })
    }).then(() => {
      window.close()
    })
  }

  const url = getUrl()
  return (
    <div className="fh fw flex-center flex-col space-evenly margin-12">
      <TextField
        id="cwp-url"
        fullWidth
        label="Party url"
        variant="outlined"
        defaultValue={url}
        inputProps={{ readOnly: true }}
        size="small"
        helperText="Share this URL with your friends to start the Crossword Party!"
      />

      <div className="fw margin-top-8">
        <Typography align="left" variant="body2">
          Connected friends
        </Typography>
        <div className={classes.root}>
          { connected.map(initial => <Avatar key={initial} alt={initial} className={classes.small}>{initial}</Avatar>) }
        </div>
      </div>

      <div className="fw flex-end">
        <Button
          color="primary"
          variant="contained"
          onClick={handleClick}
          size="small"
        >
          End Party
        </Button>
      </div>
    </div>
  )
}

