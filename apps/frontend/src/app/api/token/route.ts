import { cookies } from 'next/headers';

// App Router (Next.js 13+)
export async function GET() {
    const cookieStore = await cookies();
    const accessToken = cookieStore.get('access_token');
    
    return Response.json({ token: accessToken?.value });
}