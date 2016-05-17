package v1

type Person struct {
	UserId string
	Age int
	Username string
}

func List() (persons []Person){
	persons = make([]Person, 0)
	persons = append(persons,Person{UserId:"helloworld",Age:18,Username:"johnsun"})
	return  persons
}

func Get(userId string) (person *Person){
	person = &Person{UserId:userId,Age:18,Username:"johnsun"}
	return  person
}
