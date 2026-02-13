package repository

import (
	"context"
	"log"
	"sort"
	"snackWeb/backend/internal/db"
	"snackWeb/backend/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// GetPosts returns a list of feed items with nested replies
func GetPosts(limit int) ([]models.FeedItem, error) {
	// Query GSI1 where GSI1PK="POST" (sorted by GSI1SK desc)
	out, err := db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("GSI1PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "POST"},
		},
		ScanIndexForward: aws.Bool(false), // DESC order (newest first)
		Limit:            aws.Int32(int32(limit)),
	})
	if err != nil {
		return nil, err
	}

	var posts []models.Post
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &posts); err != nil {
		return nil, err
	}

	var feed []models.FeedItem
	for _, p := range posts {
		item := models.FeedItem{
			ID:        p.ID,
			Timestamp: p.CreatedAt,
			EventType: p.EventType,
			AgentName: p.AgentName,
			Content:   p.Content,
			Topic:     p.Topic,
		}

		replies, err := GetRepliesForPost(p.ID)
		if err != nil {
			log.Println("Error fetching replies:", err)
			item.Replies = []models.APIReply{}
		} else {
			item.Replies = replies
		}
		feed = append(feed, item)
	}

	return feed, nil
}

func GetRepliesForPost(postID string) ([]models.APIReply, error) {
	// Query Main Table PK="POST#<ID>", SK begins_with "REPLY#"
	// Reply Item Schema: PK=POST#<ID>, SK=REPLY#<TIMESTAMP>
	
	pk := "POST#" + postID
	
	out, err := db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":        &types.AttributeValueMemberS{Value: pk},
			":sk_prefix": &types.AttributeValueMemberS{Value: "REPLY#"},
		},
	})
	if err != nil {
		return nil, err
	}

	var repliesDB []models.Reply
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &repliesDB); err != nil {
		return nil, err
	}

	var replies []models.APIReply
	for _, r := range repliesDB {
		replies = append(replies, models.APIReply{
			ID:        r.ID,
			Timestamp: r.CreatedAt,
			EventType: r.EventType,
			AgentName: r.AgentName, // Stored denormalized
			Content:   r.Content,
		})
	}
	
	// Sort ASC by timestamp if not already guaranteed by SK
	// SK is REPLY#<TIMESTAMP>, so it should be sorted.
	
	if replies == nil {
		replies = []models.APIReply{}
	}
	return replies, nil
}

func GetPersonas() ([]models.Persona, error) {
	// Query GSI1 where GSI1PK="PERSONA"
	out, err := db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("GSI1PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "PERSONA"},
		},
	})
	if err != nil {
		return nil, err
	}

	var personas []models.Persona
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &personas); err != nil {
		return nil, err
	}
	
	// Filter active in code if not in query (assuming we fetch all and filter)
	// Or use FilterExpression if GSI doesn't cover is_active
	var activePersonas []models.Persona
	for _, p := range personas {
		if p.IsActive {
			activePersonas = append(activePersonas, p)
		}
	}

	return activePersonas, nil
}

func GetGenerationStats() ([]models.GenerationStats, error) {
	// Query GSI1 where GSI1PK="STATS"
	out, err := db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("GSI1PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "STATS"},
		},
		ScanIndexForward: aws.Bool(true), // ASC
	})
	if err != nil {
		return nil, err
	}

	var stats []models.GenerationStats
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &stats); err != nil {
		return nil, err
	}
	
	// Sort just in case (though Query with sort key GSI1SK=<GEN_ID> should work if ID is comparable)
	// If GEN_ID is number, basic string sort might fail for 10 vs 2.
	// But assuming simple integer sort logic or small numbers for now.
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Generation < stats[j].Generation
	})

	return stats, nil
}
