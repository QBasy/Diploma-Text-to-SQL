import api from './index';

interface LoginRequest {
    email: string;
    password: string;
}

interface User {
    id: string;
    name: string;
    email: string;
}


interface RegisterRequest {
    name: string;
    email: string;
    password: string;
}

interface ForgotPasswordRequest {
    email: string;
}

interface AuthResponse {
    token: string;
    user: User;
}

export const login = async (data: LoginRequest): Promise<AuthResponse> => {
    return api.post('/auth/login', data);
};

export const register = async (data: RegisterRequest): Promise<AuthResponse> => {
    return api.post('/auth/register', data);
};

export const getCurrentUser = async (): Promise<User> => {
    return api.get('/auth/me');
};

export const forgotPassword = async (data: ForgotPasswordRequest): Promise<void> => {
    return api.post('/auth/forgot-password', data);
};

export const logout = async (): Promise<void> => {
    return api.post('/auth/logout');
};
