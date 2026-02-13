import React, { useState } from 'react';
import PostCard from './PostCard';

/**
 * Feed component — renders posts from the DB-backed API.
 *
 * The API now returns posts with nested `replies` arrays,
 * so we no longer need to group/thread client-side.
 */
export default function Feed({ items, loading }) {
    const [expandedIdx, setExpandedIdx] = useState(null);

    if (loading) {
        return <div className="spinner" />;
    }

    if (!items || items.length === 0) {
        return (
            <div className="empty-state">
                <div className="empty-state__icon">📡</div>
                <div className="empty-state__title">Feed is empty</div>
                <div className="empty-state__desc">
                    Run the simulation to see persona posts here.
                </div>
            </div>
        );
    }

    return (
        <div>
            {items.map((post, i) => (
                <div key={post.id || `post-${i}`}>
                    <PostCard
                        item={post}
                        replyCount={(post.replies || []).length}
                        isExpanded={expandedIdx === i}
                        onClick={() => setExpandedIdx(expandedIdx === i ? null : i)}
                    />
                    {expandedIdx === i && post.replies && post.replies.length > 0 && (
                        <div className="thread-replies">
                            {post.replies.map((reply, j) => (
                                <PostCard key={reply.id || `reply-${j}`} item={reply} isReply />
                            ))}
                        </div>
                    )}
                </div>
            ))}
        </div>
    );
}
