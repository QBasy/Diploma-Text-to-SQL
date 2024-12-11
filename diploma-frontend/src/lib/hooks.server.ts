import { redirect, type Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
    const protectedRoutes = ['/query', '/items']; // Список защищённых маршрутов
    const token = event.cookies.get('token');

    if (protectedRoutes.includes(event.url.pathname) && !token) {
        throw redirect(302, '/auth'); // Перенаправление на страницу входа
    }

    return resolve(event);
};
