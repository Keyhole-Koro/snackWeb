import React, { useState } from 'react';
import PostCard from './PostCard';

/**
 * Groups feed items into threads:
 *  - Top-level: POST items (or items with no related_to)
 *  - Replies: REPLY/REACTION items grouped by related_to (matching agent_name of a post)
 */
function buildThreads(items) {
    const threads = [];
    const replyMap = {}; // agent_name -> replies[]

    // First pass: separate posts from replies
    for (const item of items) {
        const type = (item.event_type || '').toLowerCase();
        if (type === 'post' || (!item.related_to && type !== 'pass')) {
            threads.push({ post: item, replies: [] });
        }
    }

    // Second pass: attach replies to their parent post's thread
    for (const item of items) {
        const type = (item.event_type || '').toLowerCase();
        if (type === 'reply' || type === 'reaction') {
            const parentThread = threads.find(t => t.post.agent_name === item.related_to);
            if (parentThread) {
                parentThread.replies.push(item);
            } else {
                // Orphan reply — show as top-level
                threads.push({ post: item, replies: [] });
            }
        } else if (type === 'pass') {
            threads.push({ post: item, replies: [] });
        }
    }

    return threads;
}

export default function Feed({ items, loading }) {
    const [expandedIdx, setExpandedIdx] = useState(null);

    if (loading) {
        return <div className="spinner" />;
    }

    if (!items || items.length === 0) {
        return (
            <div className="empty-state">
                <div className="empty-state__icon">📡</div>
                <div className="empty-state__title">フィードが空です</div>
                <div className="empty-state__desc">
                    シミュレーションを実行すると、ペルソナの投稿がここに表示されます。
                </div>
            </div>
        );
    }

    const threads = buildThreads(items);

    return (
        <div>
            {threads.map((thread, i) => (
                <div key={`thread-${i}`}>
                    <PostCard
                        item={thread.post}
                        replyCount={thread.replies.length}
                        isExpanded={expandedIdx === i}
                        onClick={() => setExpandedIdx(expandedIdx === i ? null : i)}
                    />
                    {/* Threaded replies — only shown when expanded */}
                    {expandedIdx === i && thread.replies.length > 0 && (
                        <div className="thread-replies">
                            {thread.replies.map((reply, j) => (
                                <PostCard key={`reply-${i}-${j}`} item={reply} isReply />
                            ))}
                        </div>
                    )}
                </div>
            ))}
        </div>
    );
}
