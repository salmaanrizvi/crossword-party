var script = document.createElement('script');
script.id = '__crossword_party';
script.src = chrome.runtime.getURL('injectMiddleware.bundle.js');
(document.head || document.documentElement).appendChild(script);
console.log('loaded script')

// chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
//   console.log(sender.tab ?
//               "from a content script:" + sender.tab.url :
//               "from the extension")
//   if (request.type == "SYNC_GAME") {
//     window.postMessage({ source: '__cwp', ...request })
//     sendResponse({ received: true })
//   }
// })
