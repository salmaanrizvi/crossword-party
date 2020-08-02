var script = document.createElement('script');
script.id = '__crossword_party';
script.src = chrome.runtime.getURL('injectMiddleware.bundle.js');
(document.head || document.documentElement).appendChild(script);
console.log('loaded script')
