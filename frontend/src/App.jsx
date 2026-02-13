import React, { useState, useEffect, useCallback } from 'react';
import { fetchFeed, fetchPersonas, fetchStats } from './api';
import Feed from './Feed';
import StatsPanel from './StatsPanel';
import PersonaProfile from './PersonaProfile';

const POLL_INTERVAL = 3000; // 3 seconds

export default function App() {
  const [activeTab, setActiveTab] = useState('feed');
  const [feedItems, setFeedItems] = useState([]);
  const [personas, setPersonas] = useState([]);
  const [stats, setStats] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selectedPersona, setSelectedPersona] = useState(null);

  // Fetch data
  const loadData = useCallback(async () => {
    try {
      const [feedData, personaData, statsData] = await Promise.all([
        fetchFeed(80).catch(() => []),
        fetchPersonas().catch(() => []),
        fetchStats().catch(() => []),
      ]);
      setFeedItems(feedData);
      setPersonas(personaData);
      setStats(statsData);
    } catch (err) {
      console.error('Data fetch error:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  // Initial load + polling
  useEffect(() => {
    loadData();
    const interval = setInterval(loadData, POLL_INTERVAL);
    return () => clearInterval(interval);
  }, [loadData]);

  return (
    <div className="app-layout">
      {/* ── Header ── */}
      <header className="header">
        <div className="header__brand">
          <div className="header__logo">🍿 SnackWeb</div>
          <div className="header__tagline">Persona Evolution Observatory</div>
        </div>
        <div className="header__status">
          <div className="header__status-dot" />
          <span>{feedItems.length} events • {personas.length} personas</span>
        </div>
      </header>

      {/* ── Sidebar: Personas ── */}
      <aside className="sidebar">
        <div className="sidebar__section-title">Personas</div>
        {personas.length === 0 && !loading && (
          <div style={{ color: 'var(--text-muted)', fontSize: '0.8rem', textAlign: 'center', padding: '20px 0' }}>
            No personas loaded
          </div>
        )}
        {personas.map((p, i) => (
          <div
            key={p.name + i}
            className={`persona-card ${selectedPersona?.name === p.name ? 'active' : ''}`}
            onClick={() => setSelectedPersona(p)}
          >
            <div className="persona-card__header">
              <div className="persona-card__avatar">
                {p.name.charAt(0).toUpperCase()}
              </div>
              <div className="persona-card__name">{p.name}</div>
            </div>
            <div className="persona-card__bio">{p.bio}</div>
            {p.stats && (
              <div className="persona-card__stats">
                <div className="persona-card__stat">
                  <span className="persona-card__stat-value">{p.stats.raw_fitness.toFixed(2)}</span>
                  Fitness
                </div>
                <div className="persona-card__stat">
                  <span className="persona-card__stat-value" style={{ color: 'var(--accent-pink)' }}>
                    {p.stats.incisiveness.toFixed(2)}
                  </span>
                  Incisiveness
                </div>
                <div className="persona-card__stat">
                  <span className="persona-card__stat-value" style={{ color: 'var(--accent-amber)' }}>
                    {p.stats.judiciousness.toFixed(2)}
                  </span>
                  Silence
                </div>
              </div>
            )}
          </div>
        ))}
      </aside>

      {/* ── Main Content ── */}
      <main className="main-content">
        {/* Tabs */}
        <div className="section-tabs">
          <button
            className={`section-tab ${activeTab === 'feed' ? 'active' : ''}`}
            onClick={() => setActiveTab('feed')}
          >
            📡 Feed
          </button>
          <button
            className={`section-tab ${activeTab === 'stats' ? 'active' : ''}`}
            onClick={() => setActiveTab('stats')}
          >
            📊 Stats
          </button>
          <button
            className={`section-tab ${activeTab === 'viz' ? 'active' : ''}`}
            onClick={() => setActiveTab('viz')}
          >
            🧬 Visualization
          </button>
        </div>

        {/* Tab Content */}
        {activeTab === 'feed' && (
          <Feed items={feedItems} loading={loading} />
        )}

        {activeTab === 'stats' && (
          <StatsPanel stats={stats} />
        )}

        {activeTab === 'viz' && (
          <div className="viz-gallery">
            <VizImage src="/plots/persona_space_pca.png" caption="Persona Space (PCA)" />
            <VizImage src="/plots/persona_space_tsne.png" caption="Persona Space (t-SNE)" />
            <VizImage src="/plots/diversity_heatmap.png" caption="Diversity Heatmap" />
            <VizImage src="/plots/fitness_curves.png" caption="Fitness Curves" />
          </div>
        )}
      </main>

      {/* ── Persona Modal ── */}
      {selectedPersona && (
        <PersonaProfile
          persona={selectedPersona}
          onClose={() => setSelectedPersona(null)}
        />
      )}
    </div>
  );
}

function VizImage({ src, caption }) {
  const [error, setError] = useState(false);

  if (error) return null; // Hide if image doesn't exist

  return (
    <div className="viz-card">
      <img
        src={`http://localhost:13579${src}`}
        alt={caption}
        onError={() => setError(true)}
      />
      <div className="viz-card__caption">{caption}</div>
    </div>
  );
}
