package category

type Name string

// enumer not necessary, won't be updated often, have to have FromString anyways
const (
	UNDEFINED Name = "undefined"
	GENERAL   Name = "general"
	IMAGES    Name = "images"
	SCIENCE   Name = "science"
	THOROUGH  Name = "thorough"
)
