import { AxiosInstance } from 'axios';
import { createAPIClient, testConnection, TranzaAPIConfig } from '../clients/tranza-api';

export interface UserSession {
  userId: string;
  apiKey: string;
  apiClient: AxiosInstance;
  authenticated: boolean;
  lastActivity: Date;
  expiresAt: Date;
}

export interface SessionConfig {
  sessionTimeout: number; // in minutes
  maxSessions: number;
}

export interface AuthenticationResult {
  success: boolean;
  message: string;
  session?: UserSession;
}

// In-memory session storage
const sessions = new Map<string, UserSession>();

// Default configuration
const defaultConfig: SessionConfig = {
  sessionTimeout: 60, // 1 hour default
  maxSessions: 1000,
};

let currentConfig: SessionConfig = { ...defaultConfig };
let cleanupInterval: NodeJS.Timeout;

/**
 * Initialize the session manager with configuration
 */
export const initializeSessionManager = (config: Partial<SessionConfig> = {}): void => {
  currentConfig = { ...defaultConfig, ...config };
  
  // Start cleanup interval if not already running
  if (!cleanupInterval) {
    cleanupInterval = setInterval(() => {
      cleanupExpiredSessions();
    }, 5 * 60 * 1000); // Cleanup every 5 minutes
  }
  
  console.log('🚀 Session manager initialized with config:', currentConfig);
};

/**
 * Authenticate a user with their API key
 */
export const authenticateUser = async (userId: string, apiKey: string): Promise<AuthenticationResult> => {
  try {
    console.log(`🔐 Authenticating user: ${userId}`);

    // Create API client for testing
    const apiConfig: TranzaAPIConfig = {
      baseURL: process.env['TRANZA_API_BASE_URL'] || 'http://localhost:8080',
      apiKey: apiKey,
      timeout: 10000,
    };

    const apiClient = createAPIClient(apiConfig);

    // Test the connection
    const isValid = await testConnection(apiClient);
    
    if (!isValid) {
      console.log(`❌ Authentication failed for user: ${userId}`);
      return {
        success: false,
        message: 'Invalid API key or unable to connect to Tranza backend',
      };
    }

    // Create or update session
    const session: UserSession = {
      userId,
      apiKey,
      apiClient,
      authenticated: true,
      lastActivity: new Date(),
      expiresAt: new Date(Date.now() + currentConfig.sessionTimeout * 60 * 1000),
    };

    sessions.set(userId, session);
    console.log(`✅ User authenticated successfully: ${userId}`);
    
    return {
      success: true,
      message: 'Successfully authenticated! You can now use transfer commands.',
      session,
    };
  } catch (error) {
    console.error(`❌ Authentication error for user ${userId}:`, error);
    return {
      success: false,
      message: error instanceof Error ? error.message : 'Authentication failed',
    };
  }
};

/**
 * Get user session if authenticated and valid
 */
export const getUserSession = (userId: string): UserSession | null => {
  const session = sessions.get(userId);
  
  if (!session) {
    return null;
  }

  // Check if session is expired
  if (new Date() > session.expiresAt) {
    sessions.delete(userId);
    console.log(`⏰ Session expired for user: ${userId}`);
    return null;
  }

  // Update last activity
  session.lastActivity = new Date();
  session.expiresAt = new Date(Date.now() + currentConfig.sessionTimeout * 60 * 1000);
  
  return session;
};

/**
 * Get authenticated API client for user
 */
export const getAPIClient = (userId: string): AxiosInstance | null => {
  const session = getUserSession(userId);
  return session?.apiClient || null;
};

/**
 * Check if user is authenticated
 */
export const isUserAuthenticated = (userId: string): boolean => {
  const session = getUserSession(userId);
  return session?.authenticated || false;
};

/**
 * Logout user (remove session)
 */
export const logoutUser = (userId: string): boolean => {
  const existed = sessions.has(userId);
  sessions.delete(userId);
  
  if (existed) {
    console.log(`👋 User logged out: ${userId}`);
  }
  
  return existed;
};

/**
 * Get session info for debugging
 */
export const getSessionInfo = (userId: string): any => {
  const session = sessions.get(userId);
  
  if (!session) {
    return { authenticated: false };
  }

  return {
    authenticated: true,
    lastActivity: session.lastActivity,
    expiresAt: session.expiresAt,
    timeRemaining: session.expiresAt.getTime() - new Date().getTime(),
  };
};

/**
 * Get session statistics
 */
export const getSessionStats = (): { activeSessions: number; totalUsers: string[] } => {
  const activeSessions = sessions.size;
  const totalUsers = Array.from(sessions.keys());
  
  return { activeSessions, totalUsers };
};

/**
 * Clean up expired sessions
 */
export const cleanupExpiredSessions = (): void => {
  const now = new Date();
  let cleanedCount = 0;

  for (const [userId, session] of sessions.entries()) {
    if (now > session.expiresAt) {
      sessions.delete(userId);
      cleanedCount++;
    }
  }

  if (cleanedCount > 0) {
    console.log(`🧹 Cleaned up ${cleanedCount} expired sessions`);
  }
};

/**
 * Update session configuration
 */
export const updateSessionConfig = (config: Partial<SessionConfig>): void => {
  currentConfig = { ...currentConfig, ...config };
  console.log('⚙️ Session config updated:', currentConfig);
};

/**
 * Get current session configuration
 */
export const getSessionConfig = (): SessionConfig => {
  return { ...currentConfig };
};

/**
 * Cleanup resources and stop intervals
 */
export const destroySessionManager = (): void => {
  if (cleanupInterval) {
    clearInterval(cleanupInterval);
  }
  sessions.clear();
  console.log('🔒 Session manager destroyed');
};

// Auto-initialize with default config
initializeSessionManager();
