import { createBrowserRouter } from 'react-router-dom';
import { lazy, Suspense } from 'react';

// Lazy load pages for code splitting
const LoginPage = lazy(() => import('../pages/auth/Login'));
const RegisterPage = lazy(() => import('../pages/auth/Register'));
const DashboardPage = lazy(() => import('../pages/dashboard/Dashboard'));
const ProfilePage = lazy(() => import('../pages/profile/Profile'));
const ResumePage = lazy(() => import('../pages/generate/Resume'));
const CoverLetterPage = lazy(() => import('../pages/generate/CoverLetter'));
const HistoryPage = lazy(() => import('../pages/history/History'));

// Layouts
import RootLayout from '../shared/components/layouts/RootLayout';
import AuthLayout from '../shared/components/layouts/AuthLayout';
import DashboardLayout from '../shared/components/layouts/DashboardLayout';

// Loading component
const PageLoader = () => (
  <div className="flex items-center justify-center h-screen">
    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
  </div>
);

// Wrap lazy loaded components with Suspense
const withSuspense = (Component: React.LazyExoticComponent<() => JSX.Element>) => {
  return (
    <Suspense fallback={<PageLoader />}>
      <Component />
    </Suspense>
  );
};

export const router = createBrowserRouter([
  {
    path: '/',
    element: <RootLayout />,
    children: [
      {
        index: true,
        element: withSuspense(DashboardPage),
      },
      {
        path: 'auth',
        element: <AuthLayout />,
        children: [
          {
            path: 'login',
            element: withSuspense(LoginPage),
          },
          {
            path: 'register',
            element: withSuspense(RegisterPage),
          },
        ],
      },
      {
        path: 'dashboard',
        element: <DashboardLayout />,
        children: [
          {
            index: true,
            element: withSuspense(DashboardPage),
          },
          {
            path: 'profile',
            element: withSuspense(ProfilePage),
          },
          {
            path: 'generate',
            children: [
              {
                path: 'resume',
                element: withSuspense(ResumePage),
              },
              {
                path: 'cover-letter',
                element: withSuspense(CoverLetterPage),
              },
            ],
          },
          {
            path: 'history',
            element: withSuspense(HistoryPage),
          },
        ],
      },
    ],
  },
]);
