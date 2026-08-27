// 命令 migrate：将服务器旧博客数据（zblog 库）迁移到新单体博客库（kuonji_new）。
//
// 设计要点：
//   - 本机运行，直连服务器 MySQL（3306 已对公网开放），不改动 zblog 源数据。
//   - 目标表结构手写 SQL 创建（绕开 GORM AutoMigrate 与 MySQL uniqueIndex 的兼容性 bug），
//     与 internal/model/model.go 完全一致（单数表名、deleted_at 软删除、article_tags 关联表）。
//   - 保留原 ID（article_tags / 评论外键不断裂）。
//   - 迁移映射：
//     articles.content_md   -> content
//     articles.excerpt      -> summary
//     articles.cover        -> cover_image
//     articles.series(名称) -> series_id（按名称匹配 series 表）
//     articles.status 1     -> 1 (published)
//     published_at          = created_at
//     author_id             = 1（新管理员 Kanon）
//     丢弃 slug / content_html / is_top / is_announcement
//     comments.author       -> nickname
//     comments.parent_author -> 丢弃，user_agent 补空
//   - 用户：不迁移旧用户（旧密码哈希未知），全新创建 Kanon（管理员，密码由 -admin-password 指定）。
//
// 用法：
//
//	go run ./cmd/migrate \
//	  -src-dsn "root:xxx@tcp(8.218.174.69:3306)/zblog?charset=utf8mb4&parseTime=true" \
//	  -dst-dsn "root:xxx@tcp(8.218.174.69:3306)/kuonji_new?charset=utf8mb4&parseTime=true" \
//	  -admin-user Kanon -admin-password 231310627 -admin-email kanon@kuonji.xyz
//
// 默认 -src-dsn / -dst-dsn 即服务器地址；-force 会先 DROP 目标表重建。
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultSrcDSN = "root:231310627@tcp(8.218.174.69:3306)/zblog?charset=utf8mb4&parseTime=true&loc=Local"
	defaultDstDSN = "root:231310627@tcp(8.218.174.69:3306)/kuonji_new?charset=utf8mb4&parseTime=true&loc=Local"
)

// 目标库建表语句（与 model.go 完全对应；GORM 配置 SingularTable: true，故表名为单数）。
var createTables = []string{
	`CREATE TABLE IF NOT EXISTS user (
		id         bigint unsigned NOT NULL AUTO_INCREMENT,
		created_at datetime(3) NULL,
		updated_at datetime(3) NULL,
		deleted_at datetime(3) NULL,
		username   varchar(50)  NOT NULL,
		password   varchar(255) NOT NULL,
		email      varchar(100) NOT NULL,
		role       varchar(20)  DEFAULT 'admin',
		PRIMARY KEY (id),
		UNIQUE KEY idx_user_username (username),
		UNIQUE KEY idx_user_email (email),
		KEY idx_user_deleted_at (deleted_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

	`CREATE TABLE IF NOT EXISTS category (
		id          bigint unsigned NOT NULL AUTO_INCREMENT,
		created_at  datetime(3) NULL,
		updated_at  datetime(3) NULL,
		deleted_at  datetime(3) NULL,
		name        varchar(100) NOT NULL,
		description varchar(500),
		parent_id   bigint unsigned NULL,
		sort        bigint DEFAULT 0,
		PRIMARY KEY (id),
		KEY idx_category_deleted_at (deleted_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

	`CREATE TABLE IF NOT EXISTS tag (
		id         bigint unsigned NOT NULL AUTO_INCREMENT,
		created_at datetime(3) NULL,
		updated_at datetime(3) NULL,
		deleted_at datetime(3) NULL,
		name       varchar(50) NOT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY idx_tag_name (name),
		KEY idx_tag_deleted_at (deleted_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

	`CREATE TABLE IF NOT EXISTS series (
		id          bigint unsigned NOT NULL AUTO_INCREMENT,
		created_at  datetime(3) NULL,
		updated_at  datetime(3) NULL,
		deleted_at  datetime(3) NULL,
		name        varchar(200) NOT NULL,
		description varchar(500),
		cover_image varchar(500),
		sort        bigint DEFAULT 0,
		PRIMARY KEY (id),
		UNIQUE KEY idx_series_name (name),
		KEY idx_series_deleted_at (deleted_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

	`CREATE TABLE IF NOT EXISTS article (
		id           bigint unsigned NOT NULL AUTO_INCREMENT,
		created_at   datetime(3) NULL,
		updated_at   datetime(3) NULL,
		deleted_at   datetime(3) NULL,
		title        varchar(200) NOT NULL,
		content      longtext,
		summary      varchar(500),
		cover_image  varchar(500),
		status       tinyint DEFAULT 0,
		published_at datetime(3) NULL,
		author_id    bigint unsigned,
		category_id  bigint unsigned,
		series_id    bigint unsigned NULL,
		series_order bigint DEFAULT 0,
		view_count   bigint DEFAULT 0,
		PRIMARY KEY (id),
		KEY idx_article_deleted_at (deleted_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

	`CREATE TABLE IF NOT EXISTS article_tags (
		article_id bigint unsigned NOT NULL,
		tag_id     bigint unsigned NOT NULL,
		PRIMARY KEY (article_id, tag_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

	`CREATE TABLE IF NOT EXISTS comment (
		id         bigint unsigned NOT NULL AUTO_INCREMENT,
		created_at datetime(3) NULL,
		updated_at datetime(3) NULL,
		deleted_at datetime(3) NULL,
		article_id bigint unsigned,
		parent_id  bigint unsigned NULL,
		nickname   varchar(50) NOT NULL,
		email      varchar(100) NOT NULL,
		content    text NOT NULL,
		ip         varchar(45),
		user_agent varchar(500),
		status     tinyint DEFAULT 0,
		PRIMARY KEY (id),
		KEY idx_comment_deleted_at (deleted_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

	`CREATE TABLE IF NOT EXISTS visit_log (
		id         bigint unsigned NOT NULL AUTO_INCREMENT,
		article_id bigint unsigned,
		ip         varchar(45),
		user_agent varchar(500),
		created_at datetime(3) NULL,
		PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
}

func main() {
	var (
		srcDSN = flag.String("src-dsn", defaultSrcDSN, "source zblog DSN")
		dstDSN = flag.String("dst-dsn", defaultDstDSN, "target kuonji_new DSN")
		force  = flag.Bool("force", false, "drop target tables before creating")
		adminU = flag.String("admin-user", "Kanon", "admin username")
		adminP = flag.String("admin-password", "231310627", "admin password (bcrypt hashed)")
		adminE = flag.String("admin-email", "kanon@kuonji.xyz", "admin email")
	)
	flag.Parse()

	// 确保目标库存在（DSN 中的库名）
	// 从 DSN 中提取库名：@tcp(host:port)/dbname
	dbName, err := dbNameFromDSN(*dstDSN)
	if err != nil {
		log.Fatalf("parse dst db name: %v", err)
	}

	// 连接一个不依赖具体库的连接来建库
	adminDSN, err := dsnWithoutDB(*dstDSN)
	if err != nil {
		log.Fatalf("parse dst dsn: %v", err)
	}
	adminDB, err := sql.Open("mysql", adminDSN)
	if err != nil {
		log.Fatalf("open admin conn: %v", err)
	}
	defer adminDB.Close()
	if err := adminDB.Ping(); err != nil {
		log.Fatalf("ping admin conn (%s): %v", adminDSN, err)
	}
	if _, err := adminDB.Exec(fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)); err != nil {
		log.Fatalf("create database %s: %v", dbName, err)
	}
	log.Printf("database %s ready", dbName)
	adminDB.Close()

	// 连接目标库（建库之后）
	dst, err := sql.Open("mysql", *dstDSN)
	if err != nil {
		log.Fatalf("open dst: %v", err)
	}
	defer dst.Close()
	if err := dst.Ping(); err != nil {
		log.Fatalf("ping dst (%s): %v", *dstDSN, err)
	}
	log.Println("connected to target MySQL")

	if *force {
		for _, t := range []string{"article_tags", "comment", "article", "series", "tag", "category", "user", "visit_log"} {
			if _, err := dst.Exec("DROP TABLE IF EXISTS " + t); err != nil {
				log.Fatalf("drop %s: %v", t, err)
			}
		}
		log.Println("dropped existing target tables (-force)")
	}
	for _, stmt := range createTables {
		if _, err := dst.Exec(stmt); err != nil {
			log.Fatalf("create table: %v\nstmt: %s", err, stmt)
		}
	}
	log.Println("target tables created")

	// 源库连接
	src, err := sql.Open("mysql", *srcDSN)
	if err != nil {
		log.Fatalf("open src: %v", err)
	}
	defer src.Close()
	if err := src.Ping(); err != nil {
		log.Fatalf("ping src (%s): %v", *srcDSN, err)
	}
	log.Println("connected to source MySQL (zblog)")

	// ── 1. 用户：创建管理员 Kanon ──
	var userCount int
	if err := dst.QueryRow("SELECT COUNT(*) FROM user").Scan(&userCount); err != nil {
		log.Fatalf("count users: %v", err)
	}
	if userCount == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte(*adminP), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("bcrypt: %v", err)
		}
		if _, err := dst.Exec(
			"INSERT INTO user (id, username, password, email, role, created_at, updated_at) VALUES (1, ?, ?, ?, 'admin', NOW(3), NOW(3))",
			*adminU, string(hash), *adminE); err != nil {
			log.Fatalf("insert admin: %v", err)
		}
		log.Printf("admin user created: %s", *adminU)
	} else {
		log.Printf("user table not empty (%d), skip admin creation", userCount)
	}

	// ── 2. 分类（保留 ID；zblog 无 updated_at，用 created_at 补）──
	if _, err := dst.Exec(`
		INSERT INTO category (id, name, description, parent_id, sort, created_at, updated_at)
		SELECT id, name, COALESCE(NULLIF(description, ''), NULL), NULL, 0, created_at, created_at
		FROM zblog.categories`); err != nil {
		log.Fatalf("migrate categories: %v", err)
	}
	log.Println("categories migrated")

	// ── 3. 标签（保留 ID；zblog 无 updated_at，用 created_at 补）──
	if _, err := dst.Exec(`
		INSERT INTO tag (id, name, created_at, updated_at)
		SELECT id, name, created_at, created_at FROM zblog.tags`); err != nil {
		log.Fatalf("migrate tags: %v", err)
	}
	log.Println("tags migrated")

	// ── 4. 系列（保留 ID）──
	if _, err := dst.Exec(`
		INSERT INTO series (id, name, description, cover_image, sort, created_at, updated_at)
		SELECT id, name, COALESCE(NULLIF(description, ''), NULL), COALESCE(NULLIF(cover, ''), NULL), 0, created_at, updated_at
		FROM zblog.series`); err != nil {
		log.Fatalf("migrate series: %v", err)
	}
	log.Println("series migrated")

	// ── 5. 文章（保留 ID，series 名称 -> series_id）──
	if _, err := dst.Exec(`
		INSERT INTO article
			(id, title, content, summary, cover_image, status, published_at,
			 author_id, category_id, series_id, series_order, view_count, created_at, updated_at)
		SELECT
			a.id,
			a.title,
			a.content_md,
			COALESCE(NULLIF(a.excerpt, ''), NULL),
			COALESCE(NULLIF(a.cover, ''), NULL),
			CASE WHEN a.status = 0 THEN 0 ELSE 1 END,
			a.created_at,
			1,
			COALESCE(a.category_id, 0),
			s.id,
			COALESCE(a.series_order, 0),
			COALESCE(a.view_count, 0),
			a.created_at,
			a.updated_at
		FROM zblog.articles a
		LEFT JOIN zblog.series s ON s.name = a.series`); err != nil {
		log.Fatalf("migrate articles: %v", err)
	}
	log.Println("articles migrated")

	// ── 6. 文章-标签关联（保留原始对）──
	if _, err := dst.Exec(`
		INSERT INTO article_tags (article_id, tag_id)
		SELECT article_id, tag_id FROM zblog.article_tags`); err != nil {
		log.Fatalf("migrate article_tags: %v", err)
	}
	log.Println("article_tags migrated")

	// ── 7. 评论（author -> nickname，user_agent 补空）──
	if _, err := dst.Exec(`
		INSERT INTO comment (id, article_id, parent_id, nickname, email, content, ip, user_agent, status, created_at, updated_at)
		SELECT
			id, article_id, parent_id, author, COALESCE(email, ''), content,
			COALESCE(ip, ''), '', CASE WHEN status = 0 THEN 0 ELSE 1 END, created_at, created_at
		FROM zblog.comments`); err != nil {
		log.Fatalf("migrate comments: %v", err)
	}
	log.Println("comments migrated")

	// ── 8. 校验 ──
	verify(dst)

	log.Println("MIGRATION COMPLETED")
}

func verify(dst *sql.DB) {
	checks := []struct {
		name string
		q    string
	}{
		{"user", "SELECT COUNT(*) FROM user"},
		{"category", "SELECT COUNT(*) FROM category"},
		{"tag", "SELECT COUNT(*) FROM tag"},
		{"series", "SELECT COUNT(*) FROM series"},
		{"article", "SELECT COUNT(*) FROM article"},
		{"article_tags", "SELECT COUNT(*) FROM article_tags"},
		{"comment", "SELECT COUNT(*) FROM comment"},
	}
	for _, c := range checks {
		var n int
		if err := dst.QueryRow(c.q).Scan(&n); err != nil {
			log.Printf("verify %s: %v", c.name, err)
			continue
		}
		log.Printf("verify %-12s %d", c.name, n)
	}

	// 检查文章是否都指向存在的分类/系列/作者
	var orphanCats, orphanSeries, orphanAuthors, orphanTags int
	_ = dst.QueryRow(`SELECT COUNT(*) FROM article a LEFT JOIN category c ON c.id = a.category_id WHERE a.category_id <> 0 AND c.id IS NULL`).Scan(&orphanCats)
	_ = dst.QueryRow(`SELECT COUNT(*) FROM article a LEFT JOIN series s ON s.id = a.series_id WHERE a.series_id IS NOT NULL AND s.id IS NULL`).Scan(&orphanSeries)
	_ = dst.QueryRow(`SELECT COUNT(*) FROM article a LEFT JOIN user u ON u.id = a.author_id WHERE u.id IS NULL`).Scan(&orphanAuthors)
	_ = dst.QueryRow(`SELECT COUNT(*) FROM article_tags t LEFT JOIN tag tg ON tg.id = t.tag_id WHERE tg.id IS NULL`).Scan(&orphanTags)
	log.Printf("orphans -> categories:%d series:%d authors:%d tags:%d", orphanCats, orphanSeries, orphanAuthors, orphanTags)
}

func dbNameFromDSN(dsn string) (string, error) {
	// 简单解析：/dbname?...
	slash := -1
	for i := 0; i+1 < len(dsn); i++ {
		if dsn[i] == '/' && dsn[i+1] != '/' { // 排除协议头
			slash = i
			break
		}
	}
	if slash < 0 {
		return "", fmt.Errorf("no db name in dsn: %s", dsn)
	}
	rest := dsn[slash+1:]
	for i, c := range rest {
		if c == '?' || c == '/' {
			return rest[:i], nil
		}
	}
	return rest, nil
}

func dsnWithoutDB(dsn string) (string, error) {
	slash := -1
	for i := 0; i+1 < len(dsn); i++ {
		if dsn[i] == '/' && dsn[i+1] != '/' {
			slash = i
			break
		}
	}
	if slash < 0 {
		return "", fmt.Errorf("no db in dsn: %s", dsn)
	}
	prefix := dsn[:slash]
	rest := dsn[slash+1:]
	q := -1
	for i, c := range rest {
		if c == '?' {
			q = i
			break
		}
	}
	if q >= 0 {
		return prefix + "/?" + rest[q+1:], nil
	}
	return prefix + "/", nil
}
