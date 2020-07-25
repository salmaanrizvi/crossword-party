
// '4Vwr': 
function (e, t, r) {
  'use strict';
  function n(e) {
    return e && e.__esModule ? e : {
    default:
      e
    }
  }
  Object.defineProperty(t, '__esModule', {
    value: !0
  });
  var a = 'function' == typeof Symbol && 'symbol' == typeof Symbol.iterator ? function (e) {
    return typeof e
  }
   : function (e) {
    return e && 'function' == typeof Symbol && e.constructor === Symbol && e !== Symbol.prototype ? 'symbol' : typeof e
  },
  o = r('2KeS'), // o is redux object with createStore, compose, applyMiddleware
  s = r('bEzl'),
  i = n(s),
  l = r('zyqB'),
  u = n(l),
  c = r('1L3B'),
  d = n(c),
  _ = r('/6kf'), // _.potentialSyncActions is array of dispatched actions (looks like)
  f = n(_),
  p = r('r2rr'),
  m = n(p),
  E = r('La8e'), // E.todaysDateReducer(), e.toggleMenu(), plus other action constants like GET_LEADERBOARDS_PROFILE, GET_PUZZLE_COLLECTIONS, GET_STATS_AND_STREAKS, PUZZLE_FETCH_FAILED, UPDATE_PUZZLES_PROGRESS
  g = n(E),
  h = r('lGlm'), // lots of useful selectors, getCurrentUser, getIsMini, getAllCells etc
  O = r('aV6u'),
  T = n(O),
  A = 'object' === ('undefined' == typeof window ? 'undefined' : a(window)),
  y = A && void 0 !== window.devToolsExtension,
  R = function (e) { // this is the actual function that gets called when mounting the app
    /* Example input:
{
  "isMounted": false,
  "openMenu": null,
  "user": {
    "info": {
      "gdpr": false,
      "ads": true,
      "optOut": "show",
      "ccpa": "show",
      "isLoggedIn": true,
      "hasXwd": true,
      "hasDigi": true,
      "inShortzMode": false,
      "isErsatzShortz": false,
      "displayName": "¯\\_(ツ)_/¯"
    },
    "printPrefs": {
      "opacity": 100
    },
    "settings": {
      "onSwitch": "stay",
      "jumpBack": true,
      "skipFilled": true,
      "skipPenciled": false,
      "showTimer": true,
      "suppressDisqualificationWarnings": false,
      "spaceTriggers": "toggle",
      "backspaceAcrossWords": true,
      "autoAdvance": true,
      "soundOn": false
    }
  },
  "device": {
    "isMobile": false,
    "isTablet": false,
    "isPhone": false
  },
  "printOptions": {
    "selectedPuzzle": {},
    "isGeneratable": false,
    "version": "puzzle",
    "showBlack": false,
    "showSpoiler": true
  },
  "error": null,
  "profile": {
    "pending": true
  },
  "newGames": [],
  "adUnitPath": "crossword",
  "todaysDate": null,
  "gamePageData": {
    "assets": {},
    "board": {
      "name": "svg",
      "attributes": {
        "xmlns": "http://www.w3.org/2000/svg",
        "viewBox": "0 0 501.00 501.00"
      },
      "children": [], // svg data
      "styles": [
        {
          "name": "font-family",
          "value": "helvetica,arial,sans-serif"
        }
      ]
    },
    "clueLists": [{}, ...],
    "meta": {
      "category": 0,
      "constructors": [
        "Hal Moore"
      ],
      "copyright": "2020",
      "editor": "Will Shortz",
      "id": 18370,
      "lastUpdated": "2020-06-24 19:51:29 +0000 UTC",
      "publicationDate": "2020-07-03",
      "relatedContent": {
        "text": "Read and comment on the Wordplay blog",
        "url": "https://www.nytimes.com/2020/07/02/crosswords/daily-puzzle-2020-07-03.html"
      },
      "goldStarCutoff": "2020-07-03T23:59:59-07:00",
      "publishStream": "daily"
    },
    "dimensions": {
      "rowCount": 15,
      "columnCount": 15,
      "aspectRatio": 1,
      "cellSize": 33
    },
    "overlays": {
      "beforeStart": false,
      "afterSolve": false
    },
    "cells": [
      {
        "type": 1,
        "clues": [
          0,
          33
        ],
        "answer": "A",
        "label": "1",
        "index": 0
      },
      {
        "type": 1,
        "clues": [
          0,
          34
        ],
        "answer": "S",
        "label": "2",
        "index": 1
      },
      {
        "type": 1,
        "clues": [
          0,
          35
        ],
        "answer": "S",
        "label": "3",
        "index": 2
      },
      ...
    ],
    "clues": [
      {
        "list": 0,
        "cells": [
          0,
          1,
          2,
          3
        ],
        "direction": "Across",
        "label": "1",
        "text": "Abbr. in some job titles",
        "index": 0,
        "unfilledCount": 4,
        "prev": 69,
        "next": 1
      },
      {
        "list": 0,
        "cells": [
          5,
          6,
          7,
          8
        ],
        "direction": "Across",
        "label": "5",
        "text": "Sustain",
        "index": 1,
        "unfilledCount": 4,
        "prev": 0,
        "next": 2
      },
      {
        "list": 0,
        "cells": [
          10,
          11,
          12,
          13,
          14
        ],
        "direction": "Across",
        "label": "9",
        "text": "Singles player in the 1950s",
        "index": 2,
        "unfilledCount": 5,
        "prev": 1,
        "next": 3
      },
      ...
    ],
    "modal": {
      "type": "START_VEIL",
      "config": {},
      "isModal": false
    },
    "selection": {
      "cell": null,
      "clueCells": [],
      "clue": null,
      "cellClues": [],
      "clueList": null,
      "relatedCells": [],
      "relatedClues": []
    },
    "status": {
      "firsts": {},
      "isSolved": false,
      "isFilled": false,
      "autocheckEnabled": false,
      "blankCells": 191,
      "incorrectCells": 191,
      "lastCommitID": null,
      "currentProgress": 0
    },
    "timer": {
      "totalElapsedTime": 0,
      "sessionElapsedTime": 0,
      "resetSinceLastCommit": false
    },
    "toolbar": {
      "activeMenu": null,
      "inPencilMode": false,
      "inRebusMode": false
    },
    "transient": {
      "isObscured": false,
      "isReady": false,
      "isSynced": false,
      "doEscape": false
    }
  },
  "collections": {},
  "puzzles": {},
  "stats": {},
  "kenkenAndSetRetired": false
}
    */
    var t = (0, h.getPuzzleMetadata) (e) || {
    },
    r = [
      (0, o.applyMiddleware) (i.default,
      d.default)]; return t.id && r.push((0, o.applyMiddleware) (f.default,
      u.default,
      m.default,
      T.default)),
      y && r.push(window.devToolsExtension()),
      o.compose.apply(void 0, r) (o.createStore) (g.default,
      e)
    };
    t.default = R;
    !function () {
      'undefined' != typeof __REACT_HOT_LOADER__ && (__REACT_HOT_LOADER__.register(A, 'isBrowser', '/drone/src/src/store/configureStore.js'), __REACT_HOT_LOADER__.register(y, 'hasDevTools', '/drone/src/src/store/configureStore.js'), __REACT_HOT_LOADER__.register(R, 'default', '/drone/src/src/store/configureStore.js'))
    }()
  }