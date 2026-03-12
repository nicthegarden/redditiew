import express from 'express';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const app = express();
const PORT = process.env.PORT || 5174;
const API_BASE = process.env.VITE_API_BASE_URL || 'http://localhost:8765/api';

// Serve static files from dist directory
app.use(express.static(join(__dirname, 'dist')));

// Proxy API requests to the backend
app.use('/api', (req, res) => {
  const url = new URL(req.originalUrl.replace(/^\/api/, ''), API_BASE);
  
  fetch(url.toString(), {
    method: req.method,
    headers: req.headers,
    body: req.method !== 'GET' ? req.body : undefined,
  })
    .then(apiRes => apiRes.text().then(body => ({ status: apiRes.status, headers: apiRes.headers, body })))
    .then(({ status, headers, body }) => {
      res.status(status);
      headers.forEach((value, key) => {
        if (!key.toLowerCase().startsWith('content-encoding')) {
          res.setHeader(key, value);
        }
      });
      res.send(body);
    })
    .catch(err => {
      console.error('API proxy error:', err);
      res.status(500).json({ error: 'API request failed' });
    });
});

// Serve index.html for all other routes (SPA)
app.use((req, res) => {
  res.sendFile(join(__dirname, 'dist', 'index.html'));
});

app.listen(PORT, '0.0.0.0', () => {
  console.log(`RedditView Web Server running on http://0.0.0.0:${PORT}`);
  console.log(`API Base URL: ${API_BASE}`);
});
