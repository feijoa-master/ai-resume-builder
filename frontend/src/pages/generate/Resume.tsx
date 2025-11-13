import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../../shared/components/ui/Card';

const Resume = () => {
  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Generate Resume</CardTitle>
          <CardDescription>Create a professional resume with AI assistance</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-gray-600 dark:text-gray-400">Resume generation form will be implemented here</p>
        </CardContent>
      </Card>
    </div>
  );
};

export default Resume;
