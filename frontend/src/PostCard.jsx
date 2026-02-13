import React from 'react';

function formatTime(isoString) {
    try {
        const d = new Date(isoString);
        return d.toLocaleTimeString('ja-JP', { hour: '2-digit', minute: '2-digit' });
    } catch {
        return '';
    }
}

function getInitial(name) {
    return name ? name.charAt(0).toUpperCase() : '?';
}

function getTypeBadgeClass(type) {
    const t = (type || '').toLowerCase();
    if (t === 'post') return 'post-card__type-badge--post';
    if (t === 'reply') return 'post-card__type-badge--reply';
    if (t === 'reaction') return 'post-card__type-badge--reaction';
    if (t === 'pass') return 'post-card__type-badge--pass';
    return 'post-card__type-badge--post';
}

export default function PostCard({ item, replyCount = 0, isExpanded = false, onClick, isReply = false }) {
    const type = (item.event_type || 'post').toLowerCase();
    const isPass = type === 'pass';

    const cardClasses = [
        'post-card',
        isPass ? 'post-card--pass' : '',
        isReply ? 'post-card--reply' : '',
        !isReply && replyCount > 0 ? 'post-card--clickable' : '',
    ].filter(Boolean).join(' ');

    return (
        <div className={cardClasses} onClick={!isReply ? onClick : undefined}>
            <div className="post-card__header">
                <div className={`post-card__avatar ${isReply ? 'post-card__avatar--small' : ''}`}>
                    {getInitial(item.agent_name)}
                </div>
                <div className="post-card__meta">
                    <div className="post-card__author">{item.agent_name}</div>
                    <div className="post-card__time">{formatTime(item.timestamp)}</div>
                </div>
                <span className={`post-card__type-badge ${getTypeBadgeClass(type)}`}>
                    {type}
                </span>
            </div>

            {isPass ? (
                <div className="post-card__content" style={{ fontStyle: 'italic', opacity: 0.6 }}>
                    (Opted for silence)
                </div>
            ) : (
                <div className="post-card__content">{item.content}</div>
            )}

            {/* Reply count indicator — only on top-level posts */}
            {!isReply && replyCount > 0 && (
                <div className="post-card__reply-indicator">
                    <span className="post-card__reply-icon">{isExpanded ? '▾' : '▸'}</span>
                    <span>{replyCount} {replyCount === 1 ? 'reply' : 'replies'}</span>
                </div>
            )}
        </div>
    );
}
