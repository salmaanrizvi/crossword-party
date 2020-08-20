const Chrome = {}

Chrome.get = args => new Promise(resolve => {
  chrome.storage.sync.get(args, resolve)
})

Chrome.set = obj => new Promise(resolve => {
  chrome.storage.sync.set(obj, () => resolve(obj))
})

Chrome.activeTab = () => new Promise(resolve => {
  var query = { active: true, currentWindow: true };
  chrome.tabs.query(query, tabs => resolve(tabs[0]))
})

Chrome.updateTab = (tabId, updateProperties) => new Promise(resolve => {
  chrome.tabs.update(tabId, updateProperties, resolve)
})

export default Chrome
