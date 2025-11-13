import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../../shared/components/ui/Card';

const CoverLetter = () => {
  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Generate Cover Letter</CardTitle>
          <CardDescription>Create a compelling cover letter tailored to your job application</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-gray-600 dark:text-gray-400">Cover letter generation form will be implemented here</p>
        </CardContent>
      </Card>
    </div>
  );
};

export default CoverLetter;
