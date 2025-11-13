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
const registerSchema = z.object({
  full_name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.string().email('Invalid email address'),
  password: z.string()
    .min(8, 'Password must be at least 8 characters')
    .regex(/[A-Z]/, 'Password must contain at least one uppercase letter')
    .regex(/[0-9]/, 'Password must contain at least one number'),
  confirm_password: z.string(),
}).refine((data) => data.password === data.confirm_password, {
  message: "Passwords don't match",
  path: ['confirm_password'],
});

type RegisterFormData = z.infer<typeof registerSchema>;

const Register = () => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const register = useAuthStore((state) => state.register);
  
  const {
    register: registerField,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
  });

  const onSubmit = async (data: RegisterFormData) => {
    setIsLoading(true);
    try {
      await register(data);
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
        <CardTitle className="text-2xl text-center">Create an account</CardTitle>
        <CardDescription className="text-center">
          Enter your details to create your account
        </CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit(onSubmit)}>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Input
              {...registerField('full_name')}
              label="Full Name"
              placeholder="John Doe"
              error={!!errors.full_name}
              helperText={errors.full_name?.message}
              disabled={isLoading}
            />
          </div>
          <div className="space-y-2">
            <Input
              {...registerField('email')}
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
              {...registerField('password')}
              type="password"
              label="Password"
              placeholder="Create a strong password"
              error={!!errors.password}
              helperText={errors.password?.message}
              disabled={isLoading}
            />
          </div>
          <div className="space-y-2">
            <Input
              {...registerField('confirm_password')}
              type="password"
              label="Confirm Password"
              placeholder="Confirm your password"
              error={!!errors.confirm_password}
              helperText={errors.confirm_password?.message}
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
            Create Account
          </Button>
          <div className="text-sm text-center">
            <p className="text-gray-600 dark:text-gray-400">
              Already have an account?{' '}
              <Link
                to="/auth/login"
                className="text-primary hover:underline"
              >
                Sign in
              </Link>
            </p>
          </div>
          <div className="text-xs text-center text-gray-500 dark:text-gray-400">
            By creating an account, you agree to our{' '}
            <Link to="/terms" className="text-primary hover:underline">
              Terms of Service
            </Link>{' '}
            and{' '}
            <Link to="/privacy" className="text-primary hover:underline">
              Privacy Policy
            </Link>
          </div>
        </CardFooter>
      </form>
    </Card>
  );
};

export default Register;
