import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../features/auth/store/authStore';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../../shared/components/ui/Card';
import { Button } from '../../shared/components/ui/Button';
import { 
  FileText, 
  PenTool, 
  User, 
  TrendingUp,
  Clock,
  CheckCircle,
  AlertCircle
} from 'lucide-react';

const Dashboard = () => {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);

  const stats = [
    {
      title: 'Documents Created',
      value: '0',
      icon: FileText,
      color: 'text-blue-600',
      bgColor: 'bg-blue-100',
    },
    {
      title: 'Credits Remaining',
      value: user?.credits_remaining || 2,
      icon: TrendingUp,
      color: 'text-green-600',
      bgColor: 'bg-green-100',
    },
    {
      title: 'Profile Completion',
      value: '0%',
      icon: User,
      color: 'text-purple-600',
      bgColor: 'bg-purple-100',
    },
    {
      title: 'Account Type',
      value: user?.is_premium ? 'Premium' : 'Free',
      icon: user?.is_premium ? CheckCircle : AlertCircle,
      color: user?.is_premium ? 'text-yellow-600' : 'text-gray-600',
      bgColor: user?.is_premium ? 'bg-yellow-100' : 'bg-gray-100',
    },
  ];

  const quickActions = [
    {
      title: 'Generate Resume',
      description: 'Create a professional resume with AI',
      icon: FileText,
      action: () => navigate('/dashboard/generate/resume'),
      color: 'bg-blue-600 hover:bg-blue-700',
    },
    {
      title: 'Generate Cover Letter',
      description: 'Write a compelling cover letter',
      icon: PenTool,
      action: () => navigate('/dashboard/generate/cover-letter'),
      color: 'bg-green-600 hover:bg-green-700',
    },
    {
      title: 'Update Profile',
      description: 'Keep your information up to date',
      icon: User,
      action: () => navigate('/dashboard/profile'),
      color: 'bg-purple-600 hover:bg-purple-700',
    },
    {
      title: 'View History',
      description: 'See all your generated documents',
      icon: Clock,
      action: () => navigate('/dashboard/history'),
      color: 'bg-gray-600 hover:bg-gray-700',
    },
  ];

  return (
    <div className="space-y-6">
      {/* Welcome Section */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
          Welcome back, {user?.full_name || 'User'}!
        </h1>
        <p className="mt-2 text-gray-600 dark:text-gray-400">
          Ready to create your next professional document?
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat, index) => {
          const Icon = stat.icon;
          return (
            <Card key={index}>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">
                  {stat.title}
                </CardTitle>
                <div className={`p-2 rounded-full ${stat.bgColor}`}>
                  <Icon className={`h-4 w-4 ${stat.color}`} />
                </div>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stat.value}</div>
              </CardContent>
            </Card>
          );
        })}
      </div>

      {/* Quick Actions */}
      <div>
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          Quick Actions
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {quickActions.map((action, index) => {
            const Icon = action.icon;
            return (
              <Card 
                key={index}
                className="cursor-pointer hover:shadow-lg transition-shadow"
                onClick={action.action}
              >
                <CardHeader>
                  <div className={`inline-flex p-3 rounded-lg text-white ${action.color}`}>
                    <Icon className="h-6 w-6" />
                  </div>
                  <CardTitle className="mt-4">{action.title}</CardTitle>
                  <CardDescription>{action.description}</CardDescription>
                </CardHeader>
              </Card>
            );
          })}
        </div>
      </div>

      {/* Recent Documents (Placeholder) */}
      <Card>
        <CardHeader>
          <CardTitle>Recent Documents</CardTitle>
          <CardDescription>Your recently generated documents will appear here</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="text-center py-8 text-gray-500 dark:text-gray-400">
            <FileText className="h-12 w-12 mx-auto mb-4 opacity-50" />
            <p>No documents yet</p>
            <p className="text-sm mt-2">Start by generating your first resume or cover letter</p>
            <Button
              className="mt-4"
              onClick={() => navigate('/dashboard/generate/resume')}
            >
              Generate Your First Document
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Upgrade Banner for Free Users */}
      {!user?.is_premium && (
        <Card className="bg-gradient-to-r from-purple-600 to-blue-600">
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div className="text-white">
                <h3 className="text-xl font-bold mb-2">Upgrade to Premium</h3>
                <p className="opacity-90">
                  Get unlimited document generations and access to premium templates
                </p>
              </div>
              <Button
                variant="secondary"
                size="lg"
                onClick={() => navigate('/dashboard/upgrade')}
              >
                Upgrade Now
              </Button>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
};

export default Dashboard;
