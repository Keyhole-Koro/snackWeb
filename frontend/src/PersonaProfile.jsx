import React from 'react';

function getInitial(name) {
    return name ? name.charAt(0).toUpperCase() : '?';
}

export default function PersonaProfile({ persona, onClose }) {
    if (!persona) return null;

    return (
        <div style={{
            position: 'fixed', inset: 0, zIndex: 200,
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            background: 'rgba(0,0,0,0.7)', backdropFilter: 'blur(8px)',
        }} onClick={onClose}>
            <div style={{
                background: 'var(--bg-secondary)',
                border: '1px solid var(--border-accent)',
                borderRadius: 'var(--radius-xl)',
                padding: '32px',
                maxWidth: '480px',
                width: '90%',
                boxShadow: 'var(--shadow-glow-cyan)',
                animation: 'fadeInUp 0.3s ease-out',
            }} onClick={(e) => e.stopPropagation()}>
                {/* Header */}
                <div style={{ display: 'flex', alignItems: 'center', gap: '16px', marginBottom: '20px' }}>
                    <div className="post-card__avatar" style={{ width: '56px', height: '56px', fontSize: '1.4rem' }}>
                        {getInitial(persona.name)}
                    </div>
                    <div>
                        <h2 style={{ fontSize: '1.3rem', fontWeight: 700, margin: 0 }}>{persona.name}</h2>
                        <div style={{ color: 'var(--text-muted)', fontSize: '0.8rem' }}>Persona Agent</div>
                    </div>
                </div>

                {/* Bio */}
                <div style={{
                    background: 'var(--bg-glass)',
                    borderRadius: 'var(--radius-md)',
                    padding: '16px',
                    marginBottom: '20px',
                    border: '1px solid var(--border-subtle)',
                }}>
                    <div style={{ fontSize: '0.7rem', color: 'var(--text-muted)', textTransform: 'uppercase', letterSpacing: '1px', marginBottom: '8px' }}>
                        Bio / ストーリー
                    </div>
                    <div style={{ fontSize: '0.88rem', lineHeight: 1.7, color: 'var(--text-secondary)' }}>
                        {persona.bio}
                    </div>
                </div>

                {/* Stats */}
                {persona.stats && (
                    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '12px' }}>
                        <StatBox label="Post Quality" value={persona.stats.post_quality} />
                        <StatBox label="Incisiveness" value={persona.stats.incisiveness} accent="var(--accent-pink)" />
                        <StatBox label="Judiciousness" value={persona.stats.judiciousness} accent="var(--accent-amber)" />
                        <StatBox label="Raw Fitness" value={persona.stats.raw_fitness} accent="var(--accent-green)" />
                    </div>
                )}

                {/* Close */}
                <button onClick={onClose} style={{
                    width: '100%', marginTop: '20px', padding: '10px',
                    background: 'var(--bg-glass)', border: '1px solid var(--border-subtle)',
                    borderRadius: 'var(--radius-sm)', color: 'var(--text-secondary)',
                    cursor: 'pointer', fontFamily: 'inherit', fontSize: '0.85rem',
                }}>
                    閉じる
                </button>
            </div>
        </div>
    );
}

function StatBox({ label, value, accent }) {
    return (
        <div style={{
            background: 'var(--bg-glass)', borderRadius: 'var(--radius-sm)',
            padding: '12px', border: '1px solid var(--border-subtle)', textAlign: 'center',
        }}>
            <div style={{ fontSize: '0.65rem', color: 'var(--text-muted)', textTransform: 'uppercase', letterSpacing: '1px' }}>
                {label}
            </div>
            <div style={{ fontSize: '1.3rem', fontWeight: 700, color: accent || 'var(--accent-cyan)', marginTop: '4px' }}>
                {typeof value === 'number' ? value.toFixed(2) : value}
            </div>
        </div>
    );
}
