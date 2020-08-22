import '../img/icon-128.png'
import '../img/icon-34.png'
import { TableSortLabel } from '@material-ui/core';

var lastTabId = 0;

chrome.tabs.onSelectionChanged.addListener(function(tabId) {
  lastTabId = tabId;
  console.log('tab changed to', tabId)
  chrome.tabs.query({ url: 'https://www.nytimes.com/crosswords*'}, tabs => {
    const currentTab = tabs.find(tab => tab.id === tabId)
    if (currentTab) {
      chrome.pageAction.show(lastTabId);
    } else {
      chrome.pageAction.hide(lastTabId)
    }
  })
});

chrome.tabs.query({active: true, currentWindow: true}, function(tabs) {
  lastTabId = tabs[0].id;
  if (tabs[0].url.indexOf('https://www.nytimes.com/crosswords') > -1) {
    chrome.pageAction.show(lastTabId);
  }
});
