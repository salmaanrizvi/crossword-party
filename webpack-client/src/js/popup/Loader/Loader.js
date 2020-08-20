import React from 'react'
import { CircularProgress } from '@material-ui/core'

const Loader = () => {
  return (
    <div className="fw fh flex-center">
      <CircularProgress disableShrink />
    </div>
  )
}

export default Loader
