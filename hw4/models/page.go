package models

// Page - страница доступная шаблонизатору
type Page struct {
	Title      string
	Posts      BlogPosts
	Categories Categories
}
