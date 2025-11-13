import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../../shared/components/ui/Card';

const History = () => {
  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Document History</CardTitle>
          <CardDescription>View and manage all your generated documents</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-gray-600 dark:text-gray-400">Document history will be displayed here</p>
        </CardContent>
      </Card>
    </div>
  );
};

export default History;
