import http from 'http'
import https from 'https'
import { URL } from 'url'

const proxyPort = 3001
const cache = new Map<string, { time: number; data: string }>()
const CACHE_TTL = 60000
const MAX_RETRIES = 3
const RETRY_DELAY = 1000

interface RetryState {
  count: number
  nextRetry?: NodeJS.Timeout
}

const retryStates = new Map<string, RetryState>()

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

function isRateLimited(statusCode?: number): boolean {
  return statusCode === 429 || statusCode === 503
}

function getRetryState(key: string): RetryState {
  if (!retryStates.has(key)) {
    retryStates.set(key, { count: 0 })
  }
  return retryStates.get(key)!
}

const server = http.createServer((req, res) => {
  if (!req.url) {
    res.writeHead(400, { 'Content-Type': 'application/json' })
    res.end(JSON.stringify({ error: 'Bad request' }))
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
  
  const makeRequest = (attempt = 0) => {
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
        // Rate limited - retry with exponential backoff
        if (isRateLimited(proxyRes.statusCode)) {
          const retryState = getRetryState(cacheKey)
          if (retryState.count < MAX_RETRIES) {
            retryState.count++
            const delay = RETRY_DELAY * Math.pow(2, retryState.count - 1)
            console.log(`Rate limited (${proxyRes.statusCode}). Retrying in ${delay}ms (attempt ${retryState.count}/${MAX_RETRIES})`)
            retryState.nextRetry = setTimeout(() => makeRequest(retryState.count), delay)
            return
          } else {
            console.log(`Rate limited after ${MAX_RETRIES} retries`)
            res.writeHead(429, {
              'Content-Type': 'application/json',
              'Access-Control-Allow-Origin': '*',
              'Retry-After': '60'
            })
            res.end(JSON.stringify({ 
              error: 'Rate limited by Reddit',
              retry_after: 60,
              message: 'Please try again in a minute'
            }))
            retryStates.delete(cacheKey)
            return
          }
        }

        // Success - reset retry count
        retryStates.delete(cacheKey)

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
      console.error(`Request error (attempt ${attempt + 1}):`, err.message)
      retryStates.delete(cacheKey)
      res.writeHead(502, {
        'Content-Type': 'application/json',
        'Access-Control-Allow-Origin': '*'
      })
      res.end(JSON.stringify({
        error: 'Proxy error',
        message: err instanceof Error ? err.message : 'Unknown error'
      }))
    })

    proxyReq.on('timeout', () => {
      proxyReq.destroy()
      console.error(`Request timeout: ${target}`)
      res.writeHead(504, {
        'Content-Type': 'application/json',
        'Access-Control-Allow-Origin': '*'
      })
      res.end(JSON.stringify({
        error: 'Gateway timeout',
        message: 'Request took too long'
      }))
    })

    proxyReq.setTimeout(10000) // 10s timeout
    req.pipe(proxyReq)
  }

  makeRequest()
})

server.listen(proxyPort, () => {
  console.log(`✓ Proxy running on http://localhost:${proxyPort}`)
  console.log(`  Features: caching (${CACHE_TTL}ms TTL), rate-limit retry (${MAX_RETRIES} attempts)`)
})

// Graceful shutdown
process.on('SIGINT', () => {
  console.log('\nShutting down proxy...')
  server.close(() => {
    console.log('Proxy stopped')
    process.exit(0)
  })
})

