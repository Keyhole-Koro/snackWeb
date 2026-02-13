const API_BASE = 'http://localhost:8000';

export async function fetchFeed(limit = 50) {
    const res = await fetch(`${API_BASE}/api/feed?limit=${limit}`);
    if (!res.ok) throw new Error('Failed to fetch feed');
    return res.json();
}

export async function fetchPersonas() {
    const res = await fetch(`${API_BASE}/api/personas`);
    if (!res.ok) throw new Error('Failed to fetch personas');
    return res.json();
}

export async function fetchStats() {
    const res = await fetch(`${API_BASE}/api/stats`);
    if (!res.ok) throw new Error('Failed to fetch stats');
    return res.json();
}
