import React from 'react';

export default function StatsPanel({ stats }) {
    if (!stats || stats.length === 0) {
        return (
            <div className="empty-state">
                <div className="empty-state__icon">📊</div>
                <div className="empty-state__title">No statistical data</div>
                <div className="empty-state__desc">
                    Run the evolution simulation to see generation-level stats.
                </div>
            </div>
        );
    }

    const latest = stats[stats.length - 1];

    return (
        <div>
            {/* Summary cards */}
            <div className="stats-panel">
                <div className="stat-card">
                    <div className="stat-card__label">GEN</div>
                    <div className="stat-card__value">{latest.generation}</div>
                    <div className="stat-card__subtitle">Latest Generation</div>
                </div>
                <div className="stat-card">
                    <div className="stat-card__label">DIVERSITY</div>
                    <div className="stat-card__value">{latest.population_diversity.toFixed(2)}</div>
                    <div className="stat-card__subtitle">Population Diversity</div>
                </div>
                <div className="stat-card">
                    <div className="stat-card__label">MEAN FITNESS</div>
                    <div className="stat-card__value">{latest.fitness_mean.toFixed(2)}</div>
                    <div className="stat-card__subtitle">Mean Fitness</div>
                </div>
            </div>

            {/* History table */}
            <div style={{
                background: 'var(--bg-card)',
                border: '1px solid var(--border-subtle)',
                borderRadius: 'var(--radius-lg)',
                overflow: 'hidden',
            }}>
                <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.8rem' }}>
                    <thead>
                        <tr style={{ borderBottom: '1px solid var(--border-subtle)' }}>
                            <th style={thStyle}>Gen</th>
                            <th style={thStyle}>Diversity</th>
                            <th style={thStyle}>Fitness Mean</th>
                        </tr>
                    </thead>
                    <tbody>
                        {stats.map((s, i) => (
                            <tr key={i} style={{ borderBottom: '1px solid var(--border-subtle)' }}>
                                <td style={tdStyle}>{s.generation}</td>
                                <td style={tdStyle}>{s.population_diversity.toFixed(3)}</td>
                                <td style={tdStyle}>{s.fitness_mean.toFixed(3)}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}

const thStyle = {
    padding: '12px 16px',
    textAlign: 'left',
    color: 'var(--text-muted)',
    fontWeight: 600,
    textTransform: 'uppercase',
    letterSpacing: '0.5px',
    fontSize: '0.7rem',
};

const tdStyle = {
    padding: '10px 16px',
    color: 'var(--text-secondary)',
};
