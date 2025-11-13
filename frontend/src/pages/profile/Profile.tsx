import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../../shared/components/ui/Card';

const Profile = () => {
  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Profile Management</CardTitle>
          <CardDescription>Update your personal information and experience</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-gray-600 dark:text-gray-400">Profile form will be implemented here</p>
        </CardContent>
      </Card>
    </div>
  );
};

export default Profile;
