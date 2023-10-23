package model

type NewPeopleInfoReq struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

type UpdatePeopleInfoReq struct {
	ID         int     `json:"ID"`
	Name       *string `json:"name"`
	Surname    *string `json:"surname"`
	Patronymic *string `json:"patronymic"`
}

type NewPeople struct {
	Name       string `db:"name"`
	Surname    string `db:"surname"`
	Patronymic string `db:"patronymic"`
	Age        int    `db:"age"`
	Gender     string `db:"gender"`
	CountryID  string `db:"country_id"`
}

type People struct {
	ID         int    `json:"id" db:"people_id"`
	Name       string `json:"name" db:"name"`
	Surname    string `json:"surname" db:"surname"`
	Patronymic string `json:"patronymic,omitempty" db:"patronymic"`
	Age        int    `json:"age" db:"age"`
	Gender     string `json:"gender" db:"gender"`
	CountryID  string `json:"countryID" db:"country_id"`
}

type PeopleDataProbable struct {
	Age     int    `json:"age"`
	Gender  string `json:"gender"`
	Country [1]struct {
		CountryId string `json:"country_id"`
	} `json:"country"`
}

type PeopleFIO struct {
	Id         int
	Name       string
	Surname    string
	Patronymic string
}

type Params struct {
	MinAge int
	MaxAge int
	Gender string
	Name   string

	Limit  int
	Offset int
}

func (p *People) UpdateFIO(data UpdatePeopleInfoReq) {
	if data.Name != nil {
		p.Name = *data.Name
	}
	if data.Surname != nil {
		p.Surname = *data.Surname
	}
	if data.Patronymic != nil {
		p.Patronymic = *data.Patronymic
	}
}
