/**
 * @redditview/core - Main entry point
 * Exports all public APIs and types
 */

// Models
export * from './models/index.js'

// API Client
export { RedditApiClient, apiClient } from './api/index.js'
export type { ApiClientConfig } from './api/index.js'

// Cache
export { PostCache, LocalStorageCache } from './cache/index.js'
export type { CacheConfig } from './cache/index.js'

// Utils
export {
  formatTime,
  formatTimeAgo,
  formatNum,
  getSubredditFromInput,
  isTextPost,
  getPostThumbnail,
} from './utils/index.js'
