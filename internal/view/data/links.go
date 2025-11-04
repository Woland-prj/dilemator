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
		{Name: "Platform option 1", Href: "/pl1"},
		{Name: "Platform option 2", Href: "/pl2"},
	}
}

func EditorMenuLinks() []Link {
	return []Link{
		{Name: "Editor option 1", Href: "/ed1"},
		{Name: "Editor option 2", Href: "/ed2"},
	}
}
