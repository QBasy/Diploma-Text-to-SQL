import { writable } from 'svelte/store';
import { getHistory, clearHistory } from '$lib/api/history';

interface QueryHistory {
    id: string;
    type: 'SQL' | 'Natural Language';
    query: string;
    result: any;
    timestamp: string;
}

export const historyStore = writable<QueryHistory[]>([]);

export const loadHistory = async () => {
    try {
        const history = await getHistory();
        historyStore.set(history); // Сохраняем историю
    } catch (error) {
        throw error;
    }
};

export const clearHistoryStore = async () => {
    try {
        await clearHistory();
        historyStore.set([]); // Очищаем историю
    } catch (error) {
        throw error;
    }
};