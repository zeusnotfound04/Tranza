import { cookies } from 'next/headers';
import { NextRequest, NextResponse } from 'next/server';

// App Router (Next.js 13+)
export async function GET() {
    try {
        console.log("Fetching token...");
        const cookieStore = await cookies();
        
        const accessToken = cookieStore.get('access_token');
        
        if (!accessToken) {
            return NextResponse.json({
                success: false,
                error: 'No token found'
            }, { status: 401 });
        }
        
        return NextResponse.json({
            success: true,
            data: {
                token: accessToken.value
            }
        });
    } catch (error) {
        console.error('Token API error:', error);
        return NextResponse.json({
            success: false,
            error: 'Internal server error'
        }, { status: 500 });
    }
}