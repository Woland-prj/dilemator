package data

type Link struct {
	Name string
	Href string
}

func LandingMenuLinks() []Link {
	return []Link{
		{Name: "How It Works", Href: "#how-it-works"},
		{Name: "Use Cases", Href: "#use-cases"},
		{Name: "Features", Href: "#features"},
		{Name: "Pricing", Href: "#pricing"},
		{Name: "Contact", Href: "#contact"},
	}
}

func PlatformMenuLinks() []Link {
	return []Link{
		{Name: "Test1", Href: "/test"},
		{Name: "Tes2", Href: "/test2"},
	}
}
