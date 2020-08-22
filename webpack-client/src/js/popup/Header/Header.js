import React from 'react'
import PropTypes from 'prop-types'

import { AppBar, Toolbar, Typography } from '@material-ui/core'

const Header = ({ loading }) => {
  return (
    <AppBar position={ loading ? "fixed" : "relative" }>
      <Toolbar variant="dense">
        <Typography variant="h6" className="header">
          Crossword Party
        </Typography>
      </Toolbar>
    </AppBar>
  )
}

Header.propTypes = {
  loading: PropTypes.bool
}

Header.defaultProps = {
  loading: false,
}

export default Header
