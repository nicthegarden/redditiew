// Configuration loader for web app
interface AppConfig {
  tui: {
    default_subreddit: string
    posts_per_page: number
    list_height: number
    max_title_length: number
  }
  web: {
    default_subreddit: string
    posts_per_page: number
    theme: string
  }
  api: {
    base_url: string
    timeout_seconds: number
  }
}

let config: AppConfig | null = null

export async function loadConfig(): Promise<AppConfig> {
  if (config) return config

  try {
    const response = await fetch('/config.json')
    if (!response.ok) {
      throw new Error(`Failed to load config: ${response.statusText}`)
    }
    config = await response.json()
  } catch (error) {
    // Fallback to defaults if config.json not found
    console.warn('Could not load config.json, using defaults:', error)
    config = {
      tui: {
        default_subreddit: 'sysadmin',
        posts_per_page: 50,
        list_height: 10,
        max_title_length: 80,
      },
      web: {
        default_subreddit: 'sysadmin',
        posts_per_page: 20,
        theme: 'dark',
      },
      api: {
        base_url: '/api',
        timeout_seconds: 10,
      },
    }
  }

  return config
}

export function getConfig(): AppConfig {
  if (!config) {
    throw new Error('Config not loaded. Call loadConfig() first.')
  }
  return config
}
