import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '../../shared/components/ui/Button';
import { Input } from '../../shared/components/ui/Input';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '../../shared/components/ui/Card';
import { useAuthStore } from '../../features/auth/store/authStore';

// Validation schema
const loginSchema = z.object({
  email: z.string().email('Invalid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
});

type LoginFormData = z.infer<typeof loginSchema>;

const Login = () => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const login = useAuthStore((state) => state.login);
  
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  });

  const onSubmit = async (data: LoginFormData) => {
    setIsLoading(true);
    try {
      await login(data);
      navigate('/dashboard');
    } catch (error) {
      // Error is handled in the store with toast
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card>
      <CardHeader className="space-y-1">
        <CardTitle className="text-2xl text-center">Sign in</CardTitle>
        <CardDescription className="text-center">
          Enter your email and password to sign in
        </CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit(onSubmit)}>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Input
              {...register('email')}
              type="email"
              label="Email"
              placeholder="john@example.com"
              error={!!errors.email}
              helperText={errors.email?.message}
              disabled={isLoading}
            />
          </div>
          <div className="space-y-2">
            <Input
              {...register('password')}
              type="password"
              label="Password"
              placeholder="Enter your password"
              error={!!errors.password}
              helperText={errors.password?.message}
              disabled={isLoading}
            />
          </div>
        </CardContent>
        <CardFooter className="flex flex-col space-y-4">
          <Button 
            type="submit" 
            className="w-full" 
            loading={isLoading}
            disabled={isLoading}
          >
            Sign in
          </Button>
          <div className="text-sm text-center space-y-2">
            <Link
              to="/auth/forgot-password"
              className="text-primary hover:underline"
            >
              Forgot your password?
            </Link>
            <p className="text-gray-600 dark:text-gray-400">
              Don't have an account?{' '}
              <Link
                to="/auth/register"
                className="text-primary hover:underline"
              >
                Sign up
              </Link>
            </p>
          </div>
        </CardFooter>
      </form>
    </Card>
  );
};

export default Login;
