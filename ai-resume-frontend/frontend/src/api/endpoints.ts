// API Endpoints configuration

export const API_ENDPOINTS = {
  // Authentication
  AUTH: {
    REGISTER: '/auth/register',
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    REFRESH: '/auth/refresh',
    VERIFY_EMAIL: '/auth/verify',
    RESET_PASSWORD: '/auth/reset-password',
  },
  
  // User Profile
  PROFILE: {
    GET: '/profile',
    UPDATE: '/profile',
    DELETE: '/profile',
    
    // Experience
    EXPERIENCE: {
      LIST: '/profile/experience',
      CREATE: '/profile/experience',
      UPDATE: (id: string) => `/profile/experience/${id}`,
      DELETE: (id: string) => `/profile/experience/${id}`,
    },
    
    // Education
    EDUCATION: {
      LIST: '/profile/education',
      CREATE: '/profile/education',
      UPDATE: (id: string) => `/profile/education/${id}`,
      DELETE: (id: string) => `/profile/education/${id}`,
    },
    
    // Skills
    SKILLS: {
      LIST: '/profile/skills',
      UPDATE: '/profile/skills',
    },
    
    // Projects
    PROJECTS: {
      LIST: '/profile/projects',
      CREATE: '/profile/projects',
      UPDATE: (id: string) => `/profile/projects/${id}`,
      DELETE: (id: string) => `/profile/projects/${id}`,
    },
  },
  
  // Document Generation
  GENERATE: {
    RESUME: '/generate/resume',
    COVER_LETTER: '/generate/cover-letter',
    CHECK_CREDITS: '/generate/credits',
  },
  
  // Document Management
  DOCUMENTS: {
    LIST: '/documents',
    GET: (id: string) => `/documents/${id}`,
    UPDATE: (id: string) => `/documents/${id}`,
    DELETE: (id: string) => `/documents/${id}`,
    EXPORT_PDF: (id: string) => `/documents/${id}/export/pdf`,
    DUPLICATE: (id: string) => `/documents/${id}/duplicate`,
  },
  
  // Templates
  TEMPLATES: {
    LIST: '/templates',
    GET: (id: string) => `/templates/${id}`,
  },
  
  // Billing & Subscription
  BILLING: {
    SUBSCRIPTION: '/billing/subscription',
    UPGRADE: '/billing/upgrade',
    HISTORY: '/billing/history',
  },
} as const;

export default API_ENDPOINTS;
