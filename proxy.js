import http from 'http';
import https from 'https';
import { URL } from 'url';

const proxyPort = 3001;
const cache = new Map();
const CACHE_TTL = 60000;

function getCache(key) {
  const entry = cache.get(key);
  if (entry && Date.now() - entry.time < CACHE_TTL) {
    return entry.data;
  }
  return null;
}

function setCache(key, data) {
  cache.set(key, { time: Date.now(), data });
}

const server = http.createServer((req, res) => {
  const cacheKey = req.url;
  const cached = getCache(cacheKey);
  
  if (cached && !req.url.includes('/search.json')) {
    res.writeHead(200, {
      'Access-Control-Allow-Origin': '*',
      'X-Cache': 'HIT'
    });
    res.end(cached);
    return;
  }

  let target, hostname;
  let path = req.url;
  
  if (req.url.startsWith('/api/')) {
    path = req.url.replace('/api', '');
    hostname = 'www.reddit.com';
    target = 'https://www.reddit.com' + path;
  } else if (req.url.startsWith('/search/')) {
    path = '/search' + req.url.replace('/search/', '');
    hostname = 'www.reddit.com';
    target = 'https://www.reddit.com' + path;
  } else {
    hostname = 'old.reddit.com';
    target = 'https://old.reddit.com' + req.url;
  }
  
  const options = {
    hostname: hostname,
    path: path,
    method: 'GET',
    headers: {
      'User-Agent': 'redditview/1.0'
    }
  };
  
  const proxyReq = https.request(options, (proxyRes) => {
    let data = '';
    proxyRes.on('data', chunk => data += chunk);
    proxyRes.on('end', () => {
      if (proxyRes.statusCode === 200) {
        setCache(cacheKey, data);
      }
      const headers = { ...proxyRes.headers };
      delete headers['x-frame-options'];
      delete headers['content-security-policy'];
      res.writeHead(proxyRes.statusCode, {
        ...headers,
        'Access-Control-Allow-Origin': '*',
        'X-Cache': 'MISS'
      });
      res.end(data);
    });
  });
  
  proxyReq.on('error', (err) => {
    res.writeHead(500);
    res.end('Error: ' + err.message);
  });
  
  req.pipe(proxyReq);
});

server.listen(proxyPort, () => {
  console.log(`Proxy running on http://localhost:${proxyPort} (with caching + search)`);
});
