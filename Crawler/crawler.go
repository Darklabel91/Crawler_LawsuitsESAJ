package Crawler

import "github.com/tebeka/selenium"

const (
	InitWebSite    = "https://esaj.tjsp.jus.br/cpo"
	EndWebSite     = "g/open.do"
	xpathRecaptcha = "//*[@id=\"headerNmUsuarioLogado\"]"
)

type EntireLawsuit struct {
	LawsuitNumber string
	Warning       string
	FirstDegree   Lawsuit
	SecondDegree  Lawsuit
}

type Lawsuit struct {
	Warning   string
	Cover     LawsuitCover
	Persons   []Person
	Movements []Movement
}

func Craw(driver selenium.WebDriver, lawsuitNumber string, login string, password string) (EntireLawsuit, error) {

	var e string

	degree := "p"
	searchLink := InitWebSite + degree + EndWebSite
	fdLawsuit, err := SingleCraw(driver, searchLink, lawsuitNumber, degree, login, password)
	if err != nil {
		e += "primeiro " + err.Error()
	}

	degree = "s"
	searchLink = InitWebSite + degree + EndWebSite
	sdLawsuit, err := SingleCraw(driver, searchLink, lawsuitNumber, degree, login, password)
	if err != nil {
		e += "segundo " + err.Error()
	}

	return EntireLawsuit{
		Warning:       e,
		LawsuitNumber: lawsuitNumber,
		FirstDegree:   fdLawsuit,
		SecondDegree:  sdLawsuit,
	}, nil
}

func SingleCraw(driver selenium.WebDriver, searchLink string, lawsuit string, degree string, login string, password string) (Lawsuit, error) {
	htmlPgSrc, err := SearchLawsuit(driver, searchLink, lawsuit, degree, login, password)
	if err != nil {
		return Lawsuit{
			Warning:   err.Error(),
			Cover:     LawsuitCover{},
			Persons:   nil,
			Movements: nil,
		}, nil
	}

	secrecy := GetSecrecy(htmlPgSrc)

	if secrecy != true {
		cover, err := GetLawsuitCover(htmlPgSrc, degree)
		if err != nil {
			return Lawsuit{}, err
		}

		persons, err := GetLawsuitPersons(htmlPgSrc)
		if err != nil {
			return Lawsuit{}, err
		}

		movements, err := GetLawsuitMovements(htmlPgSrc)
		if err != nil {
			return Lawsuit{}, err
		}

		return Lawsuit{
			Warning:   "",
			Cover:     cover,
			Persons:   persons,
			Movements: movements,
		}, nil
	}

	warning := "lawsuit " + lawsuit + " is private to persons involved"
	return Lawsuit{
		Warning:   warning,
		Cover:     LawsuitCover{},
		Persons:   nil,
		Movements: nil,
	}, nil

}
