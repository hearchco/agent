package category

type Name string

// enumer not necessary, won't be updated often, have to have FromString anyways
const (
	UNDEFINED Name = "undefined"
	GENERAL   Name = "general"
	IMAGES    Name = "images"
	QUICK     Name = "quick"
	SCIENCE   Name = "science"
	BROAD     Name = "broad"
)
