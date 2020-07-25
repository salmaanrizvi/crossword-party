console.log('chrome.runtime', chrome.runtime);

// add redux into DOM i guess
// var reduxScript = document.createElement('script');
// reduxScript.id = "__crossword_party_redux";
// reduxScript.src = chrome.runtime.getURL('redux.js');
// (document.head || document.documentElement).appendChild(reduxScript)

// add our script into dom, i guess
var script = document.createElement('script');
script.id = 'salmaan';
script.src = chrome.runtime.getURL('postAction.js');
(document.head || document.documentElement).appendChild(script);

console.log('appended scripts');