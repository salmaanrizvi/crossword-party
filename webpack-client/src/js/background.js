import '../img/icon-128.png'
import '../img/icon-34.png'

var lastTabId = 0;
const urlRegex = new RegExp('https:\/\/www\.nytimes\.com\/crosswords\/game\/.*')


chrome.tabs.onActivated.addListener(({ tabId }) => {
  lastTabId = tabId;
  chrome.tabs.query({ url: 'https://www.nytimes.com/crosswords/game/*'}, tabs => {
    const currentTab = tabs.find(tab => tab.id === tabId)
    if (currentTab) {
      chrome.pageAction.show(lastTabId);
    } else {
      chrome.pageAction.hide(lastTabId)
    }
  })
});

chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
  if (urlRegex.test(tab.url)) {
    chrome.pageAction.show(tab.id)
  } else {
    chrome.pageAction.hide(tab.id)
  }
})

chrome.tabs.query({active: true, currentWindow: true}, function(tabs) {
  const tab = tabs[0]
  lastTabId = tab.id

  if (urlRegex.test(tab.url)) {
    chrome.pageAction.show(lastTabId);
  } else {
    chrome.pageAction.hide(lastTabId)
  }
});
