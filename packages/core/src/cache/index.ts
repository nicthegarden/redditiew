/**
 * Caching logic for Reddit data
 * Prevents unnecessary API calls and improves performance
 */

import type { RedditPost, CacheEntry } from '../models/index.js'

export interface CacheConfig {
  ttl?: number // Time to live in milliseconds
  maxSize?: number // Maximum number of cached entries
}

const DEFAULT_TTL = 3600000 // 1 hour

export class PostCache {
  private cache: Map<string, CacheEntry>
  private ttl: number
  private maxSize: number

  constructor(config: CacheConfig = {}) {
    this.cache = new Map()
    this.ttl = config.ttl ?? DEFAULT_TTL
    this.maxSize = config.maxSize ?? 100
  }

  get(key: string): RedditPost[] | null {
    const entry = this.cache.get(key)
    if (!entry) return null

    if (Date.now() - entry.time > this.ttl) {
      this.cache.delete(key)
      return null
    }

    return entry.posts
  }

  set(key: string, posts: RedditPost[]): void {
    if (this.cache.size >= this.maxSize) {
      // Remove oldest entry
      const firstKey = this.cache.keys().next().value
      if (firstKey) this.cache.delete(firstKey)
    }

    this.cache.set(key, { posts, time: Date.now() })
  }

  clear(): void {
    this.cache.clear()
  }

  delete(key: string): void {
    this.cache.delete(key)
  }

  size(): number {
    return this.cache.size
  }

  /**
   * Check if cache entry exists and is still valid
   */
  isValid(key: string): boolean {
    const entry = this.cache.get(key)
    if (!entry) return false
    return Date.now() - entry.time < this.ttl
  }

  /**
   * Get remaining TTL in milliseconds
   */
  getRemaining(key: string): number {
    const entry = this.cache.get(key)
    if (!entry) return 0
    const remaining = this.ttl - (Date.now() - entry.time)
    return Math.max(0, remaining)
  }
}

/**
 * Browser localStorage adapter for caching
 */
export class LocalStorageCache {
  private prefix: string

  constructor(prefix = 'redditview:') {
    this.prefix = prefix
  }

  private getKey(subreddit: string): string {
    return `${this.prefix}posts:${subreddit}`
  }

  get(subreddit: string): RedditPost[] | null {
    try {
      const item = localStorage.getItem(this.getKey(subreddit))
      return item ? JSON.parse(item) : null
    } catch {
      return null
    }
  }

  set(subreddit: string, posts: RedditPost[]): void {
    try {
      localStorage.setItem(this.getKey(subreddit), JSON.stringify(posts))
    } catch {
      // Silently fail if localStorage is full
    }
  }

  clear(): void {
    try {
      const keys = Object.keys(localStorage)
      keys.forEach(key => {
        if (key.startsWith(this.prefix)) {
          localStorage.removeItem(key)
        }
      })
    } catch {
      // Silently fail
    }
  }

  delete(subreddit: string): void {
    try {
      localStorage.removeItem(this.getKey(subreddit))
    } catch {
      // Silently fail
    }
  }
}
