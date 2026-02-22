/**
 * API Server for RedditView
 * Exposes @redditview/core functionality via HTTP
 * Used by both React (through proxy) and Go TUI (direct calls)
 * 
 * Port: 3002
 */

import http from 'http'
import url from 'url'
import https from 'https'

const cache = new Map()
const CACHE_TTL = 60000 // 1 minute
const API_PORT = 3002

function getCache(key) {
  const entry = cache.get(key)
  if (!entry) return null
  if (Date.now() - entry.time > CACHE_TTL) {
    cache.delete(key)
    return null
  }
  return entry.data
}

function setCache(key, data) {
  cache.set(key, { data, time: Date.now() })
}

function setHeaders(res, headers) {
  res.setHeader('Content-Type', 'application/json')
  res.setHeader('Access-Control-Allow-Origin', '*')
  res.setHeader('Access-Control-Allow-Methods', 'GET, OPTIONS')
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type')
  res.setHeader('Access-Control-Max-Age', '3600')
  
  if (headers) {
    Object.entries(headers).forEach(([key, value]) => {
      res.setHeader(key, value)
    })
  }
}

function sendJson(res, data, statusCode = 200) {
  setHeaders(res)
  res.writeHead(statusCode)
  res.end(JSON.stringify(data))
}

function sendError(res, message, statusCode = 400) {
  sendJson(res, { error: message, message }, statusCode)
}

function fetchFromReddit(target) {
  return new Promise((resolve, reject) => {
    const requestUrl = new URL(target)
    
    const options = {
      hostname: requestUrl.hostname,
      path: requestUrl.pathname + requestUrl.search,
      method: 'GET',
      headers: {
        'User-Agent': 'redditview/1.0'
      }
    }

    const req = https.request(options, (res) => {
      let data = ''
      res.on('data', chunk => data += chunk)
      res.on('end', () => {
        resolve({ status: res.statusCode || 500, data })
      })
    })

    req.on('error', reject)
    req.setTimeout(10000)
    req.end()
  })
}

const server = http.createServer(async (req, res) => {
  // Handle OPTIONS for CORS preflight
  if (req.method === 'OPTIONS') {
    setHeaders(res)
    res.writeHead(200)
    res.end()
    return
  }

  if (!req.url) {
    sendError(res, 'No URL provided')
    return
  }

  const parsedUrl = url.parse(req.url, true)
  const pathname = parsedUrl.pathname || ''

  console.log(`${req.method} ${pathname}`)

  try {
    // Handle comments FIRST (must come before generic subreddit match)
    // Pattern: /api/r/:subreddit/comments/:id
    const commentsMatch = pathname.match(/^\/api(\/r\/[^/]+\/comments\/[^/]+\/)/)
    if (commentsMatch) {
      const path = commentsMatch[1]
      const redditUrl = `https://www.reddit.com${path}.json`
      const cacheKey = redditUrl

      // Check cache
      const cached = getCache(cacheKey)
      if (cached) {
        console.log(`  [CACHE HIT]`)
        setHeaders(res, { 'X-Cache': 'HIT' })
        res.writeHead(200)
        res.end(cached)
        return
      }

      // Fetch from Reddit
      const result = await fetchFromReddit(redditUrl)
      
      if (result.status === 200) {
        setCache(cacheKey, result.data)
      }

      setHeaders(res, { 'X-Cache': 'MISS' })
      res.writeHead(result.status)
      res.end(result.data)
      return
    }

    // Parse subreddit from /api/r/:subreddit or /api/r/:subreddit.json
    const subredditMatch = pathname.match(/^\/api\/r\/([^/.]+)(\.json)?(?:\/|$)/)
    if (subredditMatch) {
      const subreddit = subredditMatch[1]
      const limit = parsedUrl.query.limit || '50'
      const after = parsedUrl.query.after ? `&after=${parsedUrl.query.after}` : ''
      
      const redditUrl = `https://www.reddit.com/r/${subreddit}.json?limit=${limit}${after}`
      const cacheKey = redditUrl
      
      // Check cache
      const cached = getCache(cacheKey)
      if (cached) {
        console.log(`  [CACHE HIT]`)
        setHeaders(res, { 'X-Cache': 'HIT' })
        res.writeHead(200)
        res.end(cached)
        return
      }

      // Fetch from Reddit
      const result = await fetchFromReddit(redditUrl)
      
      if (result.status === 200) {
        setCache(cacheKey, result.data)
      }

      setHeaders(res, { 'X-Cache': 'MISS' })
      res.writeHead(result.status)
      res.end(result.data)
      return
    }

    // Handle search: /api/search.json?q=...
    if (pathname === '/api/search.json') {
      const query = parsedUrl.query.q
      if (!query || typeof query !== 'string') {
        sendError(res, 'Missing search query parameter')
        return
      }

      const limit = parsedUrl.query.limit || '50'
      const type = parsedUrl.query.type || 'link'
      
      const redditUrl = `https://www.reddit.com/search.json?q=${encodeURIComponent(query)}&type=${type}&limit=${limit}`
      const cacheKey = redditUrl

      // Search results are not cached
      const result = await fetchFromReddit(redditUrl)
      
      setHeaders(res, { 'X-Cache': 'NONE' })
      res.writeHead(result.status)
      res.end(result.data)
      return
    }

    // Health check
    if (pathname === '/health') {
      sendJson(res, { status: 'ok', cache_size: cache.size })
      return
    }

    // Stats endpoint
    if (pathname === '/api/stats') {
      sendJson(res, {
        cache_size: cache.size,
        cache_ttl: CACHE_TTL,
        uptime: process.uptime()
      })
      return
    }

    sendError(res, 'Not found', 404)
  } catch (err) {
    console.error('Error:', err)
    sendError(res, err instanceof Error ? err.message : 'Unknown error', 500)
  }
})

server.listen(API_PORT, () => {
  console.log(`
╭─────────────────────────────────────╮
│  RedditView API Server              │
│  http://localhost:${API_PORT}           │
├─────────────────────────────────────┤
│  Endpoints:                         │
│  GET /api/r/:subreddit              │
│  GET /api/r/:subreddit/comments/:id │
│  GET /api/search.json?q=:query      │
│  GET /health                        │
│  GET /api/stats                     │
╰─────────────────────────────────────╯
  `)
})

// Graceful shutdown
process.on('SIGINT', () => {
  console.log('\nShutting down API server...')
  server.close(() => {
    console.log('API server stopped')
    process.exit(0)
  })
})
