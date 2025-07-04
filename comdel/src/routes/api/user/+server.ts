// src/routes/api/user/info/+server.ts
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ cookies }) => {
    const token = cookies.get('jwt');

    const response = await fetch("http://localhost:8080/user/info", {
        headers: {
            Cookie: `jwt=${token}`
        },
        method: "GET"
    });

    const userInfo = await response.json();

    return new Response(JSON.stringify(userInfo), {
        headers: { 'Content-Type': 'application/json' }
    });
};
