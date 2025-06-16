package db

import (
	"context"
	"math/rand"

	"fmt"
	"log"

	"github.com/amir-amirov/go-social-media/internal/store"
)

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)

	for _, user := range users {
		if err := store.Users.Create(ctx, &user); err != nil {
			log.Println("Error creating a user:", err)
			return
		}
	}

	posts := generatePosts(200, users)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, &post); err != nil {
			log.Println("Error creating a post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, &comment); err != nil {
			log.Println("Error creating a comment:", err)
			return
		}
	}

}

func generateUsers(num int) []store.User {
	users := make([]store.User, num)

	for i := 0; i < num; i++ {
		users[i] = store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%v", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%v", i) + "@example.com",
			Password: "123123",
		}
	}

	return users
}

func generatePosts(num int, users []store.User) []store.Post {
	posts := make([]store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = store.Post{
			UserID:  user.ID,
			Title:   titles[i%len(titles)],
			Content: contents[i%len(contents)],
			Tags: []string{
				tags[i%len(contents)],
				tags[(i+1)%len(contents)],
				tags[(i+2)%len(contents)],
			},
		}

	}

	return posts
}

func generateComments(num int, users []store.User, posts []store.Post) []store.Comment {
	comments := make([]store.Comment, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(users))]

		comments[i] = store.Comment{
			UserID:  user.ID,
			PostID:  post.ID,
			Content: post_comments[i%len(post_comments)],
		}

	}

	return comments
}
