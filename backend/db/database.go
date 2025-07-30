package db

import (
	"context"
	"go-chat-backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// Initialize connects to MongoDB and returns a MongoDB instance
func Initialize(mongoURI, dbName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	database := client.Database(dbName)

	// Create MongoDB instance
	mongoDB := &MongoDB{
		Client:   client,
		Database: database,
	}

	// Initialize collections and indexes
	if err := mongoDB.initializeCollections(); err != nil {
		return nil, err
	}

	return mongoDB, nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

// initializeCollections creates collections and indexes
func (m *MongoDB) initializeCollections() error {
	ctx := context.Background()

	// Create indexes for rooms collection
	roomsCollection := m.Database.Collection("rooms")
	_, err := roomsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{"name", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// Create indexes for users collection
	usersCollection := m.Database.Collection("users")
	_, err = usersCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{"username", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{"email", 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return err
	}

	// Create indexes for messages collection
	messagesCollection := m.Database.Collection("messages")
	_, err = messagesCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{"room_id", 1}, {"created_at", -1}},
		},
		{
			Keys: bson.D{{"created_at", -1}},
		},
	})
	if err != nil {
		return err
	}

	// Insert default rooms if they don't exist
	return m.insertDefaultRooms()
}

// insertDefaultRooms creates default chat rooms
func (m *MongoDB) insertDefaultRooms() error {
	ctx := context.Background()
	roomsCollection := m.Database.Collection("rooms")

	defaultRooms := []models.Room{
		{
			Name:        "general",
			Description: "General chat room for everyone",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "random",
			Description: "Random discussions",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "tech",
			Description: "Technology discussions",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, room := range defaultRooms {
		filter := bson.M{"name": room.Name}
		update := bson.M{
			"$setOnInsert": room,
		}
		opts := options.Update().SetUpsert(true)
		_, err := roomsCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetRooms retrieves all chat rooms
func (m *MongoDB) GetRooms() ([]models.Room, error) {
	ctx := context.Background()
	collection := m.Database.Collection("rooms")

	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{"name", 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rooms []models.Room
	if err := cursor.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

// CreateRoom creates a new chat room
func (m *MongoDB) CreateRoom(room models.Room) (*models.Room, error) {
	ctx := context.Background()
	collection := m.Database.Collection("rooms")

	room.ID = primitive.NewObjectID()
	room.CreatedAt = time.Now()
	room.UpdatedAt = time.Now()

	result, err := collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	room.ID = result.InsertedID.(primitive.ObjectID)
	return &room, nil
}

// GetRoomByID retrieves a room by its ID
func (m *MongoDB) GetRoomByID(roomID string) (*models.Room, error) {
	ctx := context.Background()
	collection := m.Database.Collection("rooms")

	objectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, err
	}

	var room models.Room
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&room)
	if err != nil {
		return nil, err
	}

	return &room, nil
}

// SaveMessage saves a message to the database
func (m *MongoDB) SaveMessage(message models.Message) (*models.Message, error) {
	ctx := context.Background()
	collection := m.Database.Collection("messages")

	message.ID = primitive.NewObjectID()
	message.CreatedAt = time.Now()

	result, err := collection.InsertOne(ctx, message)
	if err != nil {
		return nil, err
	}

	message.ID = result.InsertedID.(primitive.ObjectID)
	return &message, nil
}

// GetRoomMessages retrieves messages for a specific room
func (m *MongoDB) GetRoomMessages(roomID string, limit int) ([]models.Message, error) {
	ctx := context.Background()
	collection := m.Database.Collection("messages")

	objectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, err
	}

	opts := options.Find().
		SetSort(bson.D{{"created_at", -1}}).
		SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, bson.M{"room_id": objectID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// CreateUser creates a new user
func (m *MongoDB) CreateUser(user models.User) (*models.User, error) {
	ctx := context.Background()
	collection := m.Database.Collection("users")

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (m *MongoDB) GetUserByUsername(username string) (*models.User, error) {
	ctx := context.Background()
	collection := m.Database.Collection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (m *MongoDB) GetUserByEmail(email string) (*models.User, error) {
	ctx := context.Background()
	collection := m.Database.Collection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (m *MongoDB) GetUserByID(userID string) (*models.User, error) {
	ctx := context.Background()
	collection := m.Database.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
