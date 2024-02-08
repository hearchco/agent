package category

type Name string

// enumer not necessary, won't be updated often, have to have FromString anyways
const (
	UNDEFINED Name = "undefined"
	GENERAL   Name = "general"
	INFO      Name = "info"
	SCIENCE   Name = "science"
	NEWS      Name = "news"
	BLOG      Name = "blog"
	SURF      Name = "surf"
	NEWNEWS   Name = "newnews"
)
