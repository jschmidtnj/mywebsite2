# dns setup

## cdn CNAME

- create CNAME with value of `cdn` pointing to `mywebsite2static.storage.googleapis.com`
- create service worker with the following script:

```javascript
const base = 'https://<bucketname>.storage.googleapis.com'

async function handleRequest(request) {
  const parsedUrl = new URL(request.url)
  let path = parsedUrl.pathname
  return fetch(base + path)
}

addEventListener('fetch', event => {
  event.respondWith(handleRequest(event.request))
})

```

follow [this](https://support.cloudflare.com/hc/en-us/articles/360013791312-Fetching-object-storage-assets-through-the-Cloudflare-CDN-using-a-Cloudflare-Worker) for more information.
