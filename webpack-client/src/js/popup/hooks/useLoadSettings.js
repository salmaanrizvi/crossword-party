import { useEffect } from 'react'

import Chrome from '../chrome'
import { ACTIONS } from '../reducer'

const __CROSSWORD_PARTY_CHANNEL_PARAMETER = 'cwp_channel'

const getChannelId = url => {
  if (!url) {
    return null
  }

  const [_, search] = url.split('?')
  if (!search) return null

  const params = search.split('&')

  for (const param of params) {
    const [key, value] = param.split('=')
    if (key === __CROSSWORD_PARTY_CHANNEL_PARAMETER) {
      return value
    }
  }

  return null
}

export const trimChannelId = url => {
  if (!url || url.indexOf('?') === -1) {
    return url
  }

  const [path, search] = url.split('?')
  const qParams = search.split('&')
  const filtered = qParams.filter(qParam => {
    return qParam.indexOf(__CROSSWORD_PARTY_CHANNEL_PARAMETER) === -1
  })

  if (!filtered.length) {
    return path
  }

  return `${path}?${filtered.join('&')}`
}

export const useLoadSettings = dispatch => {
  useEffect(() => {
    dispatch(ACTIONS.LoadSettings)
    const getData = async () => {
      let tab = {}
      let data = {}

      try {
        [tab, data] = await Promise.all([Chrome.activeTab(), Chrome.get(['channelId'])])
      } catch(e) {
        // fall through
      }
      
      // Prefer channel id from url
      const channelId = getChannelId(tab.url)
      if (channelId) {
        data.channelId = channelId
      }

      setTimeout(() => {
        dispatch({ type: ACTIONS.SyncSettings, data })
      }, 350)
    }

    getData()
  }, [])
}