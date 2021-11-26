package bot

import (
	"context"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Backend struct {
	Driver *mongo.Client
	Config *BackendConf
}

type UserBack struct {
	HLogin  string
	Command string
}

func NewBackend(back BackendConf) *Backend {
	return &Backend{
		Config: &back,
	}
}

func (b *Backend) Init() error {
	ctx, err := b.Open()
	if err != nil {
		return err
	}
	defer b.Close(ctx)
	// db := b.Driver.Database(b.Config.DatabaseName)
	// db.Collection("tmp")
	// db.Collection("questions")
	return nil
}

func (b *Backend) Open() (context.Context, error) {
	clientOptions := options.Client().ApplyURI(b.Config.Connection)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	b.Driver = client
	return ctx, nil
}

func (b *Backend) Close(ctx context.Context) error {
	return b.Driver.Disconnect(ctx)
}

func (b *Backend) GetLastCommand(hlogin string) (string, error) {
	ctx, err := b.Open()
	if err != nil {
		return "", err
	}
	defer b.Close(ctx)
	db := b.Driver.Database(b.Config.DatabaseName)
	backend := db.Collection("tmp")

	filter := bson.D{
		{Key: "HLogin", Value: hlogin},
	}
	var user UserBack
	if err := backend.FindOne(ctx, filter).Decode(&user); err != nil {
		if err.Error() != "mongo: no documents in result" {
			return "", err
		}
	}

	if _, err := backend.DeleteMany(ctx, filter); err != nil {
		return "", err
	}

	return user.Command, nil
}

// func (b *Backend) FlashCommands(hlogin string) error {
// 	ctx, err := b.Open()
// 	if err != nil {
// 		return err
// 	}
// 	defer b.Close(ctx)
// 	db := b.Driver.Database(b.Config.DatabaseName)
// 	backend := db.Collection("tmp")

// 	filter := bson.D{
// 		{Key: "HLogin", Value: hlogin},
// 	}

// 	if _, err := backend.DeleteMany(ctx, filter); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (b *Backend) PutCommand(hlogin, command string) error {
	ctx, err := b.Open()
	if err != nil {
		return err
	}
	defer b.Close(ctx)
	db := b.Driver.Database(b.Config.DatabaseName)
	backend := db.Collection("tmp")

	user := UserBack{
		HLogin:  hlogin,
		Command: command,
	}
	if _, err := backend.InsertOne(ctx, bson.D{
		{Key: "HLogin", Value: user.HLogin},
		{Key: "Command", Value: user.Command},
	}); err != nil {
		return err
	}
	return nil
}

func (b *Backend) PutPage(data Page, url string) error {
	ctx, err := b.Open()
	if err != nil {
		return err
	}
	defer b.Close(ctx)
	db := b.Driver.Database(b.Config.DatabaseName)
	backend := db.Collection("questions")
	filter := bson.D{
		{Key: "url", Value: url},
	}
	if _, err := backend.UpdateOne(ctx, filter, bson.M{
		"$set": bson.M{
			"data": data,
		},
	}); err != nil {
		return err
	}
	return nil
}

func (b *Backend) PutNewPage(name, url string) error {
	ctx, err := b.Open()
	if err != nil {
		return err
	}
	defer b.Close(ctx)
	db := b.Driver.Database(b.Config.DatabaseName)
	backend := db.Collection("questions")
	page := bson.D{
		{Key: "name", Value: name},
		{Key: "url", Value: url},
		{Key: "data", Value: Page{}},
	}
	if _, err := backend.InsertOne(ctx, page); err != nil {
		return err
	}
	return nil
}

func (b *Backend) GetJsonByUrl(url string) (Page, error) {
	ctx, err := b.Open()
	if err != nil {
		return Page{}, err
	}
	defer b.Close(ctx)
	db := b.Driver.Database(b.Config.DatabaseName)
	backend := db.Collection("questions")
	filter := bson.D{
		{Key: "url", Value: url},
	}
	var data Page
	if err := backend.FindOne(ctx, filter).Decode(&data); err != nil {
		if err.Error() != "mongo: no documents in result" {
			return Page{}, err
		}
	}
	return data, err
}

func (b *Backend) GetJsonByName(name string) (Page, error) {
	ctx, err := b.Open()
	if err != nil {
		return Page{}, err
	}
	defer b.Close(ctx)
	db := b.Driver.Database(b.Config.DatabaseName)
	backend := db.Collection("questions")
	filter := bson.D{
		{Key: "name", Value: name},
	}
	var data Page
	if err := backend.FindOne(ctx, filter).Decode(&data); err != nil {
		if err.Error() != "mongo: no documents in result" {
			return Page{}, err
		}
	}
	return data, err
}

func (b *Backend) GetAllServiceName() ([]string, error) {
	ctx, err := b.Open()
	if err != nil {
		return []string{}, err
	}
	defer b.Close(ctx)
	db := b.Driver.Database(b.Config.DatabaseName)
	backend := db.Collection("questions")

	cur, err := backend.Find(ctx, bson.D{})
	var names []string
	for cur.Next(ctx) {
		data := struct {
			name string
			url  string
			page Page
		}{}
		if err := cur.Decode(&data); err != nil {
			return []string{}, err
		}
		names = append(names, data.name)
	}
	if err := cur.Err(); err != nil {
		return []string{}, err
	}
	return names, err
}

func (b *Backend) GetPageUrls() ([]string, error) {
	ctx, err := b.Open()
	if err != nil {
		return []string{}, err
	}
	defer b.Close(ctx)
	db := b.Driver.Database(b.Config.DatabaseName)
	backend := db.Collection("questions")

	cur, err := backend.Find(ctx, bson.D{})
	var urls []string
	for cur.Next(ctx) {
		data := struct {
			name string
			url  string
			page Page
		}{}
		if err := cur.Decode(&data); err != nil {
			return []string{}, err
		}
		urls = append(urls, data.url)
	}
	if err := cur.Err(); err != nil {
		return []string{}, err
	}
	return urls, err
}
