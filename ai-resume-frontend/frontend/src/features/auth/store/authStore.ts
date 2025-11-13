import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import apiClient from '../../../api/client';
import API_ENDPOINTS from '../../../api/endpoints';
import { User, AuthTokens, LoginCredentials, RegisterCredentials } from '../../../shared/types';
import toast from 'react-hot-toast';

interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  
  // Actions
  login: (credentials: LoginCredentials) => Promise<void>;
  register: (credentials: RegisterCredentials) => Promise<void>;
  logout: () => void;
  updateUser: (user: User) => void;
  checkAuth: () => Promise<void>;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      isLoading: false,

      login: async (credentials: LoginCredentials) => {
        set({ isLoading: true });
        try {
          const response = await apiClient.post<AuthTokens & { user: User }>(
            API_ENDPOINTS.AUTH.LOGIN,
            credentials
          );
          
          const { access_token, refresh_token, user } = response.data;
          
          // Save tokens to localStorage
          localStorage.setItem('access_token', access_token);
          localStorage.setItem('refresh_token', refresh_token);
          
          set({
            user,
            accessToken: access_token,
            refreshToken: refresh_token,
            isAuthenticated: true,
            isLoading: false,
          });
          
          toast.success('Login successful!');
        } catch (error: any) {
          set({ isLoading: false });
          toast.error(error.response?.data?.error || 'Login failed');
          throw error;
        }
      },

      register: async (credentials: RegisterCredentials) => {
        set({ isLoading: true });
        try {
          const response = await apiClient.post<{ message: string; user: User }>(
            API_ENDPOINTS.AUTH.REGISTER,
            {
              email: credentials.email,
              password: credentials.password,
              full_name: credentials.full_name,
            }
          );
          
          toast.success('Registration successful! Please login.');
          set({ isLoading: false });
          
          // Auto-login after registration
          await get().login({
            email: credentials.email,
            password: credentials.password,
          });
        } catch (error: any) {
          set({ isLoading: false });
          toast.error(error.response?.data?.error || 'Registration failed');
          throw error;
        }
      },

      logout: () => {
        // Clear tokens from localStorage
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        
        // Clear state
        set({
          user: null,
          accessToken: null,
          refreshToken: null,
          isAuthenticated: false,
        });
        
        toast.success('Logged out successfully');
      },

      updateUser: (user: User) => {
        set({ user });
      },

      checkAuth: async () => {
        const token = localStorage.getItem('access_token');
        
        if (!token) {
          set({ isAuthenticated: false });
          return;
        }
        
        try {
          const response = await apiClient.get<User>(API_ENDPOINTS.PROFILE.GET);
          set({
            user: response.data,
            isAuthenticated: true,
            accessToken: token,
          });
        } catch (error) {
          // Token is invalid or expired
          get().logout();
        }
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);
