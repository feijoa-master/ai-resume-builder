// Global Types for AI Resume Builder

// User & Authentication
export interface User {
  id: string;
  email: string;
  full_name: string;
  phone?: string;
  location?: string;
  website?: string;
  linkedin?: string;
  github?: string;
  summary?: string;
  is_premium: boolean;
  credits_remaining: number;
  created_at: string;
  updated_at: string;
}

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterCredentials extends LoginCredentials {
  full_name: string;
  confirm_password: string;
}

// Profile
export interface Experience {
  id: string;
  company: string;
  position: string;
  location?: string;
  start_date: string;
  end_date?: string;
  is_current: boolean;
  description: string;
  achievements?: string[];
}

export interface Education {
  id: string;
  institution: string;
  degree: string;
  field_of_study: string;
  start_date: string;
  end_date?: string;
  is_current: boolean;
  gpa?: number;
  achievements?: string[];
}

export interface Skill {
  id: string;
  name: string;
  category: 'technical' | 'soft' | 'language' | 'other';
  proficiency?: 'beginner' | 'intermediate' | 'advanced' | 'expert';
}

export interface Project {
  id: string;
  name: string;
  description: string;
  technologies: string[];
  url?: string;
  github_url?: string;
  start_date?: string;
  end_date?: string;
}

export interface Profile {
  user: User;
  experiences: Experience[];
  education: Education[];
  skills: Skill[];
  projects: Project[];
}

// Document Generation
export interface GenerationRequest {
  job_description: string;
  job_title?: string;
  company_name?: string;
  template_id?: string;
  tone?: 'professional' | 'casual' | 'creative';
  length?: 'short' | 'medium' | 'long';
}

export interface Document {
  id: string;
  user_id: string;
  type: 'resume' | 'cover_letter';
  title: string;
  content: string;
  metadata: {
    job_title?: string;
    company_name?: string;
    template_id?: string;
    ai_model?: string;
    tokens_used?: number;
  };
  created_at: string;
  updated_at: string;
}

export interface Template {
  id: string;
  name: string;
  type: 'resume' | 'cover_letter';
  preview_url: string;
  is_premium: boolean;
  category: string;
}

// API Response types
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

// Form types
export interface ProfileFormData {
  full_name: string;
  email: string;
  phone?: string;
  location?: string;
  website?: string;
  linkedin?: string;
  github?: string;
  summary?: string;
}

export interface ExperienceFormData {
  company: string;
  position: string;
  location?: string;
  start_date: string;
  end_date?: string;
  is_current: boolean;
  description: string;
  achievements: string;
}

export interface EducationFormData {
  institution: string;
  degree: string;
  field_of_study: string;
  start_date: string;
  end_date?: string;
  is_current: boolean;
  gpa?: string;
  achievements?: string;
}

// UI State
export interface LoadingState {
  isLoading: boolean;
  error: string | null;
}

export interface ModalState {
  isOpen: boolean;
  data?: any;
}
