package change

import (
	"database/sql"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"gorm.io/gorm"
)

// FromDomainPost 将领域层对象转为dao层对象
func FromDomainPost(p domain.Post) dao.Post {
	return dao.Post{
		Model: gorm.Model{
			ID:        p.ID,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			DeletedAt: gorm.DeletedAt(p.DeletedAt),
		},
		Title:        p.Title,
		Content:      p.Content,
		Uid:          p.Uid,
		Status:       p.Status,
		PlateID:      p.PlateID,
		Slug:         p.Slug,
		CategoryID:   p.CategoryID,
		Tags:         p.Tags,
		CommentCount: p.CommentCount,
		IsSubmit:     p.IsSubmit,
	}
}

// FromDomainSlicePubPostList 将dao层对象转为领域层对象
func FromDomainSlicePubPostList(posts []dao.PubPost) []domain.Post {
	if len(posts) == 0 {
		return []domain.Post{}
	}

	domainPosts := make([]domain.Post, len(posts))
	for i, post := range posts {
		domainPosts[i] = domain.Post{
			ID:           post.ID,
			Title:        post.Title,
			Content:      post.Content,
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
			Status:       post.Status,
			Uid:          post.Uid,
			PlateID:      post.PlateID,
			Slug:         post.Slug,
			CategoryID:   post.CategoryID,
			Tags:         post.Tags,
			CommentCount: post.CommentCount,
		}
	}
	return domainPosts
}

// FromDomainSlicePost 将dao层对象转为领域层对象
func FromDomainSlicePost(posts []dao.Post) []domain.Post {
	if len(posts) == 0 {
		return []domain.Post{}
	}

	domainPosts := make([]domain.Post, len(posts))
	for i, post := range posts {
		domainPosts[i] = domain.Post{
			ID:           post.ID,
			Title:        post.Title,
			Content:      post.Content,
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
			DeletedAt:    sql.NullTime(post.DeletedAt),
			Status:       post.Status,
			Uid:          post.Uid,
			PlateID:      post.PlateID,
			Slug:         post.Slug,
			CategoryID:   post.CategoryID,
			Tags:         post.Tags,
			CommentCount: post.CommentCount,
		}
	}
	return domainPosts
}

// ToDomainPost 将dao层转化为领域层
func ToDomainPost(post dao.Post) domain.Post {
	return domain.Post{
		ID:           post.ID,
		Title:        post.Title,
		Content:      post.Content,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
		DeletedAt:    sql.NullTime(post.DeletedAt),
		Status:       post.Status,
		Uid:          post.Uid,
		PlateID:      post.PlateID,
		Slug:         post.Slug,
		CategoryID:   post.CategoryID,
		Tags:         post.Tags,
		CommentCount: post.CommentCount,
		IsSubmit:     post.IsSubmit,
	}
}

// ToDomainPubPost 将dao层转化为领域层
func ToDomainPubPost(post dao.PubPost) domain.Post {
	return domain.Post{
		ID:           post.ID,
		Title:        post.Title,
		Content:      post.Content,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
		DeletedAt:    sql.NullTime(post.DeletedAt),
		Status:       post.Status,
		Uid:          post.Uid,
		PlateID:      post.PlateID,
		Slug:         post.Slug,
		CategoryID:   post.CategoryID,
		Tags:         post.Tags,
		CommentCount: post.CommentCount,
	}
}

// ToDomainListPubPost 将dao层转化为领域层
func ToDomainListPubPost(post dao.PubPost) domain.Post {
	return domain.Post{
		ID:           post.ID,
		Title:        post.Title,
		Content:      post.Content,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
		Status:       post.Status,
		Uid:          post.Uid,
		PlateID:      post.PlateID,
		Slug:         post.Slug,
		CategoryID:   post.CategoryID,
		Tags:         post.Tags,
		CommentCount: post.CommentCount,
	}
}
