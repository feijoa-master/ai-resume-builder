import { Outlet, Navigate } from 'react-router-dom';
import { useAuthStore } from '../../../features/auth/store/authStore';

const AuthLayout = () => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  // If already authenticated, redirect to dashboard
  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <h2 className="text-3xl font-extrabold text-gray-900 dark:text-white">
            AI Resume Builder
          </h2>
          <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
            Create professional resumes and cover letters with AI
          </p>
        </div>
        <Outlet />
      </div>
    </div>
  );
};

export default AuthLayout;
