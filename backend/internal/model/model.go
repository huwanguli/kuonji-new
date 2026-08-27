package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User 用户模型
type User struct {
	BaseModel
	Username string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
	Email    string `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Role     string `gorm:"size:20;default:admin" json:"role"`
}

// Article 文章模型
type Article struct {
	BaseModel
	Title       string        `gorm:"size:200;not null" json:"title"`
	Content     string        `gorm:"type:longtext" json:"content"`
	Summary     string        `gorm:"size:500" json:"summary"`
	CoverImage  string        `gorm:"size:500" json:"cover_image"`
	Status      ArticleStatus `gorm:"type:tinyint;default:0" json:"status"`
	PublishedAt *time.Time    `json:"published_at"`
	AuthorID    uint          `json:"author_id"`
	CategoryID  uint          `json:"category_id"`
	SeriesID    *uint         `json:"series_id"`
	SeriesOrder int           `gorm:"default:0" json:"series_order"`
	ViewCount   int           `gorm:"default:0" json:"view_count"`
	Tags        []Tag         `gorm:"many2many:article_tags;" json:"tags"`
	Category    Category      `json:"category"`
	Series      *Series       `json:"series,omitempty"`
	Comments    []Comment     `json:"comments,omitempty"`
}

// ArticleStatus 文章状态
type ArticleStatus int

const (
	ArticleStatusDraft     ArticleStatus = 0
	ArticleStatusPublished ArticleStatus = 1
	ArticleStatusScheduled ArticleStatus = 2
)

// CommentStatus 评论状态
const (
	CommentStatusPending  int8 = 0
	CommentStatusApproved int8 = 1
	CommentStatusSpam     int8 = 2
)

// Category 分类模型
type Category struct {
	BaseModel
	Name        string     `gorm:"size:100;not null" json:"name"`
	Description string     `gorm:"size:500" json:"description"`
	ParentID    *uint      `json:"parent_id"`
	Sort        int        `gorm:"default:0" json:"sort"`
	Parent      *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Articles    []Article  `json:"articles,omitempty"`
}

// Tag 标签模型
type Tag struct {
	BaseModel
	Name     string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Articles []Article `gorm:"many2many:article_tags;" json:"articles,omitempty"`
}

// Series 系列模型
type Series struct {
	BaseModel
	Name        string    `gorm:"size:200;not null" json:"name"`
	Description string    `gorm:"size:500" json:"description"`
	CoverImage  string    `gorm:"size:500" json:"cover_image"`
	Sort        int       `gorm:"default:0" json:"sort"`
	Articles    []Article `json:"articles,omitempty"`
}

// Comment 评论模型
type Comment struct {
	BaseModel
	ArticleID uint      `json:"article_id"`
	ParentID  *uint     `json:"parent_id"`
	Nickname  string    `gorm:"size:50;not null" json:"nickname"`
	Email     string    `gorm:"size:100;not null" json:"email"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	IP        string    `gorm:"size:45" json:"ip"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	Status    int8      `gorm:"type:tinyint;default:0" json:"status"` // 0: pending, 1: approved, 2: spam
	Article   Article   `json:"article,omitempty"`
	Parent    *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children  []Comment `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// VisitLog 访问日志模型
type VisitLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	ArticleID uint      `json:"article_id"`
	IP        string    `gorm:"size:45" json:"ip"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}
