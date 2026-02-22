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
  // CORS headers for all responses
  const corsHeaders = {
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Methods': 'GET, OPTIONS',
    'Access-Control-Allow-Headers': 'Content-Type',
    'Access-Control-Max-Age': '3600'
  }

  // Handle preflight OPTIONS requests
  if (req.method === 'OPTIONS') {
    res.writeHead(200, corsHeaders)
    res.end()
    return
  }

  if (!req.url) {
    res.writeHead(400, { 'Content-Type': 'application/json', ...corsHeaders })
    res.end(JSON.stringify({ error: 'Bad request' }))
    return
  }

  const cacheKey = req.url
  const cached = getCache(cacheKey)
  
  if (cached && !req.url.includes('/search.json')) {
    res.writeHead(200, {
      'Content-Type': 'application/json',
      'X-Cache': 'HIT',
      ...corsHeaders
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
      
      // Handle redirects
      if (proxyRes.statusCode === 301 || proxyRes.statusCode === 302) {
        const location = proxyRes.headers.location
        if (location) {
          console.log(`Redirect ${proxyRes.statusCode}: ${location}`)
          const url = new URL(location)
          const redirectOptions = {
            hostname: url.hostname,
            path: url.pathname + url.search,
            method: 'GET' as const,
            headers: {
              'User-Agent': 'redditview/1.0'
            }
          }
          
          const redirectReq = https.request(redirectOptions, (redirectRes) => {
            let redirectData = ''
            redirectRes.on('data', chunk => redirectData += chunk)
            redirectRes.on('end', () => {
              res.writeHead(200, {
                'Content-Type': redirectRes.headers['content-type'] || 'application/json',
                ...corsHeaders
              })
              res.end(redirectData)
            })
          })
          redirectReq.on('error', () => {
            res.writeHead(500, corsHeaders)
            res.end('Redirect error')
          })
          redirectReq.setTimeout(10000)
          return
        }
      }

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
              'Retry-After': '60',
              ...corsHeaders
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
          'Content-Type': 'application/json',
          'X-Cache': 'MISS',
          ...corsHeaders
        })
        res.end(data)
      })
    })
    
    proxyReq.on('error', (err) => {
      console.error(`Request error (attempt ${attempt + 1}):`, err.message)
      retryStates.delete(cacheKey)
      res.writeHead(502, {
        'Content-Type': 'application/json',
        ...corsHeaders
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
        ...corsHeaders
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

