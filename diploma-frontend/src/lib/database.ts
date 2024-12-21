import axios from 'axios';
import { isAuthenticated } from "./stores/authStore";
import { apiClient } from "./api";

const API_BASE_URL = 'http://localhost:5001';


export const getUserDatabase = async (): Promise<any> => {
    try {
        const response = await apiClient.get('/db/get-database');
        return response.data;
    } catch (error: any) {
        console.error('Error fetching user database:', error.response?.data || error.message);
        throw error;
    }
};

// Функция для создания таблицы
export const createTable = async (tableData: Record<string, any>): Promise<any> => {
    try {
        const response = await apiClient.post('/db/create-table', tableData);
        return response.data; // Возвращаем данные от сервера
    } catch (error: any) {
        console.error('Error creating table:', error.response?.data || error.message);
        throw error;
    }
};

// Функция для выполнения произвольного SQL-запроса
export const executeQuery = async (queryData: Record<string, any>): Promise<any> => {
    try {
        const response = await apiClient.post('/db/execute-query', queryData);
        return response.data; // Возвращаем данные от сервера
    } catch (error: any) {
        console.error('Error executing query:', error.response?.data || error.message);
        throw error;
    }
};