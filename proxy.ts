import http from 'http'
import https from 'https'
import { URL } from 'url'

const proxyPort = 3001
const cache = new Map<string, { time: number; data: string }>()
const CACHE_TTL = 60000

function getCache(key: string): string | null {
  const entry = cache.get(key)
  if (entry && Date.now() - entry.time < CACHE_TTL) {
    return entry.data
  }
  return null
}

function setCache(key: string, data: string): void {
  cache.set(key, { time: Date.now(), data })
}

const server = http.createServer((req, res) => {
  if (!req.url) {
    res.writeHead(400)
    res.end('Bad request')
    return
  }

  const cacheKey = req.url
  const cached = getCache(cacheKey)
  
  if (cached && !req.url.includes('/search.json')) {
    res.writeHead(200, {
      'Access-Control-Allow-Origin': '*',
      'X-Cache': 'HIT'
    })
    res.end(cached)
    return
  }

  let target: string
  let hostname: string
  let path = req.url
  
  if (req.url.startsWith('/api/')) {
    path = req.url.replace('/api', '')
    hostname = 'www.reddit.com'
    target = 'https://www.reddit.com' + path
  } else if (req.url.startsWith('/search/')) {
    path = '/search' + req.url.replace('/search/', '')
    hostname = 'www.reddit.com'
    target = 'https://www.reddit.com' + path
  } else {
    hostname = 'old.reddit.com'
    target = 'https://old.reddit.com' + req.url
  }
  
  const options = {
    hostname: hostname,
    path: path,
    method: 'GET' as const,
    headers: {
      'User-Agent': 'redditview/1.0'
    }
  }
  
  const proxyReq = https.request(options, (proxyRes) => {
    let data = ''
    proxyRes.on('data', chunk => data += chunk)
    proxyRes.on('end', () => {
      if (proxyRes.statusCode === 200) {
        setCache(cacheKey, data)
      }
      const headers = { ...proxyRes.headers }
      delete headers['x-frame-options']
      delete headers['content-security-policy']
      res.writeHead(proxyRes.statusCode || 500, {
        ...headers,
        'Access-Control-Allow-Origin': '*',
        'X-Cache': 'MISS'
      })
      res.end(data)
    })
  })
  
  proxyReq.on('error', (err) => {
    res.writeHead(500)
    res.end('Error: ' + (err instanceof Error ? err.message : 'Unknown error'))
  })
  
  req.pipe(proxyReq)
})

server.listen(proxyPort, () => {
  console.log(`Proxy running on http://localhost:${proxyPort} (with caching + search)`)
})
