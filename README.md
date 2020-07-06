#HTTP REST API viewer

Allows viewing item list and item details provided by HTTP API with GET methods.

Supports specifing API list URL and item URL pattern along with max fetching attemptions count (3 by default)

Example

```
./bin/client_emp -list_url=http://localhost/list -item_url=http://localhost/item/%s -max_attempt=10

```

##Build

```
go build -o bin/client_emp main.go

./bin/client_emp


```

or

```
go install morrah77.com/client_emp

client_emp

interact with viewer...

rm -rf `which client_emp`
```

##Test

```
go test ./...

```

##Operate

```
./bin/client_emp
```

Right after start the program displays the list of items fetched from specified API.

Example:

```
{
  "status": "success",
  "data": [
    {
      "id": "1",
      "employee_name": "Tiger Nixon",
      "employee_salary": "320800",
      "employee_age": "61",
      "profile_image": ""
    },
    {
      "id": "2",
      "employee_name": "Garrett Winters",
      "employee_salary": "170750",
      "employee_age": "63",
      "profile_image": ""
    }
  ]
}
```

Follow the hints displayed b program:

```
Please type ID to view item details,
`l` to view list,
`q` to quit

```
type an item ID to see the item details fetched from the specified API

Example:

Input:

```
24

```

Output:

```
{
  "status": "success",
  "data": {
    "id": "24",
    "employee_name": "Doris Wilder",
    "employee_salary": "85600",
    "employee_age": "23",
    "profile_image": ""
  }
}
```

When the typed ID is invalid or the program cannot fetch appropriate item details it displays the error message.

Type `l` to fetch the item list.

Tpe `q` to quit.
