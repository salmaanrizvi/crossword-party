const networkFilters = { urls: ["*://*.nytimes.com/*"] };

const onCompleted = details => {
  // console.log('url', details.url)
  if (details.url.indexOf('https://www.nytimes.com/games-assets/main-') != -1) {
    console.log('details obj is', details)
  }
}

chrome.webRequest.onBeforeRequest.addListener(
  function (details) {
    const { url = "", initiator = "" } = details;
    const jsRequest = url.indexOf('https://www.nytimes.com/games-assets/main-') !== -1
    const siteInitiated = initiator.indexOf("chrome-extension://") === -1
    if (jsRequest && siteInitiated) {
      // const response = fetch(details.url, { method: details.method })
      // console.log('response', response);
      // response.then(res => console.log("res", res.)

      // const request = new XMLHttpRequest();
      // request.open(details.method, details.url, true);

      // request.onload = (...args) => {
      //   console.log(request.response);
      // }

      // request.send(null)
      return { cancel: false }
      // return { redirectUrl: 'http://localhost:12312/games-assets/main.js' };
    }

    return { cancel: false };
    // var javascriptCode = loadSynchronously(details.url);

    // modify javascriptCode here
    // return { redirectUrl: "data:text/javascript," + encodeURIComponent(javascriptCode) };
  },
  networkFilters,
  ["blocking"]
);

chrome.webRequest.onCompleted.addListener(onCompleted, networkFilters);

console.log('added listener');
